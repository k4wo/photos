package db

import (
	"database/sql"
	model "photos/model"
)

func filesScanner(rows *sql.Rows) ([]model.File, error) {
	var images []model.File
	for rows.Next() {
		image := model.File{}
		err := rows.Scan(
			&image.ID,
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
		} else {
			return images, err
		}
	}

	return images, nil
}

func fileScanner(row *sql.Row) (model.File, error) {
	file := model.File{}
	err := row.Scan(
		&file.ID,
		&file.Owner,
		&file.Name,
		&file.Hash,
		&file.Size,
		&file.Extension,
		&file.MimeType,
		&file.Latitude,
		&file.Longitude,
		&file.Orientation,
		&file.Model,
		&file.Camera,
		&file.Iso,
		&file.FocalLength,
		&file.ExposureTime,
		&file.FNumber,
		&file.Height,
		&file.Width,
		&file.Date,
	)

	return file, err
}

func albumScannerWithoutCover(row *sql.Row) (model.Album, error) {
	album := model.Album{}

	err := row.Scan(
		&album.ID,
		&album.Owner,
		&album.Name,
		&album.Size,
		&album.Cover,
		&album.UpdatedAt,
		&album.CreatedAt,
	)

	return album, err
}

func albumsScanner(rows *sql.Rows) ([]model.Album, error) {
	var albums []model.Album

	for rows.Next() {
		album := model.Album{}
		file := model.File{}

		err := rows.Scan(
			&album.ID,
			&album.Owner,
			&album.Name,
			&album.Size,
			&album.UpdatedAt,
			&album.CreatedAt,

			&file.ID,
			&file.Owner,
			&file.Name,
			&file.Hash,
			&file.Size,
			&file.Extension,
			&file.MimeType,
			&file.Latitude,
			&file.Longitude,
			&file.Orientation,
			&file.Model,
			&file.Camera,
			&file.Iso,
			&file.FocalLength,
			&file.ExposureTime,
			&file.FNumber,
			&file.Height,
			&file.Width,
			&file.Date,
		)

		if err != nil {
			return albums, err
		}
		if file.ID.Valid {
			album.File = model.Cover{true, file}
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func albumScanner(row *sql.Row) (model.Album, error) {
	album := model.Album{}
	file := model.File{}

	err := row.Scan(
		album.ID,
		album.Owner,
		album.Name,
		album.Size,
		album.UpdatedAt,
		album.CreatedAt,

		&file.ID,
		&file.Owner,
		&file.Name,
		&file.Hash,
		&file.Size,
		&file.Extension,
		&file.MimeType,
		&file.Latitude,
		&file.Longitude,
		&file.Orientation,
		&file.Model,
		&file.Camera,
		&file.Iso,
		&file.FocalLength,
		&file.ExposureTime,
		&file.FNumber,
		&file.Height,
		&file.Width,
		&file.Date,
	)

	return album, err
}
