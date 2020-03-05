package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	image "photos/image"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

// FileField is name of the field with a file
const FileField = "file"

// UploadDir point to the place where files are stored
const UploadDir = "./files/"

var db *sql.DB

func createFileName(name, user string) (string, error) {
	today := time.Now()
	now := today.UnixNano()
	fileName := fmt.Sprintf("%s_%s_%d", name, user, now)

	h := sha1.New()
	_, err := io.WriteString(h, fileName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func saveFile(w http.ResponseWriter, file multipart.File, FileHeader *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}

	// BUG: don't use real name for file name, can be overrided
	err = ioutil.WriteFile(UploadDir+FileHeader.Filename, data, 0666)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}

	fileInfo, _ := image.ExtractExif(data)
	fileInfo.Name = FileHeader.Filename
	fileInfo.Hash, err = createFileName(FileHeader.Filename, "k4wo")
	if err != nil {
		// TODO: Find better approach!
		fileInfo.Hash, err = createFileName(FileHeader.Filename, "k4wo")
	}

	saveImage(&fileInfo)
	image.ResizeImage(data, fileInfo, UploadDir)
	jsonResponse(w, http.StatusCreated, STRINGS["uploadedSuccessfully"])
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, message)
}

func fetchImages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)
	image, err := getImages()
	if err != nil {
		fmt.Println("handle error")
	}

	json.NewEncoder(w).Encode(image)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	gotenv.Load()
	db = dbConnection()

	router := httprouter.New()

	router.POST("/upload", UploadFile)
	router.GET("/images", fetchImages)

	router.GET("/albums", fetchAlbums)
	router.POST("/albums", addNewAlbum)

	router.DELETE("/album/:id", deleteAlbum)
	router.GET("/album/:id", fetchAlbumContent)

	router.PUT("/files/delete", deleteFileRoute)

	router.ServeFiles("/files/*filepath", http.Dir("./files"))

	log.Println("Running")
	http.ListenAndServe(":8080", router)
}
