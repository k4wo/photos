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

type dbStruct struct {
	connection *sql.DB
}

func (store *dbStruct) saveImage(image *ImageInfo) {
	sql := `
		INSERT INTO files (
			type, owner, name, hash, size, extension, 
			mime, latitude, longitude, orientation, 
			model, camera, iso, focal_length, 
			exposure_time, f_number, y_dimension, 
			x_dimension, date
		) 
		VALUES 
			(
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
				$11, $12, $13, $14, $15, $16, $17, $18, $19
			)	
	`

	_, err := store.connection.Exec(
		sql,
		fileType["image"],
		1, // TODO: change it
		image.name,
		image.hash,
		image.size,
		image.extension,
		image.mimeType,
		image.latitude,
		image.longitude,
		image.orientation,
		image.model,
		image.camera,
		image.iso,
		image.focalLength,
		image.exposureTime,
		image.fNumber,
		image.yDimension,
		image.xDimension,
		image.time,
	)

	if err != nil {
		panic(err)
	}
}

func dbConnection() dbStruct {
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

	return dbStruct{connection: database}
}
