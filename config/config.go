package config

import (
	"fmt"
	"os"
)

func DatabaseConfig() string {
	var (
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbName     = os.Getenv("DB_NAME")
		timeZone   = os.Getenv("TIMEZONE")
	)
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=%s",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		timeZone,
	)
}
