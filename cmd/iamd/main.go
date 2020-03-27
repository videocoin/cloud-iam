package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/cloud-iam/helpers"
	"github.com/videocoin/cloud-iam/service"
	"github.com/videocoin/common/grpcutil"
	log "github.com/videocoin/common/log"
	"github.com/videocoin/runtime/security"
	iam "github.com/videocoin/videocoinapis/videocoin/iam/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var ServiceName = "iam"

type Config struct {
	LogLevel        string `default:"info"`
	RPCAddr         string `default:"0.0.0.0:5000"`
	Hostname        string `default:"iam.videocoin.network"`
	DBURI           string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local"`
	AuthTokenSecret string `required:"true"`
	entry           *logrus.Entry
}

func main() {
	cfg := new(Config)
	if err := envconfig.Process(ServiceName, cfg); err != nil {
		fatal(err)
	}

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		fatal(err)
	}
	logger := log.NewLogrus(lvl)

	cfg.entry = logrus.NewEntry(logger)
	log.SetGlobal(logger)

	if err := run(cfg); err != nil {
		log.Fatalln(err)
	}
}

func run(cfg *Config) error {
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

		pubKeyFunc := func(ctx context.Context, subject string, keyID string) (interface{}, error) {
			key, err := ds.GetUserKey(subject, keyID)
			if err != nil {
				return nil, err
			}
			return helpers.PubKeyFromBytesPEM(key.PublicKeyData)
		}
		grpcSrv = grpc.NewServer(grpcutil.DefaultServerOptsWithAuth(cfg.entry, security.Authnz(cfg.Hostname, cfg.AuthTokenSecret, pubKeyFunc))...)

		iam.RegisterIAMServer(grpcSrv, service.NewServer(ds))
		healthpb.RegisterHealthServer(grpcSrv, healthSrv)
		healthSrv.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", ServiceName), healthpb.HealthCheckResponse_SERVING)

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

	return errgrp.Wait()
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}
