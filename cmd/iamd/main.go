package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kelseyhightower/envconfig"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/cloud-iam/service"
	"github.com/videocoin/common/logging"
	"github.com/videocoin/common/server"
	iam "github.com/videocoin/videocoinapis/videocoin/iam/v1"
	v1 "github.com/videocoin/videocoinapis/videocoin/iam/v1"
)

const serviceName = "iam"

type Config struct {
	HTTPListenPort  int    `default:"8080"`
	GRPCListenPort  int    `default:"9095"`
	LogLevel        string `default:"info"`
	DBURI           string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local"`
	Hostname        string `default:"iam.videocoin.network"`
	AuthTokenSecret string `required:"true"`
	KeyLimitPerUser int    `default:"10"`
}

func main() {
	cfg := new(Config)
	if err := envconfig.Process(serviceName, cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := logging.Setup(cfg.LogLevel); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		logging.Fatalln(err)
	}
}

func run(cfg *Config) error {
	servercfg := server.DefaultConfig
	servercfg.HTTPListenPort = cfg.HTTPListenPort
	servercfg.GRPCListenPort = cfg.GRPCListenPort
	servercfg.Log = logging.Global()

	srv, err := server.New(servercfg)
	if err != nil {
		return err
	}
	defer srv.Shutdown()

	ds, err := datastore.Open(cfg.DBURI)
	if err != nil {
		return err
	}
	defer ds.Close()

	/*
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
				auth.WithAuthorization(auth.RBAC()),
			}
	*/

	svc := service.New(ds)
	iam.RegisterIAMServer(srv.GRPC, svc)
	healthpb.RegisterHealthServer(srv.GRPC, svc)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv.HTTP.HandleFunc("/health", handleHealthCheck(svc))

	mux := runtime.NewServeMux()
	if err := v1.RegisterIAMHandlerServer(ctx, mux, svc); err != nil {
		return err
	}
	srv.HTTP.PathPrefix("/").Handler(mux)

	return srv.Run()
}

func handleHealthCheck(srv healthpb.HealthServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := srv.Check(context.Background(), new(healthpb.HealthCheckRequest))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if resp.Status == healthpb.HealthCheckResponse_NOT_SERVING {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
	}
}
