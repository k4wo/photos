package main

func getImages() ([]File, error) {
	query := `
		SELECT
			owner,
			name,
			hash,
			size,
			extension,
			mime,
			latitude,
			longitude,
			orientation,
			model,
			camera,
			iso,
			focal_length,
			exposure_time,
			f_number,
			height,
			width,
			date
		FROM files`
	rows, _ := db.Query(query)
	defer rows.Close()

	return filesScanner(rows)
}

func saveImage(image *File) {
	sql := `
		INSERT INTO files (
			type, owner, name, hash, size, extension, 
			mime, latitude, longitude, orientation, 
			model, camera, iso, focal_length, 
			exposure_time, f_number, height, 
			width, date
		) 
		VALUES 
			(
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
				$11, $12, $13, $14, $15, $16, $17, $18, $19
			)	
	`

	_, err := db.Exec(
		sql,
		fileType["image"],
		1, // TODO: change it
		image.Name,
		image.Hash,
		image.Size,
		image.Extension,
		image.MimeType,
		image.Latitude,
		image.Longitude,
		image.Orientation,
		image.Model,
		image.Camera,
		image.Iso,
		image.FocalLength,
		image.ExposureTime,
		image.FNumber,
		image.Height,
		image.Width,
		image.Date,
	)

	if err != nil {
		panic(err)
	}
}
