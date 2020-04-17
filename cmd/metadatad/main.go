package main

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/pseidemann/finish"
	log "github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
)

const serviceName = "metadata"

type Config struct {
	LogLevel string `default:"info"`
	DBURI    string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local"`
	Addr     string `default:"0.0.0.0:8080"`
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
	ds, err := datastore.Open(cfg.DBURI)
	if err != nil {
		return err
	}
	defer ds.Close()

	handler := &server{ds: ds, router: httprouter.New()}
	handler.routes()

	srv := &http.Server{Addr: cfg.Addr, Handler: handler}

	fin := &finish.Finisher{
		Timeout: 15 * time.Second,
		Log:     log.StandardLogger(),
		Signals: finish.DefaultSignals,
	}
	fin.Add(srv)

	go func() {
		log.Infof("Starting server at %s", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fin.Wait()

	return nil
}
