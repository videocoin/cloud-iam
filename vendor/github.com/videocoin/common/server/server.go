package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // anonymous import to get the pprof handler registered

	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	node_https "github.com/prometheus/node_exporter/https"
	"github.com/weaveworks/common/instrument"
	"golang.org/x/net/context"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/videocoin/common/logging"
	"github.com/videocoin/common/middleware"
	"github.com/videocoin/common/signals"
)

// Server wraps a HTTP and gRPC server, and some common initialization.
//
// Servers will be automatically instrumented for Prometheus metrics.
type Server struct {
	cfg          Config
	handler      *signals.Handler
	grpcListener net.Listener
	httpListener net.Listener

	HTTP       *mux.Router
	HTTPServer *http.Server
	GRPC       *grpc.Server
	Log        logging.Interface
}

// New makes a new Server
func New(cfg Config) (*Server, error) {
	// Setup listeners first, so we can fail early if the port is in use.
	httpListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.HTTPListenAddress, cfg.HTTPListenPort))
	if err != nil {
		return nil, err
	}
	if cfg.HTTPConnLimit > 0 {
		httpListener = netutil.LimitListener(httpListener, cfg.HTTPConnLimit)
	}

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPCListenAddress, cfg.GRPCListenPort))
	if err != nil {
		return nil, err
	}

	if cfg.GRPCConnLimit > 0 {
		grpcListener = netutil.LimitListener(grpcListener, cfg.GRPCConnLimit)
	}

	// If user doesn't supply a logging implementation, by default instantiate
	// logrus.
	log := cfg.Log
	if log == nil {
		log = logging.NewLogrus(cfg.LogLevel)
	}

	// Setup TLS
	var httpTLSConfig *tls.Config
	if len(cfg.HTTPTLSConfig.TLSCertPath) > 0 && len(cfg.HTTPTLSConfig.TLSKeyPath) > 0 {
		// Note: ConfigToTLSConfig from prometheus/node_exporter is awaiting security review.
		httpTLSConfig, err = node_https.ConfigToTLSConfig(&cfg.HTTPTLSConfig)
		if err != nil {
			return nil, fmt.Errorf("error generating http tls config: %v", err)
		}
	}
	var grpcTLSConfig *tls.Config
	if len(cfg.GRPCTLSConfig.TLSCertPath) > 0 && len(cfg.GRPCTLSConfig.TLSKeyPath) > 0 {
		// Note: ConfigToTLSConfig from prometheus/node_exporter is awaiting security review.
		grpcTLSConfig, err = node_https.ConfigToTLSConfig(&cfg.GRPCTLSConfig)
		if err != nil {
			return nil, fmt.Errorf("error generating grpc tls config: %v", err)
		}
	}

	// Prometheus histograms for requests.
	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: cfg.MetricsNamespace,
		Name:      "request_duration_seconds",
		Help:      "Time (in seconds) spent serving HTTP requests.",
		Buckets:   instrument.DefBuckets,
	}, []string{"method", "route", "status_code", "ws"})
	prometheus.MustRegister(requestDuration)

	log.WithField("http", httpListener.Addr()).WithField("grpc", grpcListener.Addr()).Infof("server listening on addresses")

	// Setup gRPC server
	serverLog := middleware.GRPCServerLog{
		WithRequest: !cfg.ExcludeRequestInLog,
		Log:         log,
	}
	grpcMiddleware := []grpc.UnaryServerInterceptor{
		serverLog.UnaryServerInterceptor,
		middleware.UnaryServerInstrumentInterceptor(requestDuration),
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
	}
	grpcMiddleware = append(grpcMiddleware, cfg.GRPCMiddleware...)

	grpcStreamMiddleware := []grpc.StreamServerInterceptor{
		serverLog.StreamServerInterceptor,
		middleware.StreamServerInstrumentInterceptor(requestDuration),
		otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer()),
	}
	grpcStreamMiddleware = append(grpcStreamMiddleware, cfg.GRPCStreamMiddleware...)

	grpcKeepAliveOptions := keepalive.ServerParameters{
		MaxConnectionIdle:     cfg.GRPCServerMaxConnectionIdle,
		MaxConnectionAge:      cfg.GRPCServerMaxConnectionAge,
		MaxConnectionAgeGrace: cfg.GRPCServerMaxConnectionAgeGrace,
		Time:    cfg.GRPCServerTime,
		Timeout: cfg.GRPCServerTimeout,
	}

	grpcOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcMiddleware...,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpcStreamMiddleware...,
		)),
		grpc.KeepaliveParams(grpcKeepAliveOptions),
		grpc.MaxRecvMsgSize(cfg.GPRCServerMaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.GRPCServerMaxSendMsgSize),
		grpc.MaxConcurrentStreams(uint32(cfg.GPRCServerMaxConcurrentStreams)),
	}
	grpcOptions = append(grpcOptions, cfg.GRPCOptions...)
	if grpcTLSConfig != nil {
		grpcCreds := credentials.NewTLS(grpcTLSConfig)
		grpcOptions = append(grpcOptions, grpc.Creds(grpcCreds))
	}
	grpcServer := grpc.NewServer(grpcOptions...)

	// Setup HTTP server
	router := mux.NewRouter()
	if cfg.PathPrefix != "" {
		// Expect metrics and pprof handlers to be prefixed with server's path prefix.
		// e.g. /loki/metrics or /loki/debug/pprof
		router = router.PathPrefix(cfg.PathPrefix).Subrouter()
	}
	if cfg.RegisterInstrumentation {
		RegisterInstrumentation(router)
	}
	httpMiddleware := []middleware.Interface{
		middleware.Tracer{
			RouteMatcher: router,
		},
		middleware.Log{
			Log: log,
		},
		middleware.Instrument{
			Duration:     requestDuration,
			RouteMatcher: router,
		},
	}

	httpMiddleware = append(httpMiddleware, cfg.HTTPMiddleware...)
	httpServer := &http.Server{
		ReadTimeout:  cfg.HTTPServerReadTimeout,
		WriteTimeout: cfg.HTTPServerWriteTimeout,
		IdleTimeout:  cfg.HTTPServerIdleTimeout,
		Handler:      middleware.Merge(httpMiddleware...).Wrap(router),
	}
	if httpTLSConfig != nil {
		httpServer.TLSConfig = httpTLSConfig
	}

	return &Server{
		cfg:          cfg,
		httpListener: httpListener,
		grpcListener: grpcListener,
		handler:      signals.NewHandler(log),

		HTTP:       router,
		HTTPServer: httpServer,
		GRPC:       grpcServer,
		Log:        log,
	}, nil
}

// RegisterInstrumentation on the given router.
func RegisterInstrumentation(router *mux.Router) {
	router.Handle("/metrics", promhttp.Handler())
	router.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)
}

// Run the server; blocks until SIGTERM or an error is received.
func (s *Server) Run() error {
	errChan := make(chan error, 1)

	// Wait for a signal
	go func() {
		s.handler.Loop()
		select {
		case errChan <- nil:
		default:
		}
	}()

	go func() {
		var err error
		if s.HTTPServer.TLSConfig == nil {
			err = s.HTTPServer.Serve(s.httpListener)
		} else {
			err = s.HTTPServer.ServeTLS(s.httpListener, s.cfg.HTTPTLSConfig.TLSCertPath, s.cfg.HTTPTLSConfig.TLSKeyPath)
		}
		if err == http.ErrServerClosed {
			err = nil
		}

		select {
		case errChan <- err:
		default:
		}
	}()

	// Setup gRPC server
	// for HTTP over gRPC, ensure we don't double-count the middleware
	// httpgrpc.RegisterHTTPServer(s.GRPC, httpgrpc_server.NewServer(s.HTTP))

	go func() {
		err := s.GRPC.Serve(s.grpcListener)
		if err == grpc.ErrServerStopped {
			err = nil
		}

		select {
		case errChan <- err:
		default:
		}
	}()

	return <-errChan
}

// Stop unblocks Run().
func (s *Server) Stop() {
	s.handler.Stop()
}

// Shutdown the server, gracefully.  Should be defered after New().
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ServerGracefulShutdownTimeout)
	defer cancel() // releases resources if httpServer.Shutdown completes before timeout elapses

	s.HTTPServer.Shutdown(ctx)
	s.GRPC.GracefulStop()
}
