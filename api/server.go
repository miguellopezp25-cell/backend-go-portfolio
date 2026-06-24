package api

import (
	"context"
	"fmt"
	"log"
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
	cfg    *config.Config
	pool   *pgxpool.Pool
	store  db.Store
	svc    *visitorservice.Service
	http   *http.Server
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
	router := SetupRouter(svc)

	return &Server{
		cfg:   cfg,
		pool:  pool,
		store: store,
		svc:   svc,
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: router,
		},
	}, nil
}

func (s *Server) Start() error {
	go func() {
		log.Printf("Server listening on port %d", s.cfg.Server.Port)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return err
	}

	s.pool.Close()
	log.Println("Server exited gracefully")
	return nil
}
