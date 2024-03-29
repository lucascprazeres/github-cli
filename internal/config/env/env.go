package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file cound't be loaded")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
