package dev

import (
	"database/sql"
	"io/ioutil"
)

var db *sql.DB

func dropDatabase() {
	db.Exec("drop schema public cascade;")
	db.Exec("create schema public;")
}

func executeSQLFile(fileDestination string) {
	query, err := ioutil.ReadFile(fileDestination)
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}

// ResetDatabase removes existing db and creates new one with fake data
func ResetDatabase(database *sql.DB) {
	db = database
	dropDatabase()
	executeSQLFile("../dev/database/schema.sql")
	executeSQLFile("../dev/database/users.sql")
	executeSQLFile("../dev/database/files.sql")
	executeSQLFile("../dev/database/albums.sql")
	executeSQLFile("../dev/database/user_album.sql")
	executeSQLFile("../dev/database/album_file.sql")
}
