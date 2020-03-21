package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	logz "github.com/videocoin/common/log"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
	"golang.org/x/sync/errgroup"
)

var (
	// ServiceName is the service name.
	ServiceName = "metadata"

	// Version is the application version.
	Version = "dev"
)

// Config ...
type Config struct {
	LogLevel string `default:"info"`
	DBURI    string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local"`
	Addr     string `default:"0.0.0.0:8080"`
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
	logz.SetGlobal(logger)

	ds, err := datastore.Open(cfg.DBURI)
	if err != nil {
		return err
	}
	defer ds.Close()

	srv := newServer()
	srv.ds = ds
	srv.router = httprouter.New()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	ctx, cancel := context.WithCancel(context.Background())

	errgrp, ctx := errgroup.WithContext(ctx)
	var httpSrv *http.Server
	errgrp.Go(func() error {
		httpSrv = &http.Server{Addr: cfg.Addr, Handler: srv}
		return httpSrv.ListenAndServe()
	})

	select {
	case <-stop:
		break
	case <-ctx.Done():
		break
	}

	cancel()

	if httpSrv != nil {
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpSrv.Shutdown(ctx)
	}

	if err := errgrp.Wait(); err != nil {
		return err
	}

	return nil
}
