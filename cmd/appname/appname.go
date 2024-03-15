package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/idan-fishman/go-echo-starter/cmd/appname/config"
	"github.com/idan-fishman/go-echo-starter/pkg/log"
	"github.com/idan-fishman/go-echo-starter/pkg/middleware"
	"github.com/idan-fishman/go-echo-starter/pkg/probe"
	"github.com/idan-fishman/go-echo-starter/pkg/validation"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Initialize loggers.
	logLevel := zap.NewAtomicLevelAt(zap.FatalLevel)
	serverLogger := log.InitializeLogger(logLevel, true, true)
	requestLogger := log.InitializeLogger(logLevel, true, true)
	defer log.FlushLogger(serverLogger)
	defer log.FlushLogger(requestLogger)
	serverLogger.Info("start-up process initiated.")

	// Initialize validator
	v := validation.NewValidator()

	// Load application configuration and set the log level.
	cfg, err := config.Load(v)
	if err != nil {
		serverLogger.Fatal("failed to load application configuration, unable to proceed", zap.Error(err))
	}
	log.UpdateLogLevel(logLevel, cfg.Log.Level)

	// Initialize Redis client and health checker
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	rdbHealthChecker := &probe.RedisHealthChecker{Client: rdb}

	// Initialize HTTP server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.RequestLogger(requestLogger))
	e.Use(middleware.RequestTimeout(cfg.Server.RequestTimeoutSeconds))

	// Define probes routes
	e.GET("/healthz", probe.LivenessProbe)
	e.GET("/readyz", probe.ReadinessProbe(serverLogger, rdbHealthChecker))

	// Start HTTP server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		serverLogger.Info(fmt.Sprintf("starting http server. listening on port %d", cfg.Server.Port))
		if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
			serverLogger.Fatal("http server start-up failed", zap.Error(err))
		}
	}()

	// Wait for shutdown signal and initiate graceful shutdown
	<-ctx.Done()

	// Shutdown HTTP server
	serverLogger.Info("received shutdown signal. initiating graceful shutdown.")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.GracefulShutdownTimeoutSeconds)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		serverLogger.Fatal("error during http server shutdown. some connections may not have been properly closed.", zap.Error(err))
	} else {
		serverLogger.Info("http server has been gracefully shutdown. all connections closed.")
	}
}
