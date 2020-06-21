package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	constants "photos/constants"
	model "photos/model"
)

var selectAlbum = `
	SELECT
		albums.id,
		albums.owner,
		albums.name,
		albums.size,
		albums.updated_at,
		albums.created_at,

		files.id,
		files.owner,
		files.name,
		files.hash,
		files.size,
		files.extension,
		files.mime,
		files.latitude,
		files.longitude,
		files.orientation,
		files.model,
		files.camera,
		files.iso,
		files.focal_length,
		files.exposure_time,
		files.f_number,
		files.height,
		files.width,
		files.date
	FROM albums
	LEFT JOIN files ON albums.cover = files.id
	WHERE albums.owner = $1
`

func hasAlbumAccess(userID int, albumID string, db *sql.DB) bool {
	var count int
	rawQuery := `SELECT count(id) FROM albums WHERE owner = $1 AND id = $2`

	row := db.QueryRow(rawQuery, userID, albumID)
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return count > 0
}

func isFileInAlbum(fileID int, albumID string, db *sql.DB) bool {
	var count int
	rawQuery := `SELECT count(id) FROM album_file WHERE file = $1 AND album = $2`

	row := db.QueryRow(rawQuery, fileID, albumID)
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
		return true
	}

	return count > 0
}

func getAlbumByName(name string, userID int, db *sql.DB) (model.Album, error) {
	query := selectAlbum + " AND name = $2"
	row := db.QueryRow(query, userID, name)

	return albumScanner(row)
}

func GetAlbumContent(userID int, albumID string, db *sql.DB) ([]model.File, error) {
	rawQuery := `
		SELECT
			files.id,
			files.owner,
			files.name,
			files.hash,
			files.size,
			files.extension,
			files.mime,
			files.latitude,
			files.longitude,
			files.orientation,
			files.model,
			files.camera,
			files.iso,
			files.focal_length,
			files.exposure_time,
			files.f_number,
			files.height,
			files.width,
			files.date
		FROM
			album_file
			LEFT JOIN files ON files.id = album_file.file
			LEFT JOIN user_album ON user_album.album = album_file.album
		WHERE
			album_file. "album" = $1
			AND album_file. "user" = $2;
	`

	rows, _ := db.Query(rawQuery, albumID, userID)
	defer rows.Close()

	return filesScanner(rows)
}

func GetAlbums(userID int, db *sql.DB) ([]model.Album, error) {
	rows, _ := db.Query(selectAlbum, userID)
	defer rows.Close()

	return albumsScanner(rows)
}

func CreateAlbum(userID int, name string, db *sql.DB) (model.Album, error) {
	if name == "" {
		msg := fmt.Sprintf(constants.STRINGS["noAlbumName"])
		return model.Album{}, errors.New(msg)
	}

	album, err := getAlbumByName(name, userID, db)
	if err != sql.ErrNoRows {
		return album, err
	}

	query := `INSERT INTO albums(owner, name) VALUES($1, $2) RETURNING *`
	row := db.QueryRow(
		query,
		userID,
		name,
	)

	return albumScanner(row)
}

func AddFilesToAlbum(albumID string, userID int, files []int, db *sql.DB) int {
	hasAccess := hasAlbumAccess(userID, albumID, db)
	if !hasAccess {
		return http.StatusForbidden
	}

	for _, fileID := range files {
		if !hasFileAccess(userID, fileID, db) {
			fmt.Println("addFilesToAlbum", "no access", fileID)

			return http.StatusForbidden
		}
	}

	for _, fileID := range files {
		if !isFileInAlbum(fileID, albumID, db) {
			rawQuery := `INSERT INTO "album_file"("album", "file", "user") VALUES($1, $2, $3);`
			_, err := db.Exec(
				rawQuery,
				albumID,
				fileID,
				userID,
			)

			if err != nil {
				fmt.Println("addFilesToAlbum", err)
			}
		}
	}

	return http.StatusOK
}

func SetAlbumCover(albumID string, userID, fileID int, db *sql.DB) (int, model.File) {
	hasAccess := hasAlbumAccess(userID, albumID, db)
	if !hasAccess {
		return http.StatusForbidden, model.File{}
	}

	if !hasFileAccess(userID, fileID, db) {
		fmt.Println("setAlbumCover", "no access", fileID)
		return http.StatusForbidden, model.File{}
	}

	if isFileInAlbum(fileID, albumID, db) {
		rawQuery := `UPDATE albums SET cover = $1 WHERE id = $2`
		db.Exec(
			rawQuery,
			fileID,
			albumID,
		)

		file, err := getFileByID(fileID, db)
		if err != nil {
			fmt.Println("setAlbumCover", err)
			return http.StatusInternalServerError, model.File{}
		}

		return http.StatusOK, file
	}

	return http.StatusBadRequest, model.File{}
}

func RemoveFromAlbum(albumID string, userID, fileID int, db *sql.DB) int {
	hasAccess := hasAlbumAccess(userID, albumID, db)
	if !hasAccess {
		return http.StatusForbidden
	}

	rawQuery := `DELETE FROM "album_file" WHERE file = $1;`
	db.Exec(rawQuery, fileID)

	return http.StatusOK
}

func DeleteAlbum(albumID string, userID int, db *sql.DB) error {
	query := `DELETE FROM albums WHERE id = $1 AND owner = $2`
	_, err := db.Exec(
		query,
		albumID,
		userID,
	)

	return err
}
