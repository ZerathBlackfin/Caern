package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"caern/internal/config"
)

const (
	heartbeat      = 30 * time.Second
	maxLayoutBytes = 1 << 20
)

func New(store *config.Store, configDir string, ui fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /api/events", events(store))
	mux.HandleFunc("POST /api/layout", saveLayout(store))
	mux.Handle("GET /icons/", userFiles("/icons/", filepath.Join(configDir, "icons")))
	mux.Handle("GET /images/", userFiles("/images/", filepath.Join(configDir, "images")))
	mux.Handle("GET /", uiFiles(ui))
	return mux
}

func saveLayout(store *config.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var layout config.Layout
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxLayoutBytes)).Decode(&layout); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := store.SaveLayout(layout); err != nil {
			slog.Error("could not save the layout", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func events(store *config.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		changed, unsubscribe := store.Subscribe()
		defer unsubscribe()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Accel-Buffering", "no")
		rc := http.NewResponseController(w)

		write := func(msg string) error {
			if _, err := io.WriteString(w, msg); err != nil {
				return err
			}
			return rc.Flush()
		}
		sendState := func() error {
			data, err := json.Marshal(store.Current())
			if err != nil {
				return err
			}
			return write("event: state\ndata: " + string(data) + "\n\n")
		}

		tick := time.NewTicker(heartbeat)
		defer tick.Stop()
		err := sendState()
		for err == nil {
			select {
			case <-r.Context().Done():
				return
			case <-changed:
				err = sendState()
			case <-tick.C:
				err = write(": ping\n\n")
			}
		}
	}
}

func userFiles(prefix, dir string) http.Handler {
	files := http.StripPrefix(prefix, http.FileServerFS(os.DirFS(dir)))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		files.ServeHTTP(w, r)
	})
}

func uiFiles(ui fs.FS) http.Handler {
	if _, err := fs.Stat(ui, "index.html"); errors.Is(err, fs.ErrNotExist) {
		slog.Warn("the interface is not built, run `make ui` first")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "The Caern interface is not built.", http.StatusServiceUnavailable)
		})
	}
	files := http.FileServerFS(ui)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		files.ServeHTTP(w, r)
	})
}
