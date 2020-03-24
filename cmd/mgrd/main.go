package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	admin "github.com/videocoin/cloud-iam/api/admin/v1"
	logz "github.com/videocoin/common/log"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/common/grpcutil"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

var (
	// ServiceName is the service name.
	ServiceName = "mgr"

	// Version is the application version.
	Version = "dev"
)

// Config is the global config.
type Config struct {
	LogLevel string `default:"info"`
	RPCAddr  string `default:"0.0.0.0:5000"`
	DBURI    string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := new(Config)
	if err := envconfig.Process(ServiceName, cfg); err != nil {
		return err
	}

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return err
	}
	logger := logz.NewLogrus(lvl)
	entry := logrus.NewEntry(logger)
	logz.SetGlobal(logger)

	ctx, cancel := context.WithCancel(context.Background())
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

		grpcSrv = grpc.NewServer(grpcutil.DefaultServerOpts(entry)...)
		admin.RegisterIAMServer(grpcSrv, NewServer(ds))
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
		return err
	}

	return nil
}
