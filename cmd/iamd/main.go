package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	iam "github.com/videocoin/videocoinapis/videocoin/iam/v1"
	logz "github.com/videocoin/common/log"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/cloud-iam/service"
	"github.com/videocoin/common/grpcutil"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

var (
	// ServiceName is the service name.
	ServiceName = "iam"

	// Version is the application version.
	Version = "dev"
)

// Config is the global config.
type Config struct {
	LogLevel             string `default:"info" envconfig:"LOG_LEVEL"`
	RPCAddr              string `default:"0.0.0.0:5000"`
	DBURI                string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local"`
	EncryptionPassphrase string `required:"true"`
}

func main() {
	cfg := new(Config)
	if err := envconfig.Process(ServiceName, cfg); err != nil {
		log.Fatal(err)
	}

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger := logrus.NewEntry(logz.NewLogrus(lvl))
	logz.SetGlobal(logz.NewLogrus(lvl))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	errgrp, ctx := errgroup.WithContext(ctx)

	healthSrv := health.NewServer()
	var grpcSrv *grpc.Server
	errgrp.Go(func() error {
		ds, err := datastore.Open(cfg.DBURI)
		if err != nil {
			return err
		}
		defer ds.Close()

		grpcSrv = grpc.NewServer(grpcutil.DefaultServerOpts(logger)...)
		iam.RegisterIAMServer(grpcSrv, service.NewServer(ds, cfg.EncryptionPassphrase))
		healthpb.RegisterHealthServer(grpcSrv, healthSrv)

		healthSrv.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", ServiceName), healthpb.HealthCheckResponse_SERVING)

		grpc_prometheus.Register(grpcSrv)
		http.Handle("/metrics", promhttp.Handler())

		lis, err := net.Listen("tcp", cfg.RPCAddr)
		if err != nil {
			return err
		}
		return grpcSrv.Serve(lis)
	})

	select {
	case <-stop:
		break
	case <-ctx.Done():
		break
	}

	cancel()

	healthSrv.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", ServiceName), healthpb.HealthCheckResponse_NOT_SERVING)

	if grpcSrv != nil {
		grpcSrv.GracefulStop()
	}

	if err := errgrp.Wait(); err != nil {
		log.Fatal(err)
	}
}
