package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// FileField is name of the field with a file
const FileField = "file"

// UploadDir point to the place where files are stored
const UploadDir = "./files/"

var db dbStruct

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
	fileInfo.name = FileHeader.Filename
	fileInfo.hash, err = createFileName(FileHeader.Filename, "k4wo")
	if err != nil {
		// TODO: Find better approach!
		fileInfo.hash, err = createFileName(FileHeader.Filename, "k4wo")
	}

	db.saveImage(&fileInfo)
	jsonResponse(w, http.StatusCreated, STRINGS["uploadedSuccessfully"])
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, message)
}

func main() {
	gotenv.Load()
	db = dbConnection()

	router := httprouter.New()
	router.POST("/upload", UploadFile)
	router.ServeFiles("/*filepath", http.Dir("./public"))

	log.Println("Running")
	http.ListenAndServe(":8080", router)
}
