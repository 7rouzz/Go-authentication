package services

import (
    "context"
    "errors"
    "time"
    
    "ecommerce-backend/internal/models"
    "ecommerce-backend/internal/repositories"
    "ecommerce-backend/internal/utils"
)

type AuthService struct {
    userRepo    *repositories.UserRepository
    jwtSecret   string
    jwtExpiry   time.Duration
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string, jwtExpiry string) *AuthService {
    expiry, _ := time.ParseDuration(jwtExpiry)
    return &AuthService{
        userRepo:  userRepo,
        jwtSecret: jwtSecret,
        jwtExpiry: expiry,
    }
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
    if err := utils.ValidateStruct(req); err != nil {
        return nil, err
    }
    
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
        Role:     "user",
    }
    
    createdUser, err := s.userRepo.CreateUser(ctx, user)
    if err != nil {
        return nil, err
    }
    
    createdUser.Password = ""
    return createdUser, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
    if err := utils.ValidateStruct(req); err != nil {
        return nil, err
    }
    
    user, err := s.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("invalid email or password")
    }
    
    if !utils.CheckPasswordHash(req.Password, user.Password) {
        return nil, errors.New("invalid email or password")
    }
    
    token, err := utils.GenerateJWT(
        user.ID.Hex(),
        user.Email,
        user.Role,
        s.jwtSecret,
        s.jwtExpiry,
    )
    if err != nil {
        return nil, err
    }
    
    user.Password = ""
    
    return &models.LoginResponse{
        Token: token,
        User:  *user,
    }, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    user.Password = ""
    return user, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
    if err := utils.ValidateStruct(req); err != nil {
        return nil, err
    }
    
    updatedUser, err := s.userRepo.UpdateUser(ctx, userID, req)
    if err != nil {
        return nil, err
    }
    
    updatedUser.Password = ""
    return updatedUser, nil
}
