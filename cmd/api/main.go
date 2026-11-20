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
	log := logger.NewLogger()
	_ = godotenv.Load()
	cfg := config.LoadConfig()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database-ке қосылу мүмкін емес: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database тексеру қатесі: %v", err)
	}
	log.Info("✅ Connected to PostgreSQL successfully")

	userRepo := repository.NewUserRepository(pool)
	companyRepo := repository.NewCompanyRepository(pool)
	analysisRepo := repository.NewAnalysisRepository(pool)

	aiClient := ai_client.NewClient(cfg.AIServiceURL)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	companyService := service.NewCompanyService(companyRepo)

	analysisService := service.NewAnalysisService(analysisRepo, companyRepo, aiClient)

	authHandler := handler.NewAuthHandler(authService)
	companyHandler := handler.NewCompanyHandler(companyService)
	analysisHandler := handler.NewAnalysisHandler(analysisService)
	healthHandler := handler.NewHealthHandler(pool)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler.HealthCheck)
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	mux.HandleFunc("GET /api/companies", middleware.AuthMiddleware(companyHandler.GetCompanies, cfg.JWTSecret))
	mux.HandleFunc("POST /api/companies", middleware.AuthMiddleware(companyHandler.CreateCompany, cfg.JWTSecret))
	mux.HandleFunc("GET /api/companies/{id}", middleware.AuthMiddleware(companyHandler.GetCompany, cfg.JWTSecret))
	mux.HandleFunc("PUT /api/companies/{id}", middleware.AuthMiddleware(companyHandler.UpdateCompany, cfg.JWTSecret))
	mux.HandleFunc("DELETE /api/companies/{id}", middleware.AuthMiddleware(companyHandler.DeleteCompany, cfg.JWTSecret))

	mux.HandleFunc("POST /api/analyses", middleware.AuthMiddleware(analysisHandler.CreateAnalysis, cfg.JWTSecret))
	mux.HandleFunc("GET /api/analyses/{id}", middleware.AuthMiddleware(analysisHandler.GetAnalysis, cfg.JWTSecret))

	finalHandler := middleware.LoggerMiddleware(middleware.CORSMiddleware(mux), log)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      finalHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof("🚀 Server started on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Сервер қатесі: %v", err)
		}
	}()

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
