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
		fmt.Println(err)
		return false
	}

	return count > 0
}

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

func saveFile(image *model.File, userID int, db *sql.DB) {
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

func DeleteFiles(filesID []int, userID int, db *sql.DB) error {
	args := make([]interface{}, len(filesID))
	placeholder := ""
	for i, id := range filesID {
		placeholder = placeholder + fmt.Sprintf(", $%d", i+1)
		args[i] = id
	}

	query := "DELETE FROM files WHERE id IN (" + placeholder[2:] + ")"
	_, err := db.Exec(query, args...)

	return err
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
		return &model.File{}, err
	}

	// BUG: don't use real name for file name, can be overrided
	err = ioutil.WriteFile(uploadDir+FileHeader.Filename, data, 0666)
	if err != nil {
		return &model.File{}, err
	}

	fileInfo, _ := image.ExtractExif(data)
	fileInfo.Name = null.StringFrom(FileHeader.Filename)
	fileInfo.Hash, _ = createFileName(FileHeader.Filename, userID, db)
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
		return null.NewString("", false), err
	}

	return null.StringFrom(fmt.Sprintf("%x", h.Sum(nil))), nil
}

func ProcessFiles(files []*multipart.FileHeader, userID int, uploadDir string, db *sql.DB) int {
	for _, file := range files {
		f, err := file.Open()

		if err != nil {
			fmt.Println("processFile", err)

			return http.StatusInternalServerError
		}

		defer f.Close()
		mimeType := file.Header.Get("Content-Type")

		if mimeType == "image/jpeg" || mimeType == "image/png" {
			fileInfo, err := writeFile(f, file, userID, uploadDir, db)

			if err != nil {
				fmt.Println("processFile", err)

				return http.StatusInternalServerError
			}

			saveFile(fileInfo, userID, db)
		} else {
			return http.StatusBadRequest
		}
	}

	return http.StatusOK
}