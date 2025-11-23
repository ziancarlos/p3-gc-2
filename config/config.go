package config

import (
	"log"
	"os"
)

type Config struct {
	PortShopping          string
	PortPayment           string
	MongoURI              string
	ShoppingDBName        string
	PaymentDBName         string
	PaymentServiceBaseURI string
	JWTSecret             string
}

func LoadConfig() *Config {
	return &Config{
		PortShopping:          getEnv("PORT_SHOPPING", "9051"),
		PortPayment:           getEnv("PORT_PAYMENT", "9061"),
		MongoURI:              getEnv("MONGO_URI", "mongodb://localhost:9071"),
		ShoppingDBName:        getEnv("SHOPPING_DB_NAME", "shopping_db"),
		PaymentDBName:         getEnv("PAYMENT_DB_NAME", "payment_db"),
		PaymentServiceBaseURI: getEnv("PAYMENT_SERVICE_BASE_URI", "localhost:9061"),
		JWTSecret:             getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

