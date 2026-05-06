package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port            string
	JWTSecret       string
	HostFolder      string
	TemplatesFolder string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
}

func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "3000"),
		JWTSecret:       getEnv("JWT_SECRET", "default-secret-change-in-production"),
		HostFolder:      getEnv("HOST_FOLDER", "./published_sites"),
		TemplatesFolder: getEnv("TEMPLATES_FOLDER", "./templates"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "wawebv2"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
