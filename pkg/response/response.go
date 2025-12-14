package response

import (
    "github.com/gin-gonic/gin"
)

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
    c.JSON(status, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func Error(c *gin.Context, status int, message string) {
    c.JSON(status, Response{
        Success: false,
        Error:   message,
    })
    c.Abort()
}
