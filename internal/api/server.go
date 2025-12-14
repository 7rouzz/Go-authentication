package api

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
    
    "ecommerce-backend/internal/config"
    "ecommerce-backend/internal/database"
    "ecommerce-backend/internal/handlers"
    "ecommerce-backend/internal/repositories"
    "ecommerce-backend/internal/services"
    "github.com/gin-gonic/gin"
)

type Server struct {
    router *gin.Engine
    config *config.Config
    mongo  *database.MongoDB
}

func NewServer(cfg *config.Config) *Server {
    mongoDB, err := database.ConnectMongoDB(cfg.MongoDBURI, cfg.DBName)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    
    router := gin.Default()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    router.Use(CORSMiddleware())
    
    userRepo := repositories.NewUserRepository(mongoDB.DB)
    authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
    authHandler := handlers.NewAuthHandler(authService)
    
    SetupRoutes(router, authHandler, cfg.JWTSecret)
    
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "OK",
            "timestamp": time.Now().Unix(),
        })
    })
    
    return &Server{
        router: router,
        config: cfg,
        mongo:  mongoDB,
    }
}

func (s *Server) Start() error {
    srv := &http.Server{
        Addr:    ":" + s.config.Port,
        Handler: s.router,
    }
    
    go func() {
        log.Printf("Server starting on port %s", s.config.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server failed:", err)
        }
    }()
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := s.mongo.Disconnect(); err != nil {
        log.Println("Error disconnecting from MongoDB:", err)
    }
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exiting")
    return nil
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
