package api

import (
    "ecommerce-backend/internal/handlers"
    "ecommerce-backend/internal/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, jwtSecret string) {
    public := router.Group("/api/v1/auth")
    {
        public.POST("/register", authHandler.Register)
        public.POST("/login", authHandler.Login)
    }
    
    protected := router.Group("/api/v1")
    protected.Use(middleware.AuthMiddleware(jwtSecret))
    {
        auth := protected.Group("/auth")
        {
            auth.GET("/profile", authHandler.GetProfile)
            auth.PUT("/profile", authHandler.UpdateProfile)
        }
    }
    
    admin := protected.Group("/admin")
    admin.Use(middleware.AdminMiddleware())
    {
    }
}
