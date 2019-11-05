package main

import (
	"crypto/sha1"
	. "fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

const FileField = "file"
const UploadDir = "./files/"

func createFileName(name, user string) (string, error) {
	today := time.Now()
	now := today.UnixNano()
	fileName := Sprintf("%s_%s_%d", name, user, now)

	h := sha1.New()
	_, err := io.WriteString(h, fileName)
	if err != nil {
		return "", err
	}

	return Sprintf("%x", h.Sum(nil)), nil
}

func saveFile(w http.ResponseWriter, file multipart.File, FileHeader *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		_, _ = Fprintf(w, "%v", err)
		return
	}

	err = ioutil.WriteFile(UploadDir+FileHeader.Filename, data, 0666)
	if err != nil {
		_, _ = Fprintf(w, "%v", err)
		return
	}

	fileInfo, _ := extractExif(data)
	fileInfo.name, err = createFileName(FileHeader.Filename, "k4wo")
	if err != nil {
		// TODO: Find better approach!
		fileInfo.name, err = createFileName(FileHeader.Filename, "k4wo")
	}

	Println(fileInfo, fileInfo.extension, fileInfo.name)

	jsonResponse(w, http.StatusCreated, STRINGS["uploadedSuccessfully"])
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = Fprint(w, message)
}

func main() {
	router := httprouter.New()
	router.POST("/upload", UploadFile)
	router.ServeFiles("/*filepath", http.Dir("./public"))

	log.Println("Running")
	http.ListenAndServe(":8080", router)
}
