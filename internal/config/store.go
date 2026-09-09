package config

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

//go:embed defaults
var defaults embed.FS

const reloadDelay = 150 * time.Millisecond

type State struct {
	Page     *Page    `json:"page"`
	Problems Problems `json:"problems,omitempty"`
}

type Store struct {
	dir string

	mu    sync.RWMutex
	state State
	subs  map[chan struct{}]struct{}
}

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	for _, name := range files {
		createStarter(dir, name)
	}
	empty, _ := parse(nil, nil)
	s := &Store{
		dir:   dir,
		state: State{Page: empty},
		subs:  map[chan struct{}]struct{}{},
	}
	s.reload()
	return s, nil
}

func createStarter(dir, name string) {
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		return
	}
	data, _ := defaults.ReadFile("defaults/" + name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		slog.Warn("could not create starter file", "path", path, "err", err)
		return
	}
	slog.Info("created starter file", "path", path)
}

func (s *Store) Current() State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

func (s *Store) Subscribe() (changed <-chan struct{}, unsubscribe func()) {
	ch := make(chan struct{}, 1)
	s.mu.Lock()
	s.subs[ch] = struct{}{}
	s.mu.Unlock()
	return ch, func() {
		s.mu.Lock()
		delete(s.subs, ch)
		s.mu.Unlock()
	}
}

func (s *Store) Watch(ctx context.Context) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()
	if err := w.Add(s.dir); err != nil {
		return err
	}

	timer := time.NewTimer(reloadDelay)
	timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case ev := <-w.Events:
			if slices.Contains(files, filepath.Base(ev.Name)) {
				timer.Reset(reloadDelay)
			}
		case err := <-w.Errors:
			slog.Warn("config watcher", "err", err)
		case <-timer.C:
			s.reload()
		}
	}
}

func (s *Store) reload() {
	var problems Problems
	read := func(name string) []byte {
		data, err := os.ReadFile(filepath.Join(s.dir, name))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			slog.Error("could not read config file", "name", name, "err", err)
			problems = append(problems, Problem{File: name, Message: "could not be read, check the file permissions"})
		}
		return data
	}
	settingsData, servicesData := read(settingsFile), read(servicesFile)

	var page *Page
	if len(problems) == 0 {
		var err error
		if page, err = parse(settingsData, servicesData); err != nil {
			errors.As(err, &problems)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if len(problems) > 0 {
		slog.Error("invalid config, keeping the last valid version", "dir", s.dir, "count", len(problems))
		for _, p := range problems {
			slog.Error("config problem", "file", p.File, "line", p.Line, "message", p.Message)
		}
		s.state.Problems = problems
	} else {
		slog.Info("config loaded", "dir", s.dir)
		s.state = State{Page: page}
	}
	for ch := range s.subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
