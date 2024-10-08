package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DSorbon/effective-mobile-task/internal/app/server"
	"github.com/DSorbon/effective-mobile-task/internal/config"
	"github.com/DSorbon/effective-mobile-task/internal/repository/postgres"
	"github.com/DSorbon/effective-mobile-task/internal/service"
	"github.com/DSorbon/effective-mobile-task/internal/transport/http/handler"
	"github.com/DSorbon/effective-mobile-task/pkg/logger"
	"github.com/gookit/validate"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title Song API
// @version 1.0
// @description REST API for Song App

// @host localhost:8080
// @BasePath /api/v1/
func Run() {
	if err := config.LoadFromFile(".env"); err != nil {
		log.Fatalf("config.LoadFromFile(): %v", err)
	}

	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})

	pool, err := pgxpool.New(context.Background(), config.Values.PGDSN)
	if err != nil {
		log.Fatalf("couldn't connect to the database: %v", err)
	}
	log.Printf("successfully connected to db")

	logger.Init(zapCore(zapAtomicLevel()))

	songRepository := postgres.NewSongRepository(pool)
	songService := service.NewSongService(songRepository)

	newHandler := handler.NewHandler(songService)
	handlers := newHandler.InitRoutes()

	srv := server.NewServer()
	go srv.Run(handlers)

	log.Printf("Server started at %s port", config.Values.APIPort)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to stop server: %v", err)
	}

	pool.Close()
}

func zapCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	return zapcore.NewCore(consoleEncoder, stdout, level)
}

func zapAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(config.Values.LOG_LEVEL); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
