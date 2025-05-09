package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yourusername/bartenderapp/services/auth/internal/handlers"
	"github.com/yourusername/bartenderapp/services/auth/internal/repository"
	"github.com/yourusername/bartenderapp/services/auth/internal/service"
	"github.com/yourusername/bartenderapp/services/pkg/database"
	"github.com/yourusername/bartenderapp/services/pkg/middleware"
)

func main() {
	// Initialize database connection
	dbConfig := database.NewConfigFromEnv()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Create repositories
	userRepo := repository.NewUserRepository(db)

	// Create services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Create router
	router := mux.NewRouter()

	// Apply middleware to all routes
	router.Use(middleware.RequestLogger)
	router.Use(middleware.CORS)
	router.Use(middleware.JSONContentType)

	// Public routes
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	router.HandleFunc("/verify", authHandler.VerifyToken).Methods("POST")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		middleware.RespondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		})
	}).Methods("GET")

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// Protected routes - require authentication
	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.Authenticate)

	// User routes
	protected.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	protected.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	protected.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	protected.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	protected.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
	protected.HandleFunc("/users/me", userHandler.GetCurrentUser).Methods("GET")

	// Admin-only routes
	adminRouter := router.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.Authenticate)
	adminRouter.Use(middleware.RequireRole("admin"))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Run server in a goroutine so we can gracefully shut it down
	go func() {
		log.Printf("Auth Service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal
	<-c

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Shutdown the server
	srv.Shutdown(ctx)
	log.Println("Auth Service shutting down")
	os.Exit(0)
} 