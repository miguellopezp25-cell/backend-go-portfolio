package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/miguel/go-back-portfolo/config"
	"github.com/miguel/go-back-portfolo/database"
	"github.com/miguel/go-back-portfolo/schema/db"
	"github.com/miguel/go-back-portfolo/service/visitorservice"
)

type Server struct {
	cfg            *config.Config
	pool           *pgxpool.Pool
	store          db.Store
	visitorService *visitorservice.Service
	http           *http.Server
}

func NewServer(cfgPath string) (*Server, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	store := db.NewStore(pool)
	svc := visitorservice.NewService(store)

	s := &Server{
		cfg:            cfg,
		pool:           pool,
		store:          store,
		visitorService: svc,
	}

	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: SetupRouter(s),
	}

	return s, nil
}

func (s *Server) Start() error {
	go func() {
		slog.Info("server listening", "port", s.cfg.Server.Port)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return err
	}

	s.pool.Close()
	slog.Info("server exited gracefully")
	return nil
}
