package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func dbConnection() *sql.DB {
	dbConfig := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	database, err := sql.Open("postgres", dbConfig)
	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	log.Info().Str("database", "postgres").Msg("Connected successfully")

	return database
}
