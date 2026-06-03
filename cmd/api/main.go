package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KakoeToImya/go-ws-chat/config"
	"github.com/KakoeToImya/go-ws-chat/internal/handler"
	"github.com/KakoeToImya/go-ws-chat/internal/repository/postgres"
	"github.com/KakoeToImya/go-ws-chat/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	//конфиг
	cfg := config.MustLoad()
	logger := log.New(os.Stdout, "API: ", log.LstdFlags)

	// БД
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())

	if err != nil {
		logger.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		logger.Fatalf("Unable to ping database: %v", err)
	}
	logger.Println("Connected to PostgreSQL")
	// слои
	userRepo := postgres.NewUserRepo(pool)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	//роутер
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	// сервер
	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Printf("server started in [%s] mode on port :%s", cfg.Env, cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server failure %v", err)
		}
	}()

	// graceful отключение
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("SD server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server forced to sd %v", err)
	}

	logger.Println("server stoped")

}
