package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GolangWorld/db"
	"GolangWorld/handlers"
	"GolangWorld/routes"
	"GolangWorld/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables (Single Responsibility Principle - SRP)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Get server port from environment (SRP - Configuration handling should be separate)
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// Create new Echo instance (SRP - HTTP server initialization)
	e := echo.New()

	// Connect to MongoDB (SRP - Database connections should be handled separately)
	client, err := db.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.DisconnectMongoDB(client) // Ensure disconnection (Resource management)

	// Dependency Injection (DIP - Inject dependencies rather than hardcoding)
	userService := services.NewMongoUserService(client)
	userHandler := handlers.NewUserHandler(userService)

	// Group API routes (SRP - Routes should be managed separately)
	apiGroup := e.Group("/api")
	routes.RegisterRoutes(apiGroup, userHandler)

	// Graceful shutdown implementation (SRP - Server lifecycle management)
	go func() {
		if err := e.Start(":" + serverPort); err != nil && err != http.ErrServerClosed {
			log.Fatal("Shutting down server: ", err)
		}
	}()

	// Setup OS signal listening for SIGINT and SIGTERM (Handles system termination)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Blocks execution until a termination signal is received

	log.Println("Shutting down gracefully...")

	// Context with timeout ensures that all background processes finish before exiting
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to shut down the Echo server gracefully
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exited")
}
