package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    MongoDBURI  string
    DBName      string
    JWTSecret   string
    JWTExpiry   string
}

func LoadConfig() *Config {
    _ = godotenv.Load()
    
    return &Config{
        Port:       getEnv("PORT"),
        MongoDBURI: getEnv("MONGODB_URI"),
        DBName:     getEnv("DB_NAME"),
        JWTSecret:  getEnv("JWT_SECRET"),
        JWTExpiry:  getEnv("JWT_EXPIRY"),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
