package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"hackathon-absolutecinema/internal/ai_client" // <-- ҚОСЫЛДЫ
	"hackathon-absolutecinema/internal/config"
	"hackathon-absolutecinema/internal/handler"
	"hackathon-absolutecinema/internal/middleware"
	"hackathon-absolutecinema/internal/repository"
	"hackathon-absolutecinema/internal/service"
	"hackathon-absolutecinema/pkg/logger"
)

func main() {
	// Logger инициализация
	log := logger.NewLogger()

	// .env жүктеу
	_ = godotenv.Load()

	// Конфигурация
	cfg := config.LoadConfig()

	// PostgreSQL қосылу
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database-ке қосылу мүмкін емес: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database тексеру қатесі: %v", err)
	}
	log.Info("✅ PostgreSQL-ге сәтті қосылдық")

	// Repository инициализация
	userRepo := repository.NewUserRepository(pool)
	companyRepo := repository.NewCompanyRepository(pool)
	analysisRepo := repository.NewAnalysisRepository(pool)

	// AI Клиентін инициализациялау <-- ЖАҢА ҚОСЫЛДЫ
	aiClient := ai_client.NewClient(cfg.AIServiceURL)

	// Service инициализация <-- ТҮЗЕТІЛДІ
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	companyService := service.NewCompanyService(companyRepo)
	// Енді analysisService-ке companyRepo және aiClient беріледі
	analysisService := service.NewAnalysisService(analysisRepo, companyRepo, aiClient)

	// Handler инициализация
	authHandler := handler.NewAuthHandler(authService)
	companyHandler := handler.NewCompanyHandler(companyService)
	analysisHandler := handler.NewAnalysisHandler(analysisService)
	healthHandler := handler.NewHealthHandler(pool)

	// Router (Go 1.22+ method-based routing)
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /health", healthHandler.HealthCheck)
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// Protected routes
	mux.HandleFunc("GET /api/companies", middleware.AuthMiddleware(companyHandler.GetCompanies, cfg.JWTSecret))
	mux.HandleFunc("POST /api/companies", middleware.AuthMiddleware(companyHandler.CreateCompany, cfg.JWTSecret))
	mux.HandleFunc("GET /api/companies/{id}", middleware.AuthMiddleware(companyHandler.GetCompany, cfg.JWTSecret))
	mux.HandleFunc("PUT /api/companies/{id}", middleware.AuthMiddleware(companyHandler.UpdateCompany, cfg.JWTSecret))
	mux.HandleFunc("DELETE /api/companies/{id}", middleware.AuthMiddleware(companyHandler.DeleteCompany, cfg.JWTSecret))

	mux.HandleFunc("POST /api/analyses", middleware.AuthMiddleware(analysisHandler.CreateAnalysis, cfg.JWTSecret))
	mux.HandleFunc("GET /api/analyses/{id}", middleware.AuthMiddleware(analysisHandler.GetAnalysis, cfg.JWTSecret))

	// Middleware қолдану (CORS + Logger)
	finalHandler := middleware.LoggerMiddleware(middleware.CORSMiddleware(mux), log)

	// HTTP Server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      finalHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Серверді іске қосу
	go func() {
		log.Infof("🚀 Сервер %s портта іске қосылды", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Сервер қатесі: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("⏹️ Сервер сөндірілуде...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Серверді сөндіру қатесі: %v", err)
	}

	log.Info("✅ Сервер сәтті сөндірілді")
}
