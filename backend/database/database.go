package database

import (
	"fmt"
	"database/sql"
)

func ConnectToDatabase() {
	db, err := sql.Open("sqlite3", "./dancing.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connexion réussie à la base SQLite !")
}