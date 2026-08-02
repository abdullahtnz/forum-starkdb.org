package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
	JWTRefresh  string
	AppPort     string
	UploadDir   string
	MaxFileSize int64
	FrontendURL string
}

func Load() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", ""),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", ""),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTRefresh:  getEnv("JWT_REFRESH_SECRET", ""),
		AppPort:     getEnv("APP_PORT", "8080"),
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
		MaxFileSize: getEnvInt("MAX_FILE_SIZE", 10485760),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func (c *Config) DSN() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=" + c.DBSSLMode
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return fallback
}
