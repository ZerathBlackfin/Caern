package main

import (
	"cmp"
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"caern/internal/config"
	"caern/internal/server"
	"caern/web"
)

func main() {
	addr := cmp.Or(os.Getenv("CAERN_ADDR"), ":7676")
	dir := cmp.Or(os.Getenv("CAERN_CONFIG_DIR"), "config")
	if err := run(addr, dir); err != nil {
		slog.Error("caern stopped", "err", err)
		os.Exit(1)
	}
}

func run(addr, configDir string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	store, err := config.NewStore(configDir)
	if err != nil {
		return err
	}
	go func() {
		if err := store.Watch(ctx); err != nil {
			slog.Error("config changes will not be picked up", "err", err)
		}
	}()

	srv := &http.Server{
		Addr:              addr,
		Handler:           server.New(store, configDir, web.Files),
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext:       func(net.Listener) context.Context { return ctx },
	}
	errc := make(chan error, 1)
	go func() {
		slog.Info("caern is listening", "addr", addr, "config", configDir)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
