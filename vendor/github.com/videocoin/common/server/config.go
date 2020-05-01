package server

import (
	"math"
	"time"

	node_https "github.com/prometheus/node_exporter/https"
	"google.golang.org/grpc"

	"github.com/videocoin/common/logging"
	"github.com/videocoin/common/middleware"
)

var infinty = time.Duration(math.MaxInt64)

// DefaultConfig ...
var DefaultConfig = Config{
	HTTPListenPort:                  80,
	HTTPConnLimit:                   0,
	GRPCListenPort:                  9095,
	GRPCConnLimit:                   0,
	RegisterInstrumentation:         true,
	ServerGracefulShutdownTimeout:   30 * time.Second,
	HTTPServerReadTimeout:           30 * time.Second,
	HTTPServerWriteTimeout:          30 * time.Second,
	HTTPServerIdleTimeout:           120 * time.Second,
	GPRCServerMaxRecvMsgSize:        4 * 1024 * 1024,
	GRPCServerMaxSendMsgSize:        4 * 1024 * 1024,
	GPRCServerMaxConcurrentStreams:  100,
	GRPCServerMaxConnectionIdle:     infinty,
	GRPCServerMaxConnectionAge:      infinty,
	GRPCServerMaxConnectionAgeGrace: infinty,
	GRPCServerTime:                  2 * time.Hour,
	GRPCServerTimeout:               20 * time.Second,
}

// Config for a Server
type Config struct {
	MetricsNamespace  string `yaml:"-"`
	HTTPListenAddress string `yaml:"http_listen_address"`
	HTTPListenPort    int    `yaml:"http_listen_port"`
	HTTPConnLimit     int    `yaml:"http_listen_conn_limit"`
	GRPCListenAddress string `yaml:"grpc_listen_address"`
	GRPCListenPort    int    `yaml:"grpc_listen_port"`
	GRPCConnLimit     int    `yaml:"grpc_listen_conn_limit"`

	HTTPTLSConfig node_https.TLSStruct `yaml:"http_tls_config"`
	GRPCTLSConfig node_https.TLSStruct `yaml:"grpc_tls_config"`

	RegisterInstrumentation bool `yaml:"register_instrumentation"`
	ExcludeRequestInLog     bool `yaml:"-"`

	ServerGracefulShutdownTimeout time.Duration `yaml:"graceful_shutdown_timeout"`
	HTTPServerReadTimeout         time.Duration `yaml:"http_server_read_timeout"`
	HTTPServerWriteTimeout        time.Duration `yaml:"http_server_write_timeout"`
	HTTPServerIdleTimeout         time.Duration `yaml:"http_server_idle_timeout"`

	GRPCOptions          []grpc.ServerOption            `yaml:"-"`
	GRPCMiddleware       []grpc.UnaryServerInterceptor  `yaml:"-"`
	GRPCStreamMiddleware []grpc.StreamServerInterceptor `yaml:"-"`
	HTTPMiddleware       []middleware.Interface         `yaml:"-"`

	GPRCServerMaxRecvMsgSize        int           `yaml:"grpc_server_max_recv_msg_size"`
	GRPCServerMaxSendMsgSize        int           `yaml:"grpc_server_max_send_msg_size"`
	GPRCServerMaxConcurrentStreams  uint          `yaml:"grpc_server_max_concurrent_streams"`
	GRPCServerMaxConnectionIdle     time.Duration `yaml:"grpc_server_max_connection_idle"`
	GRPCServerMaxConnectionAge      time.Duration `yaml:"grpc_server_max_connection_age"`
	GRPCServerMaxConnectionAgeGrace time.Duration `yaml:"grpc_server_max_connection_age_grace"`
	GRPCServerTime                  time.Duration `yaml:"grpc_server_keepalive_time"`
	GRPCServerTimeout               time.Duration `yaml:"grpc_server_keepalive_timeout"`

	LogLevel logging.Level     `yaml:"log_level"`
	Log      logging.Interface `yaml:"-"`

	PathPrefix string `yaml:"http_path_prefix"`
}

/*
// RegisterFlags adds the flags required to config this to the given FlagSet
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&cfg.HTTPListenAddress, "server.http-listen-address", "", "HTTP server listen address.")
	f.StringVar(&cfg.HTTPTLSConfig.TLSCertPath, "server.http-tls-cert-path", "", "HTTP server cert path.")
	f.StringVar(&cfg.HTTPTLSConfig.TLSKeyPath, "server.http-tls-key-path", "", "HTTP server key path.")
	f.StringVar(&cfg.HTTPTLSConfig.ClientAuth, "server.http-tls-client-auth", "", "HTTP TLS Client Auth type.")
	f.StringVar(&cfg.HTTPTLSConfig.ClientCAs, "server.http-tls-ca-path", "", "HTTP TLS Client CA path.")
	f.StringVar(&cfg.GRPCTLSConfig.TLSCertPath, "server.grpc-tls-cert-path", "", "GRPC TLS server cert path.")
	f.StringVar(&cfg.GRPCTLSConfig.TLSKeyPath, "server.grpc-tls-key-path", "", "GRPC TLS server key path.")
	f.StringVar(&cfg.GRPCTLSConfig.ClientAuth, "server.grpc-tls-client-auth", "", "GRPC TLS Client Auth type.")
	f.StringVar(&cfg.GRPCTLSConfig.ClientCAs, "server.grpc-tls-ca-path", "", "GRPC TLS Client CA path.")
	f.IntVar(&cfg.HTTPListenPort, "server.http-listen-port", 80, "HTTP server listen port.")
	f.IntVar(&cfg.HTTPConnLimit, "server.http-conn-limit", 0, "Maximum number of simultaneous http connections, <=0 to disable")
	f.StringVar(&cfg.GRPCListenAddress, "server.grpc-listen-address", "", "gRPC server listen address.")
	f.IntVar(&cfg.GRPCListenPort, "server.grpc-listen-port", 9095, "gRPC server listen port.")
	f.IntVar(&cfg.GRPCConnLimit, "server.grpc-conn-limit", 0, "Maximum number of simultaneous grpc connections, <=0 to disable")
	f.BoolVar(&cfg.RegisterInstrumentation, "server.register-instrumentation", true, "Register the intrumentation handlers (/metrics etc).")
	f.DurationVar(&cfg.ServerGracefulShutdownTimeout, "server.graceful-shutdown-timeout", 30*time.Second, "Timeout for graceful shutdowns")
	f.DurationVar(&cfg.HTTPServerReadTimeout, "server.http-read-timeout", 30*time.Second, "Read timeout for HTTP server")
	f.DurationVar(&cfg.HTTPServerWriteTimeout, "server.http-write-timeout", 30*time.Second, "Write timeout for HTTP server")
	f.DurationVar(&cfg.HTTPServerIdleTimeout, "server.http-idle-timeout", 120*time.Second, "Idle timeout for HTTP server")
	f.IntVar(&cfg.GPRCServerMaxRecvMsgSize, "server.grpc-max-recv-msg-size-bytes", 4*1024*1024, "Limit on the size of a gRPC message this server can receive (bytes).")
	f.IntVar(&cfg.GRPCServerMaxSendMsgSize, "server.grpc-max-send-msg-size-bytes", 4*1024*1024, "Limit on the size of a gRPC message this server can send (bytes).")
	f.UintVar(&cfg.GPRCServerMaxConcurrentStreams, "server.grpc-max-concurrent-streams", 100, "Limit on the number of concurrent streams for gRPC calls (0 = unlimited)")
	f.DurationVar(&cfg.GRPCServerMaxConnectionIdle, "server.grpc.keepalive.max-connection-idle", infinty, "The duration after which an idle connection should be closed. Default: infinity")
	f.DurationVar(&cfg.GRPCServerMaxConnectionAge, "server.grpc.keepalive.max-connection-age", infinty, "The duration for the maximum amount of time a connection may exist before it will be closed. Default: infinity")
	f.DurationVar(&cfg.GRPCServerMaxConnectionAgeGrace, "server.grpc.keepalive.max-connection-age-grace", infinty, "An additive period after max-connection-age after which the connection will be forcibly closed. Default: infinity")
	f.DurationVar(&cfg.GRPCServerTime, "server.grpc.keepalive.time", time.Hour*2, "Duration after which a keepalive probe is sent in case of no activity over the connection., Default: 2h")
	f.DurationVar(&cfg.GRPCServerTimeout, "server.grpc.keepalive.timeout", time.Second*20, "After having pinged for keepalive check, the duration after which an idle connection should be closed, Default: 20s")
	f.StringVar(&cfg.PathPrefix, "server.path-prefix", "", "Base path to serve all API routes from (e.g. /v1/)")
	cfg.LogLevel.RegisterFlags(f)
}*/
