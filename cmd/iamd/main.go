package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/cloud-iam/helpers"
	"github.com/videocoin/cloud-iam/service"
	"github.com/videocoin/common/grpcutil"
	jwt "github.com/videocoin/jwt-go"
	"github.com/videocoin/runtime/grpc/middleware/auth"
	iam "github.com/videocoin/videocoinapis/videocoin/iam/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const serviceName = "iam"

type Config struct {
	LogLevel        string `default:"info"`
	RPCAddr         string `default:"0.0.0.0:5000"`
	DBURI           string `required:"true"`
	Hostname        string `default:"iam.videocoin.network"`
	UserInfoURI     string `default:"https://studio.dev.videocoin.network/api/v1/user"`
	AuthTokenSecret string `required:"true"`
}

func main() {
	cfg := new(Config)
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

		pubKeyFunc := func(ctx context.Context, user string, keyID string) (interface{}, error) {
			key, err := ds.GetUserKey(user, keyID)
			if err != nil {
				return nil, err
			}
			return helpers.PubKeyFromBytesPEM(key.PublicKeyData)
		}

		serviceAccountMatchingFunc := func(token *jwt.Token) bool {
			_, ok := token.Header["kid"]
			return ok
		}
		authOpts := []auth.AuthOption{
			auth.WithAuthentication(auth.HMACJWT(cfg.AuthTokenSecret)),
			auth.WithAuthentication(auth.ServiceAccount(cfg.Hostname, cfg.AuthTokenSecret, pubKeyFunc), serviceAccountMatchingFunc),
			auth.WithAuthorization(auth.RBAC(cfg.UserInfoURI)),
		}
		grpcSrv = grpc.NewServer(grpcutil.DefaultServerOptsWithAuth(log.NewEntry(log.StandardLogger()), auth.NewAuthnzHandler(authOpts...).HandleAuthnz)...)

		iam.RegisterIAMServer(grpcSrv, service.New(ds))
		healthpb.RegisterHealthServer(grpcSrv, healthSrv)
		healthSrv.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", serviceName), healthpb.HealthCheckResponse_SERVING)

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

	healthSrv.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", serviceName), healthpb.HealthCheckResponse_NOT_SERVING)

	if grpcSrv != nil {
		grpcSrv.GracefulStop()
	}

	return errgrp.Wait()
}
