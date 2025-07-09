package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"text_sharing/internal/config"
	"text_sharing/internal/http-server/handlers/auth/login"
	"text_sharing/internal/http-server/handlers/auth/register"
	"text_sharing/internal/http-server/handlers/texts/get"
	"text_sharing/internal/http-server/handlers/texts/save"
	mwJWT "text_sharing/internal/http-server/middleware/jwt"
	mwLogger "text_sharing/internal/http-server/middleware/logger"
	"text_sharing/internal/lib/logger/handlers/slogpretty"
	"text_sharing/internal/lib/logger/sl"
	"text_sharing/internal/storage/minio"
	"text_sharing/internal/storage/redis"
	"text_sharing/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting pastebin", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")

	redisClient, err := redis.NewClientForCache(ctx)
	if err != nil {
		log.Error("failed to init redis cache", sl.Err(err))
		os.Exit(1)
	}

	cache := redis.NewCache(redisClient)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	minioClient := minio.NewMinioClient()
	err = minioClient.InitMinio()
	if err != nil {
		log.Error("failed to init minio storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Group(func(r chi.Router) {
		r.Post("/pastebin/register", register.New(log, storage))
		r.Post("/pastebin/login", login.New(log, storage))
	})

	router.Group(func(r chi.Router) {
		r.Use(mwJWT.AuthorizeJWTToken)

		r.Post("/pastebin/write", save.New(log, storage, minioClient, cache))
		r.Get("/pastebin/{linkID}", get.New(log, storage, minioClient, cache))
	})

	log.Info("starting pastebin", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
