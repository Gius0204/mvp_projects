package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuración de la aplicación
type Config struct {
	DatabaseURL string
	KafkaBroker string
	Port        string
}

// Cargar variables de entorno
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del sistema")
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
		Port:        os.Getenv("PORT"),
	}
}
