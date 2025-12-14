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
        Port:       getEnv("PORT", "8080"),
        MongoDBURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
        DBName:     getEnv("DB_NAME", "ecommerce"),
        JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
        JWTExpiry:  getEnv("JWT_EXPIRY", "24h"),
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
