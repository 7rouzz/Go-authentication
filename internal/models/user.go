package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name" validate:"required,min=2,max=100"`
    Email     string             `json:"email" bson:"email" validate:"required,email"`
    Password  string             `json:"password,omitempty" bson:"password" validate:"required,min=6"`
    Role      string             `json:"role" bson:"role" validate:"required,oneof=user admin"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
    Name  string `json:"name" validate:"omitempty,min=2,max=100"`
    Email string `json:"email" validate:"omitempty,email"`
}
