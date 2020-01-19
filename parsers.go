package main

import (
	"database/sql"
)

func filesParser(rows *sql.Rows) ([]ImageInfo, error) {
	var images []ImageInfo
	for rows.Next() {
		image := ImageInfo{}
		err := rows.Scan(
			&image.Owner,
			&image.Name,
			&image.Hash,
			&image.Size,
			&image.Extension,
			&image.MimeType,
			&image.Latitude,
			&image.Longitude,
			&image.Orientation,
			&image.Model,
			&image.Camera,
			&image.Iso,
			&image.FocalLength,
			&image.ExposureTime,
			&image.FNumber,
			&image.Height,
			&image.Width,
			&image.Date,
		)

		if err == nil {
			images = append(images, image)
		}

		if err != nil {
			return images, err
		}
	}

	return images, nil
}
