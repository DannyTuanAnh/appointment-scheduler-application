package utils

import (
	"log"
	"os"
	"strconv"

	"time"

	"github.com/joho/godotenv"
)

var (
	pathEnv = ".env"
)

// loadEnv is a helper function that loads environment variables from a .env file using the godotenv package
func LoadEnv() {
	err := godotenv.Load(pathEnv)
	if err != nil {
		log.Println("No .env file found")
		panic(err)
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func GetEnvTime(key string, defaultValue int) time.Duration {
	if value := os.Getenv(key); value != "" {
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Println("Error converting environment variable to int:", err)
		} else {
			return time.Duration(val)
		}
	}

	return time.Duration(defaultValue)
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Println("Error converting environment variable to int:", err)
		} else {
			return val
		}
	}

	return defaultValue
}
