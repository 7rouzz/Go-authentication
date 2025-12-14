package main

import (
    "log"
    "ecommerce-backend/internal/api"
    "ecommerce-backend/internal/config"
)

func main() {
    cfg := config.LoadConfig()
    server := api.NewServer(cfg)
    
    if err := server.Start(); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
