package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del entorno.")
	}
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3014"
	}
	return port
}

func GetRedisAddr() string {
	return os.Getenv("REDIS_ADDR")
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}
