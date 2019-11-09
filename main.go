package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// FileField is name of the field with a file
const FileField = "file"

// UploadDir point to the place where files are stored
const UploadDir = "./files/"

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

	err = ioutil.WriteFile(UploadDir+FileHeader.Filename, data, 0666)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}

	fileInfo, _ := extractExif(data)
	fileInfo.name, err = createFileName(FileHeader.Filename, "k4wo")
	if err != nil {
		// TODO: Find better approach!
		fileInfo.name, err = createFileName(FileHeader.Filename, "k4wo")
	}

	jsonResponse(w, http.StatusCreated, STRINGS["uploadedSuccessfully"])
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, message)
}

func main() {
	gotenv.Load()

	dbConfig := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.POST("/upload", UploadFile)
	router.ServeFiles("/*filepath", http.Dir("./public"))

	log.Println("Running")
	http.ListenAndServe(":8080", router)
}
