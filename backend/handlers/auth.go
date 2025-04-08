package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
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

	router.HandleFunc("/connexion/connexion.html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "méthode de requête invalide", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		db, err := sql.Open("sqlite3", "/home/ilian/dancing-rabbit/backend/database/dancing.db")
		if err != nil {
			http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
			log.Println("Erreur de connexion à la base de données:", err)
			return
		}
		defer db.Close()

		var storedHashedPassword string
		row := db.QueryRow("SELECT password FROM users WHERE email = ?", email)
		err = row.Scan(&storedHashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Email ou mot de passe invalide", http.StatusUnauthorized)
			} else {
				http.Error(w, "Erreur lors de la requête à la base de données", http.StatusInternalServerError)
				log.Println("Erreur lors de la requête à la base de données:", err)
			}
			return
		}

		if CheckPassword(password, storedHashedPassword) {
			sessionToken := uuid.New().String()

			_, err := db.Exec("INSERT INTO sessions (token, email) VALUES (?, ?)", sessionToken, email)
			if err != nil {
				http.Error(w, "Erreur lors de la création de la session", http.StatusInternalServerError)
				log.Println("Erreur insertion session:", err)
				return
			}

			http.SetCookie(w, &http.Cookie{ // création du cookie
				Name:     "session_token",
				Value:    sessionToken,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   3600,
			})

			http.Redirect(w, r, "/main-menu/menu.html", http.StatusSeeOther)
		} else {
			http.Error(w, "Mot de passe ou email invalide", http.StatusUnauthorized)
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

		db, err := sql.Open("sqlite3", "/home/ilian/dancing-rabbit/backend/database/dancing.db")
		if err != nil {
			http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
			log.Println("Erreur de connexion DB:", err)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO users (pseudo, email, password) VALUES (?, ?, ?)")
		if err != nil {
			log.Println("Erreur de préparation DB:", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(user.Pseudo, user.Email, hashedPassword); err != nil {
			log.Println("Erreur d'insertion DB:", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Inscription réussie"))
	}).Methods("POST")
}

func ProfileHandler(router *mux.Router) {
	router.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Non autorisé (pas de cookie)", http.StatusUnauthorized)
			return
		}

		db, err := sql.Open("sqlite3", "/home/ilian/dancing-rabbit/backend/database/dancing.db")
		if err != nil {
			http.Error(w, "Erreur DB", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var email string
		err = db.QueryRow("SELECT email FROM sessions WHERE token = ?", cookie.Value).Scan(&email)
		if err != nil {
			http.Error(w, "Session invalide", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "Bienvenue sur le profil de %s !", email)
	}).Methods("GET")
}

func LogoutHandler(router *mux.Router) {
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err == nil {
			db, err := sql.Open("sqlite3", "/home/ilian/dancing-rabbit/backend/database/dancing.db")
			if err == nil {
				db.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value)
				db.Close()
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		fmt.Fprintln(w, "Déconnexion réussie")
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
