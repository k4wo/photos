package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// UploadFile uploads a file to the server
func UploadFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	enableCors(&w)

	r.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
	files := r.MultipartForm.File["files"]
	for _, file := range files {
		f, err := file.Open()

		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}
		defer f.Close()

		mimeType := file.Header.Get("Content-Type")
		switch mimeType {
		case "image/jpeg":
		case "image/png":
			saveFile(w, f, file)
		default:
			jsonResponse(w, http.StatusBadRequest, STRINGS["fileFormatInvalid"])
		}
	}

}
