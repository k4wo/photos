package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"photos/image"
	model "photos/model"
	"time"
)

func hasFileAccess(userID, fileID int) bool {
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

func getFile(fileID int) ([]model.File, error) {
	query := `
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
		WHERE owner = $1`
	rows, _ := db.Query(query, fileID)
	defer rows.Close()

	return filesScanner(rows)
}

func saveFile(image *model.File) {
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

func deleteFiles(filesID []int, userID int) error {
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

func writeFile(file multipart.File, FileHeader *multipart.FileHeader, userID int) (*model.File, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return &model.File{}, err
	}

	// BUG: don't use real name for file name, can be overrided
	err = ioutil.WriteFile(UploadDir+FileHeader.Filename, data, 0666)
	if err != nil {
		return &model.File{}, err
	}

	fileInfo, _ := image.ExtractExif(data)
	fileInfo.Name = FileHeader.Filename
	fileInfo.Hash, err = createFileName(FileHeader.Filename, userID)
	image.ResizeImage(data, fileInfo, UploadDir)

	return &fileInfo, nil
}

func createFileName(name string, user int) (string, error) {
	today := time.Now()
	now := today.UnixNano()
	fileName := fmt.Sprintf("%s_%d_%d", name, user, now)

	h := sha1.New()
	_, err := io.WriteString(h, fileName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func processFiles(files []*multipart.FileHeader, userID int) int {
	for _, file := range files {
		f, err := file.Open()

		if err != nil {
			fmt.Println("processFile", err)

			return http.StatusInternalServerError
		}

		defer f.Close()
		mimeType := file.Header.Get("Content-Type")

		if mimeType == "image/jpeg" || mimeType == "image/png" {
			fileInfo, err := writeFile(f, file, userID)

			if err != nil {
				fmt.Println("processFile", err)

				return http.StatusInternalServerError
			}

			saveFile(fileInfo)
		} else {
			return http.StatusBadRequest
		}
	}

	return http.StatusOK
}
