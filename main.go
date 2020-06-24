package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/subosito/gotenv"
)

// FileField is name of the field with a file
const FileField = "file"

// UploadDir point to the place where files are stored
const UploadDir = "./files/"

const userID = 1 // TODO: use real userID
var db *sql.DB

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, message)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

func main() {
	gotenv.Load()
	if os.Getenv("ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	db = dbConnection()

	router := httprouter.New()
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	router.POST("/upload", uploadFilesRoute)
	router.GET("/images", fetchFilesRoute)

	router.GET("/albums", fetchAlbumsRoute)
	router.POST("/albums", addNewAlbumRoute)

	router.DELETE("/album/:id", deleteAlbumRoute)
	router.GET("/album/:id", fetchAlbumContentRoute)
	router.PUT("/album/:id/files", addFilesToAlbumRoute)
	router.DELETE("/album/:id/file", removeFromAlbumRoute)
	router.PUT("/album/:id/cover", setAlbumCoverRoute)

	router.DELETE("/files/delete", deleteFileRoute)

	router.ServeFiles("/files/*filepath", http.Dir("./files"))

	log.Info().Msg("Running")
	http.ListenAndServe(":8080", router)
}
