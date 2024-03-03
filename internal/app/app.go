package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	conf "github.com/Mth-Ryan/go-yaml-cfg"
	calendarhandler "github.com/distributed-calendar/calendar-server/internal/handler/calendar"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	calendarservice "github.com/distributed-calendar/calendar-server/internal/service/calendar"
)

type errFunc func() error

type App struct {
	cfg *Config
	mux *chi.Mux

	onCleanupFuncs []errFunc
	onRunFuncs     []errFunc

	pgConnPool *pgxpool.Pool

	calendarService *calendarservice.Service
}

func (a *App) Run() {
	defer a.cleanup()

	slog.Info("starting server...")

	if err := a.run(); err != nil {
		panic(err)
	}

	slog.Info("server stopped, cleaning up...")
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
			slog.Error("on cleanup error", err)
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

	a.initMux()

	if err := a.initPostgres(); err != nil {
		return err
	}

	a.initHttpServer()

	a.initCalendarService()

	a.initCalendarHandler()

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

func (a *App) initMux() {
	a.mux = chi.NewMux()

	options := httplog.Options{
		LogLevel: slog.LevelDebug,
		JSON:     true,
	}

	a.mux.Use(
		httplog.RequestLogger(
			httplog.NewLogger("calendar-server", options),
		),
		middleware.Recoverer,
	)

	a.mux.Mount("/ping", a.pingHandler())
}

func (a *App) initPostgres() error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		a.cfg.Postgres.User,
		a.cfg.Postgres.Password,
		a.cfg.Postgres.Host,
		a.cfg.Postgres.Port,
		a.cfg.Postgres.Dbname,
	)

	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("cannot create pgx pool: %w", err)
	}

	conn.Config().MaxConns = 5

	a.pgConnPool = conn

	a.addOnCleanup(func() error {
		conn.Close()

		return nil
	})

	return nil
}

func (a *App) initCalendarHandler() {
	handler := calendarhandler.NewHandler(a.calendarService)
	a.mux.Mount("/", handler)
}

func (a *App) initHttpServer() {
	server := &http.Server{
		Addr:    ":" + a.cfg.HttpServer.Port,
		Handler: a.mux,
	}

	a.addOnRun(func() error {
		slog.Info("server started")

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
			slog.Error("cannot write to ping", err)
		}
	})
}

func newConfig(configPath string) (*Config, error) {
	if err := conf.InitializeConfigSingleton[Config](configPath); err != nil {
		return nil, err
	}

	cfg, err := conf.GetConfigFromSingleton[Config]()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewApp(configPath string) (*App, error) {
	app := &App{}

	err := app.init(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot init app: %w", err)
	}

	return app, nil
}
