package main

import (
	"log/slog"
	"net/http"
	"os"
	"text_sharing/internal/config"
	"text_sharing/internal/lib/logger/handlers/slogpretty"
	"text_sharing/internal/lib/logger/sl"
	"text_sharing/internal/storage/minio"

	"github.com/go-chi/chi/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	minioClient := minio.NewMinioClient()
	err := minioClient.InitMinio()
	if err != nil {
		log.Error("failed to init minio storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	log.Info("starting minio", slog.String("address", cfg.CfgMinio.MinioEndpoint))

	srv := &http.Server{
		Addr:         cfg.CfgMinio.MinioEndpoint,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
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
