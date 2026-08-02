package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abdullahtnz/forum-starkdb/internal/config"
	"github.com/abdullahtnz/forum-starkdb/internal/database"
	"github.com/abdullahtnz/forum-starkdb/internal/handlers"
	"github.com/abdullahtnz/forum-starkdb/internal/middleware"
	"github.com/abdullahtnz/forum-starkdb/internal/services"
)

func main() {
	cfg := config.Load()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL")

	authSvc := services.NewAuthService(pool, cfg.JWTSecret, cfg.JWTRefresh)
	postSvc := services.NewPostService(pool)
	commentSvc := services.NewCommentService(pool)
	userSvc := services.NewUserService(pool)
	adminSvc := services.NewAdminService(pool)
	tagSvc := services.NewTagService(pool)
	uploadSvc := services.NewUploadService(pool, cfg.UploadDir, cfg.MaxFileSize)

	authH := handlers.NewAuthHandler(authSvc)
	postH := handlers.NewPostHandler(postSvc)
	commentH := handlers.NewCommentHandler(commentSvc)
	userH := handlers.NewUserHandler(userSvc)
	adminH := handlers.NewAdminHandler(adminSvc)
	tagH := handlers.NewTagHandler(tagSvc)
	uploadH := handlers.NewUploadHandler(uploadSvc)

	rateLimiter := middleware.NewRateLimiter(100, 200)
	strictLimiter := middleware.NewRateLimiter(1, 5)

	mux := http.NewServeMux()

	// File server for uploads
	fs := http.FileServer(http.Dir(cfg.UploadDir))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", fs))

	// Auth routes
	mux.HandleFunc("POST /api/auth/signup", strictLimiter.Limit(http.HandlerFunc(authH.Signup)).ServeHTTP)
	mux.HandleFunc("POST /api/auth/login", strictLimiter.Limit(http.HandlerFunc(authH.Login)).ServeHTTP)
	mux.HandleFunc("POST /api/auth/refresh", http.HandlerFunc(authH.Refresh))
	mux.HandleFunc("POST /api/auth/logout", http.HandlerFunc(authH.Logout))

	// User routes
	mux.Handle("GET /api/users/me", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(userH.Me)))
	mux.Handle("PUT /api/users/me", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(userH.UpdateProfile)))
	mux.HandleFunc("GET /api/users/{id}", userH.PublicProfile)
	mux.HandleFunc("GET /api/users/{id}/posts", userH.UserPosts)

	// Post routes
	mux.HandleFunc("GET /api/posts", postH.List)
	mux.Handle("GET /api/posts/{id}", middleware.OptionalAuth(cfg.JWTSecret)(http.HandlerFunc(postH.Get)))
	mux.Handle("POST /api/posts", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(postH.Create)))
	mux.Handle("PUT /api/posts/{id}", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(postH.Update)))
	mux.Handle("DELETE /api/posts/{id}", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(postH.Delete)))
	mux.Handle("POST /api/posts/{id}/like", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(postH.ToggleLike)))
	mux.Handle("PUT /api/posts/{id}/pin", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(postH.PinPost))))
	mux.Handle("PUT /api/posts/{id}/close", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(postH.ClosePost))))

	// Comment routes
	mux.Handle("GET /api/posts/{id}/comments", middleware.OptionalAuth(cfg.JWTSecret)(http.HandlerFunc(commentH.List)))
	mux.Handle("POST /api/posts/{id}/comments", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(commentH.Create)))
	mux.Handle("POST /api/comments/{id}/reply", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(commentH.Reply)))
	mux.Handle("PUT /api/comments/{id}", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(commentH.Update)))
	mux.Handle("DELETE /api/comments/{id}", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(commentH.Delete)))
	mux.Handle("POST /api/comments/{id}/like", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(commentH.ToggleLike)))

	// Tag routes
	mux.HandleFunc("GET /api/tags", tagH.List)

	// Upload routes
	mux.Handle("POST /api/upload", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(uploadH.Upload)))

	// Admin routes
	mux.Handle("GET /api/admin/users", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(adminH.GetUsers))))
	mux.Handle("DELETE /api/admin/users/{id}", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(adminH.DeleteUser))))
	mux.Handle("PUT /api/admin/users/{id}/admin", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(adminH.ToggleAdmin))))
	mux.Handle("GET /api/admin/bad-words", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(adminH.GetBadWords))))
	mux.Handle("POST /api/admin/bad-words", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(adminH.AddBadWord))))
	mux.Handle("DELETE /api/admin/bad-words", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(adminH.RemoveBadWord))))
	mux.Handle("DELETE /api/admin/posts/{id}", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(adminH.DeletePost))))
	mux.Handle("DELETE /api/admin/comments/{id}", middleware.Auth(cfg.JWTSecret)(middleware.AdminOnly()(http.HandlerFunc(adminH.DeleteComment))))

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"starkdb-forum"}`))
	})

	handler := middleware.Logger(
		middleware.SecurityHeaders(
			middleware.CORS(cfg.FrontendURL)(
				rateLimiter.Limit(mux),
			),
		),
	)

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
