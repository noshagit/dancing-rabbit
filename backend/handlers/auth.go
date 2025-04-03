package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Pseudo          string `json:"pseudo"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func LoginHandler(router *mux.Router) {
	router.HandleFunc("/connexion/connexion.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/connexion/connexion.html")
	}).Methods("GET")

	router.HandleFunc("/connexion/connexion.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/connexion/connexion.css")
	}).Methods("GET")

	router.HandleFunc("/connexion/connexion.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/connexion/connexion.js")
	}).Methods("GET")
	router.HandleFunc("/login/login.html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		db, err := sql.Open("sqlite3", "./dancing.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			log.Println("Database connection error:", err)
			return
		}
		defer db.Close()

		var storedHashedPassword string
		row := db.QueryRow("SELECT password FROM users WHERE email = ?", email)
		err = row.Scan(&storedHashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database query error", http.StatusInternalServerError)
				log.Println("Database query error:", err)
			}
			return
		}

		if CheckPassword(password, storedHashedPassword) {
			fmt.Fprintln(w, "Login successful")
		} else {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		}
	}).Methods("POST")
}

func RegisterHandler(router *mux.Router) {

	router.HandleFunc("/inscription/inscription.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/inscription/inscription.html")
	}).Methods("GET")

	router.HandleFunc("/inscription/inscription.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/inscription/inscription.css")
	}).Methods("GET")

	router.HandleFunc("/inscription/inscription.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/inscription/inscription.js")
	}).Methods("GET")

	router.HandleFunc("/inscription/inscription.html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode invalide", http.StatusMethodNotAllowed)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Erreur de décodage JSON", http.StatusBadRequest)
			log.Println("Erreur de décodage JSON:", err)
			return
		}

		if user.Pseudo == "" || user.Email == "" || user.Password == "" || user.ConfirmPassword == "" {
			http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
			return
		}

		if user.Password != user.ConfirmPassword {
			http.Error(w, "Les mots de passe ne correspondent pas", http.StatusBadRequest)
			return
		}

		hashedPassword := hashPassword(user.Password)
		if err != nil {
			http.Error(w, "Erreur de hachage du mot de passe", http.StatusInternalServerError)
			log.Println("Erreur de hachage du mot de passe:", err)
			return
		}

		db, err := sql.Open("sqlite3", "/home/ilian/dancing-rabbit/backend/database/dancing.db")
		if err != nil {
			http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
			log.Println("Erreur de connexion DB:", err)
			return
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO users (pseudo, email, password) VALUES (?, ?, ?)", user.Pseudo, user.Email, hashedPassword)
		if err != nil {
			http.Error(w, "Erreur lors de l'insertion dans la base", http.StatusInternalServerError)
			log.Println("Erreur d'insertion DB:", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Inscription réussie"))
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
