package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PublicHost string
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

const PUBLIC_HOST = "PUBLIC_HOST"
const PORT = "PORT"
const DB_USER = "DB_USER"
const DB_PASSWORD = "DB_PASSWORD"
const DB_ADDRESS = "DB_ADDRESS"
const DB_HOST = "DB_HOST"
const DB_PORT = "DB_PORT"
const DB_NAME = "DB_NAME"

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	dbHost := getEnv(DB_HOST, "127.0.0.1")
	port := getEnv(DB_PORT, "3306")
	dbAddress := fmt.Sprintf("%s:%s", dbHost, port)

	return Config{
		PublicHost: getEnv(PUBLIC_HOST, "http://localhost"),
		Port:       getEnv(PORT, "8080"),
		DBUser:     getEnv(DB_USER, "root"),
		DBPassword: getEnv(DB_PASSWORD, "root"),
		DBAddress:  getEnv(DB_ADDRESS, dbAddress),
		DBName:     getEnv(DB_NAME, "ecom"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
