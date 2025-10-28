package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dekanayake/kart-challenge/backend-challenge/internal"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/config"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/reader"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/repository"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/routes"
)

func main() {
	config.LoadConfig()
	config.InitLogger()

	config.Logger.Info().
		Interface("config", config.AppConfig).
		Msg("Loaded environment configuration")

	port := fmt.Sprintf(":%d", config.AppConfig.Port)
	config.Logger.Info().Msgf("Starting backend-challenge service on %s...", port)

	productRepo := repository.GetProductRepository()
	orderRepo := repository.GetOrderRepository()
	fileReader, err := reader.GetFileReader("hdd", config.AppConfig.CouponCodeFolderPath, config.AppConfig.CouponCodeFilePartialIndexChunkSize, config.AppConfig.CouponCodeFileConcurrentPoolSize)
	if err != nil {
		config.Logger.Error().Err(err).Msg("Failed in creating the File reader.")
	}

	server := internal.Server{
		ProductRepo: &productRepo,
		OrderRepo:   &orderRepo,
		FileReader:  &fileReader,
	}

	r := routes.SetupRouter(server)
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		config.Logger.Info().
			Str("port", strconv.Itoa(config.AppConfig.Port)).
			Msg("Server started successfully and listening for requests")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			config.Logger.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	config.Logger.Warn().Msg("Received termination signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		config.Logger.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		config.Logger.Info().Msg("Server shut down gracefully")
	}

}
