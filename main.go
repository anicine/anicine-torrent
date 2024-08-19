package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anicine/anicine-torrent/client"
	"github.com/anicine/anicine-torrent/internal/config"
	"github.com/anicine/anicine-torrent/nyaa"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/sync/errgroup"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	config, err := config.Load(".env")
	if err != nil {
		slog.Error("cannot load the env variables", "error", err)
		os.Exit(1)
	}

	uri, err := url.Parse(config.Proxy)
	if err != nil {
		slog.Error("cannot parse the proxy URI", "error", err)
		os.Exit(1)
	}

	proxy := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
		Timeout: time.Duration(30 * time.Second),
	}

	client.SetProxy(proxy)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDB))
	if err != nil {
		cancel()
		slog.Error("cannot create the connection with the database", "error", err)
		os.Exit(1)
	}

	if err = db.Ping(ctx, readpref.Primary()); err != nil {
		cancel()
		slog.Error("cannot ping the database", "error", err)
		os.Exit(1)
	}

	conf := torrent.NewDefaultClientConfig()
	conf.DisableAggressiveUpload = true
	conf.DisableWebseeds = true
	conf.DisableIPv6 = true
	conf.DisableTCP = true
	conf.NoUpload = true
	conf.NoDHT = true
	conf.Seed = false

	bitTorrent, err := torrent.NewClient(conf)
	if err != nil {
		slog.Error("cannot start to the bittorrent client", "error", err)
		os.Exit(1)
	}

	ns := nyaa.NewNyaa(slog.Default().WithGroup("[NYAA]"), db.Database("anicine-torrent"), bitTorrent)

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	if err = ns.Init(ctx); err != nil {
		cancel()
		slog.Error("cannot migrate the database", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":5555",
		Handler: mux,
		TLSConfig: &tls.Config{
			ClientAuth: tls.VerifyClientCertIfGiven,
			MinVersion: tls.VersionTLS13,
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
		},
	}

	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		slog.Warn("shutdown the server ...")
		cancel()
	}()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return ns.Start(ctx)
	})
	group.Go(func() error {
		if config.CertFile != "" || config.KeyFile != "" {
			return server.ListenAndServeTLS(config.CertFile, config.KeyFile)
		}
		return server.ListenAndServe()
	})
	group.Go(func() error {
		<-ctx.Done()
		slog.Warn("try to stop the api ...")
		return server.Shutdown(context.Background())
	})
	group.Go(func() error {
		<-ctx.Done()
		slog.Warn("try to stop the db ...")
		ns.Save(context.Background())
		return db.Disconnect(ctx)
	})
	group.Go(func() error {
		<-ctx.Done()
		slog.Warn("try to stop the bittorrent client ...")
		if errs := bitTorrent.Close(); len(errs) > 0 {
			return errs[0]
		}
		return nil
	})

	slog.Info("server is running", "address", server.Addr)
	if err := group.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Warn("graceful shutdown complete.")
}
