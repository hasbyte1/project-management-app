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

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/hasbyte1/project-management-app/internal/config"
	"github.com/hasbyte1/project-management-app/internal/handler"
	"github.com/hasbyte1/project-management-app/internal/middleware"
	"github.com/hasbyte1/project-management-app/internal/repository"
	"github.com/hasbyte1/project-management-app/internal/service"
	"github.com/hasbyte1/project-management-app/pkg/database"
	"github.com/hasbyte1/project-management-app/pkg/jwt"
	"github.com/hasbyte1/project-management-app/pkg/redis"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting server in %s mode on port %s", cfg.Server.Environment, cfg.Server.Port)

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL database")

	// Initialize Redis
	redisClient, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		// Continue without Redis
	} else {
		defer redisClient.Close()
		log.Println("Connected to Redis")
	}

	// Initialize JWT token manager
	tokenManager := jwt.NewTokenManager(&cfg.JWT)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenManager)
	orgService := service.NewOrganizationService(orgRepo, userRepo)
	projectService := service.NewProjectService(projectRepo, taskRepo)
	taskService := service.NewTaskService(taskRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	orgHandler := handler.NewOrganizationHandler(orgService)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", authHandler.Register)
			r.Post("/auth/login", authHandler.Login)
			r.Post("/auth/refresh", authHandler.RefreshToken)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Auth
			r.Get("/auth/me", authHandler.GetCurrentUser)

			// Organizations
			r.Post("/organizations", orgHandler.Create)
			r.Get("/organizations", orgHandler.List)
			r.Get("/organizations/{organizationId}", orgHandler.GetByID)
			r.Patch("/organizations/{organizationId}", orgHandler.Update)
			r.Delete("/organizations/{organizationId}", orgHandler.Delete)

			// Organization members
			r.Get("/organizations/{organizationId}/members", orgHandler.GetMembers)
			r.Post("/organizations/{organizationId}/members", orgHandler.AddMember)
			r.Patch("/organizations/{organizationId}/members/{memberId}", orgHandler.UpdateMemberRole)
			r.Delete("/organizations/{organizationId}/members/{memberId}", orgHandler.RemoveMember)

			// Projects
			r.Post("/projects", projectHandler.Create)
			r.Get("/organizations/{organizationId}/projects", projectHandler.List)
			r.Get("/projects/{projectId}", projectHandler.GetByID)
			r.Patch("/projects/{projectId}", projectHandler.Update)
			r.Delete("/projects/{projectId}", projectHandler.Delete)
			r.Post("/projects/{projectId}/archive", projectHandler.Archive)
			r.Post("/projects/{projectId}/unarchive", projectHandler.Unarchive)

			// Project members
			r.Get("/projects/{projectId}/members", projectHandler.GetMembers)
			r.Post("/projects/{projectId}/members", projectHandler.AddMember)
			r.Patch("/projects/{projectId}/members/{memberId}", projectHandler.UpdateMemberRole)
			r.Delete("/projects/{projectId}/members/{memberId}", projectHandler.RemoveMember)

			// Tasks
			r.Post("/tasks", taskHandler.Create)
			r.Get("/tasks/{taskId}", taskHandler.GetByID)
			r.Patch("/tasks/{taskId}", taskHandler.Update)
			r.Delete("/tasks/{taskId}", taskHandler.Delete)
			r.Patch("/tasks/{taskId}/status", taskHandler.UpdateStatus)

			// Task comments
			r.Get("/tasks/{taskId}/comments", taskHandler.GetComments)
			r.Post("/tasks/{taskId}/comments", taskHandler.CreateComment)

			// Project tasks and statuses
			r.Get("/projects/{projectId}/tasks", taskHandler.List)
			r.Get("/projects/{projectId}/statuses", taskHandler.GetStatuses)
			r.Post("/projects/{projectId}/statuses", taskHandler.CreateStatus)
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on http://localhost:%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
