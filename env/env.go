package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func FetchString(key string, fallback string) string {
	response, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return response
}

func FetchInt(key string, fallback int) int {
	response, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	resp, err := strconv.Atoi(response)
	if err != nil {
		return fallback
	}
	return resp
}
