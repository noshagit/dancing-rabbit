package Backend

import (
	"fmt"
	"net/http"
	"log"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

// TODO : Routes vers page de login

func loginHandler(router *mux.Router) {
	router.HandleFunc("/login", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "login")
	}).Methods("POST")
}

// TODO : Routes vers page de register

func registerHandler(router *mux.Router) {
	router.HandleFunc("/register", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "register")
	}).Methods("POST")
}

// TODO : Routes vers page de profil

func profileHandler(router *mux.Router) {
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

func checkPassword(password, hashedPassword string) bool {
	hashedInput := hashPassword(password)
	return hashedInput == hashedPassword
}