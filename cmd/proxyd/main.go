package main

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const serviceName = "proxy"

type config struct {
	LogLevel   string `default:"info"`
	ListenAddr string `default:":8080"`
	GRPCAddr   string `default:"0.0.0.0:5000"`
}

func main() {
	cfg := new(config)
	if err := envconfig.Process(serviceName, cfg); err != nil {
		log.Fatal(err)
	}

	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			log.Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealthCheck(conn))

	gw := runtime.NewServeMux()
	mux.Handle("/", gw)

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: allowCORS(mux),
	}

	go func() {
		<-ctx.Done()
		log.Info("Shutting down the http server")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	log.Infof("Starting listening at %s", cfg.ListenAddr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Errorf("Failed to listen and serve: %v", err)
		return err
	}

	return nil
}
