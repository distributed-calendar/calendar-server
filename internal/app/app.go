package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

type errFunc func() error

type App struct {
	cfg    *Config
	mux    *chi.Mux
	logger *zap.Logger

	onCleanupFuncs []errFunc
	onRunFuncs     []errFunc
}

func (a *App) Run() {
	defer a.cleanup()

	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	errGroup, _ := errgroup.WithContext(context.Background())

	for _, f := range a.onRunFuncs {
		errGroup.Go(f)
	}

	return errGroup.Wait()
}

func (a *App) addOnRun(f errFunc) {
	a.onRunFuncs = append(a.onRunFuncs, f)
}

func (a *App) cleanup() {
	for _, f := range a.onCleanupFuncs {
		err := f()
		if err != nil {
			a.logger.Error("on cleanup error", zap.Error(err))
		}
	}
}

func (a *App) addOnCleanup(f func() error) {
	a.onCleanupFuncs = append(a.onCleanupFuncs, f)
}

func (a *App) init(configPath string) error {
	if err := a.initConfig(configPath); err != nil {
		return err
	}

	a.initLogger()
	a.initMux()
	a.initHttpServer()

	return nil
}

func (a *App) initConfig(configPath string) error {
	cfg, err := newConfig(configPath)
	if err != nil {
		return fmt.Errorf("cannot create config: %w", err)
	}

	a.cfg = cfg

	return nil
}

func (a *App) initLogger() {
	logger := zap.Must(zap.NewProduction())

	a.addOnCleanup(logger.Sync)
}

func (a *App) initMux() {
	a.mux = chi.NewMux()
	a.mux.Mount("/ping", a.pingHandler())
}

func (a *App) initHttpServer() {
	server := &http.Server{
		Addr:    ":" + a.cfg.HttpServer.Port,
		Handler: a.mux,
	}

	a.addOnRun(func() error {
		e := server.ListenAndServe()
		if errors.Is(e, http.ErrServerClosed) {
			return nil
		}

		return e
	})
}

func (a *App) pingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("pong"))
		if err != nil {
			a.logger.Error("cannot write to ping", zap.Error(err))
		}
	})
}

func newConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("cannot decode config file: %w", err)
	}

	return cfg, nil
}

func NewApp(configPath string) (*App, error) {
	app := &App{}

	err := app.init(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot init app: %w", err)
	}

	return app, nil
}
