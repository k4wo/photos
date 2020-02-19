package main

import (
	"database/sql"
	"fmt"
	"os"
)

var fileType = map[string]string{
	"image":     "IMAGE",
	"video":     "VIDEO",
	"animation": "ANIMATION",
	"collage":   "COLLAGE",
}

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

	fmt.Println("Connected to DB successfully.")

	return database
}
