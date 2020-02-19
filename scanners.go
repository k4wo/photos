package main

import (
	"database/sql"
	model "photos/model"
)

func filesScanner(rows *sql.Rows) ([]model.File, error) {
	var images []model.File
	for rows.Next() {
		image := model.File{}
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

func albumsScanner(rows *sql.Rows) ([]model.Album, error) {
	var albums []model.Album
	for rows.Next() {
		album := model.Album{}
		err := rows.Scan(
			&album.ID,
			&album.Owner,
			&album.Name,
			&album.Size,
			&album.Cover,
			&album.UpdatedAt,
			&album.CreatedAt,
		)

		if err == nil {
			albums = append(albums, album)
		}

		if err != nil {
			return albums, err
		}
	}

	return albums, nil
}

func albumScanner(row *sql.Row) (model.Album, error) {
	album := model.Album{}
	err := row.Scan(
		&album.ID,
		&album.Owner,
		&album.Name,
		&album.Size,
		&album.UpdatedAt,
		&album.CreatedAt,
		&album.Cover,
	)

	return album, err
}
