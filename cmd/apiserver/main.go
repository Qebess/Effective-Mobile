package main

import (
	"apiserver/internal/apiserver/config"
	"apiserver/internal/apiserver/handlers"
	"apiserver/store"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()
	loger := SetupLogger(config.AppEnv)

	loger.Info("Logger start...")
	loger.Debug("Debug message enabled")

	store := store.New()

	if err := store.Open(config.ConnStr); err != nil {
		loger.Error("Fail to open db:", err)
		os.Exit(1)
	}

	loger.Info("Connected to db")

	router := gin.Default()
	router.POST("/", handlers.CreateHandler(loger, store))
	router.GET("/:id", handlers.GetHandler(loger, store))
	router.PUT("/:id", handlers.UpdateHandler(loger, store))
	router.DELETE("/:id", handlers.DeleteHandler(loger, store))
	router.GET("/list", handlers.GetList(loger, store))
	router.POST("/range", handlers.GetRangeHandler(loger, store))

	connStr := ":" + config.AppPort

	router.Run(connStr)

}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
