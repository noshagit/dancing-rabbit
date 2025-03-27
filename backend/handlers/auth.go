package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(router *mux.Router) {
	router.HandleFunc("/login", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "login")
	}).Methods("POST")
}

func RegisterHandler(router *mux.Router) {
	router.HandleFunc("/register", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "register")
	}).Methods("POST")
}

func ProfileHandler(router *mux.Router) {
	router.HandleFunc("/profil", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "profil")
	}).Methods("POST")
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func CheckPassword(password, hashedPassword string) bool {
	hashedInput := hashPassword(password)
	return hashedInput == hashedPassword
}
