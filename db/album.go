package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	constants "photos/constants"
	model "photos/model"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

// hasAlbumAccess checks if an user is an owner of the album or
// the album is shared with him
func hasAlbumAccess(userID int, albumID string, db *sql.DB) bool {
	var count int
	rawQuery := `
		SELECT
			count(*)
		FROM
			albums
			LEFT JOIN user_album ON user_album.album = albums.id
		WHERE 
			(albums.owner = $1 AND albums.id = $2)
			OR (user_album.album = $2 AND user_album.user = $1);
	`

	row := db.QueryRow(rawQuery, userID, albumID)
	err := row.Scan(&count)
	if err != nil {
		log.Error().Err(err).Caller().Int("user", userID).Str("album", albumID)

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
		log.Error().Err(err).Caller().Int("file", fileID).Str("album", albumID)

		return true
	}

	return count > 0
}

func getAlbumByName(name string, userID int, db *sql.DB) (model.Album, error) {
	query := selectAlbum + " AND albums.name = $2"
	row := db.QueryRow(query, userID, name)

	return albumScanner(row)
}

// GetAlbumContent returns files from an album where a user is an owner or the album
// is shared with the user
func GetAlbumContent(userID int, albumID string, db *sql.DB) ([]model.File, error) {
	hasAccess := hasAlbumAccess(userID, albumID, db)
	if !hasAccess {
		return []model.File{}, errors.New(constants.STRINGS["noAccessToAlbum"])
	}

	// user has access to the album so take all files from the album
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
			files
			LEFT JOIN album_file ON files.id = album_file.file
		WHERE
			album_file."album" = $1;
	`

	rows, _ := db.Query(rawQuery, albumID)
	defer rows.Close()

	return filesScanner(rows)
}

// GetAlbums returns all albums with a cover where a user is an owner
func GetAlbums(userID int, db *sql.DB) ([]model.Album, error) {
	rows, _ := db.Query(selectAlbum, userID)
	defer rows.Close()

	return albumsScanner(rows)
}

// CreateAlbum creates an album and returns it without cover. Pass an id as cover
func CreateAlbum(userID int, name string, db *sql.DB) (model.Album, error) {
	if name == "" {
		return model.Album{}, errors.New(constants.STRINGS["noAlbumName"])
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

	return albumScannerWithoutCover(row)
}

// AddFilesToAlbum adds file(s) to the album where user is an owner or the album is shared with him
func AddFilesToAlbum(albumID string, userID int, files []int, db *sql.DB) int {
	hasAccess := hasAlbumAccess(userID, albumID, db)
	if !hasAccess {
		return http.StatusForbidden
	}

	for _, fileID := range files {
		if !hasFileAccess(userID, fileID, db) {
			log.Warn().
				Caller().
				Int("user", userID).
				Int("file", fileID).
				Msg("Don't have access to the file")

			return http.StatusForbidden
		}
	}

	for _, fileID := range files {
		if !isFileInAlbum(fileID, albumID, db) {
			rawQuery := `INSERT INTO "album_file"("album", "file", "added_by") VALUES($1, $2, $3);`
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

// SetAlbumCover sets an albums' cover only when a user have access to the album
// and a file is already in the album
func SetAlbumCover(albumID string, userID, fileID int, db *sql.DB) (int, model.File) {
	hasAccess := hasAlbumAccess(userID, albumID, db)
	if !hasAccess {
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
			log.Error().
				Err(err).
				Caller().
				Str("action", "getFileByID").
				Int("file", fileID)

			return http.StatusInternalServerError, model.File{}
		}

		return http.StatusOK, file
	}

	return http.StatusBadRequest, model.File{}
}

// RemoveFromAlbum removes a file from an album if a user has access to the album.
// Doesn't care whether the user is an owner of the file / album.
func RemoveFromAlbum(albumID string, userID, fileID int, db *sql.DB) int {
	hasAccess := hasAlbumAccess(userID, albumID, db)
	if !hasAccess {
		return http.StatusForbidden
	}

	rawQuery := `DELETE FROM "album_file" WHERE file = $1;`
	db.Exec(rawQuery, fileID)

	return http.StatusOK
}

// DeleteAlbum deletes an album by a user who is an owner. DB takes care of
// removing all related data from `user_album` and `album_file` tables
func DeleteAlbum(albumID string, userID int, db *sql.DB) error {
	hasAccess := hasAlbumAccess(userID, albumID, db)
	if !hasAccess {
		return errors.New(constants.STRINGS["noAccessToAlbum"])
	}

	query := `DELETE FROM albums WHERE id = $1 AND owner = $2`
	_, err := db.Exec(
		query,
		albumID,
		userID,
	)

	return err
}

func init() {
	if os.Getenv("ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
