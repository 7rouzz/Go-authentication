package handlers

import (
    "ecommerce-backend/internal/models"
    "ecommerce-backend/internal/services"
    "ecommerce-backend/pkg/response"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "Invalid request payload")
        return
    }
    
    user, err := h.authService.Register(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, 400, err.Error())
        return
    }
    
    response.Success(c, 201, "User registered successfully", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "Invalid request payload")
        return
    }
    
    loginResponse, err := h.authService.Login(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, 401, err.Error())
        return
    }
    
    response.Success(c, 200, "Login successful", loginResponse)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    user, err := h.authService.GetProfile(c.Request.Context(), userID.(string))
    if err != nil {
        response.Error(c, 404, err.Error())
        return
    }
    
    response.Success(c, 200, "Profile retrieved successfully", user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    var req models.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, "Invalid request payload")
        return
    }
    
    user, err := h.authService.UpdateProfile(c.Request.Context(), userID.(string), &req)
    if err != nil {
        response.Error(c, 400, err.Error())
        return
    }
    
    response.Success(c, 200, "Profile updated successfully", user)
}
