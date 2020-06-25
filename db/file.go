package db

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"photos/constants"
	"photos/image"
	model "photos/model"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v3"
)

var selectFile = `
	SELECT
		id,
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
	FROM files
`

func hasFileAccess(userID, fileID int, db *sql.DB) bool {
	var count int
	rawQuery := "SELECT count(id) FROM files WHERE id = $1 AND owner = $2"

	row := db.QueryRow(rawQuery, fileID, userID)
	err := row.Scan(&count)
	if err != nil {
		log.Error().Err(err).Caller().Int("user", userID).Int("file", fileID)

		return false
	}

	return count > 0
}

// GetFiles gets all files which belongs to a user
func GetFiles(userID int, db *sql.DB) ([]model.File, error) {
	query := selectFile + " WHERE owner = $1"
	rows, _ := db.Query(query, userID)
	defer rows.Close()

	return filesScanner(rows)
}

func getFileByID(fileID int, db *sql.DB) (model.File, error) {
	query := selectFile + " WHERE id = $1"
	row := db.QueryRow(query, fileID)

	return fileScanner(row)
}

func saveFile(file *model.File, userID int, db *sql.DB) bool {
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
		constants.FileType["image"],
		userID,
		file.Name,
		file.Hash,
		file.Size,
		file.Extension,
		file.MimeType,
		file.Latitude,
		file.Longitude,
		file.Orientation,
		file.Model,
		file.Camera,
		file.Iso,
		file.FocalLength,
		file.ExposureTime,
		file.FNumber,
		file.Height,
		file.Width,
		file.Date,
	)

	if err != nil {
		log.Error().Err(err).Caller().Int("user", userID).Msg("Problem with inserting a file")

		return false
	}

	return true
}

// DeleteFiles deletes only those files which are owned by a user.
// For deleting related entries in other tables `album_file`, `user_file`, etc.
// db takes care. Returning id of not inserted files.
// TODO: add physically removing from disc / storage
func DeleteFiles(filesID []int, userID int, db *sql.DB) []int {
	var notInserted []int
	query := "DELETE FROM files WHERE id = $1"

	for _, file := range filesID {
		if hasFileAccess(userID, file, db) {
			result, err := db.Exec(query, file)
			rowsNo, _ := result.RowsAffected()

			if err != nil || rowsNo == 0 {
				notInserted = append(notInserted, file)

				log.Error().
					Err(err).
					Caller().
					Int("user", userID).
					Int("file", file).
					Msg("Problem with deleting a file")
			}
		} else {
			notInserted = append(notInserted, file)

			log.Warn().
				Caller().
				Int("user", userID).
				Int("file", file).
				Msg("Problem with deleting a file")
		}
	}

	return notInserted
}

func writeFile(
	file multipart.File,
	FileHeader *multipart.FileHeader,
	userID int,
	uploadDir string,
	db *sql.DB,
) (*model.File, error) {

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Caller().Int("user", userID).Msg("Can't read a file")

		return &model.File{}, err
	}

	nullHash, _ := createFileName(FileHeader.Filename, userID, db)
	hash := nullHash.ValueOrZero()

	err = ioutil.WriteFile(uploadDir+hash, data, 0666)
	if err != nil {
		log.Error().Err(err).Caller().Int("user", userID).Str("hash", hash).Msg("Can't write a file")

		return &model.File{}, err
	}

	fileInfo, _ := image.ExtractExif(data)
	fileInfo.Name = null.StringFrom(FileHeader.Filename)
	fileInfo.Hash = nullHash
	image.ResizeImage(data, fileInfo, uploadDir)

	return &fileInfo, nil
}

func createFileName(name string, user int, db *sql.DB) (null.String, error) {
	today := time.Now()
	now := today.UnixNano()
	fileName := fmt.Sprintf("%s_%d_%d", name, user, now)

	h := sha1.New()
	_, err := io.WriteString(h, fileName)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Int("user", user).
			Str("name", name).
			Msg("Can't create a file name")

		return null.NewString("", false), err
	}

	return null.StringFrom(fmt.Sprintf("%x", h.Sum(nil))), nil
}

// ProcessFiles saves on disk file and than insert data to db. It accepts only jpeg/png so far.
func ProcessFiles(files []*multipart.FileHeader, userID int, uploadDir string, db *sql.DB) int {
	for _, file := range files {
		f, err := file.Open()

		if err != nil {
			log.Error().Err(err).Caller().Int("user", userID).Msg("Can't open a file")

			return http.StatusInternalServerError
		}

		defer f.Close()
		mimeType := file.Header.Get("Content-Type")

		if mimeType == "image/jpeg" || mimeType == "image/png" {
			fileInfo, err := writeFile(f, file, userID, uploadDir, db)

			if err != nil {
				log.Error().Err(err).Caller().Int("user", userID).Msg("Failed write a file")

				return http.StatusInternalServerError
			}

			// TODO: after fail remove the file
			saveFile(fileInfo, userID, db)
		} else {
			return http.StatusBadRequest
		}
	}

	return http.StatusOK
}
