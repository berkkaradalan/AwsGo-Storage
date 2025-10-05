package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berkkaradalan/AwsGo-Storage/config"
	"github.com/berkkaradalan/AwsGo-Storage/handlers"
	"github.com/berkkaradalan/AwsGo-Storage/repositories"
	"github.com/berkkaradalan/AwsGo-Storage/routers"
	"github.com/berkkaradalan/AwsGo-Storage/services"
)

func main() {
	dbService := config.ConnectDatabase()
	log.Println("Database connected successfully")

	env := config.LoadEnv()

	authConfig := config.NewAuthConfig(*env)

	userRepo := repositories.NewUserRepository(dbService)
	userService := services.NewUserService(userRepo, authConfig)
	userHandler := handlers.NewUserHandler(userService)


	router := routers.SetupRouter(userHandler, *env, authConfig)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
	
	_ = dbService
}