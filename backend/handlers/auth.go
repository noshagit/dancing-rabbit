package handlers

import (
	"database/sql"
	"encoding/json"
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
		http.ServeFile(w, r, "../../frontend/connexion/connexion.html")
	}).Methods("GET")

	router.HandleFunc("/connexion/connexion.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/connexion/connexion.css")
	}).Methods("GET")

	router.HandleFunc("/connexion/connexion.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/connexion/connexion.js")
	}).Methods("GET")

	router.HandleFunc("/connexion/connexion.html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode de requête invalide", http.StatusMethodNotAllowed)
			return
		}

		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Erreur de décodage JSON", http.StatusBadRequest)
			log.Println("Erreur de décodage JSON:", err)
			return
		}

		db, err := sql.Open("sqlite3", "../../backend/database/dancing.db")
		if err != nil {
			http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
			log.Println("Erreur de connexion à la base de données:", err)
			return
		}
		defer db.Close()

		var storedHashedPassword string
		row := db.QueryRow("SELECT password FROM users WHERE email = ?", credentials.Email)
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

		if bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(credentials.Password)) == nil {
			sessionToken := uuid.New().String()

			_, err := db.Exec("INSERT INTO sessions (token, email) VALUES (?, ?)", sessionToken, credentials.Email)
			if err != nil {
				http.Error(w, "Erreur lors de la création de la session", http.StatusInternalServerError)
				log.Println("Erreur insertion session:", err)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:   "session_token",
				Value:  sessionToken,
				Path:   "/",
				MaxAge: 86400,
			})

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Connexion réussie"))
		} else {
			http.Error(w, "Mot de passe ou email invalide", http.StatusUnauthorized)
		}
	}).Methods("POST")
}

func RegisterHandler(router *mux.Router) {

	router.HandleFunc("/inscription/inscription.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/inscription/inscription.html")
	}).Methods("GET")

	router.HandleFunc("/inscription/inscription.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/inscription/inscription.css")
	}).Methods("GET")

	router.HandleFunc("/inscription/inscription.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/inscription/inscription.js")
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

		db, err := sql.Open("sqlite3", "../../backend/database/dancing.db")
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

	router.HandleFunc("/profil/profil.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/profil/profil.html")
	}).Methods("GET")

	router.HandleFunc("/profil/profil.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/profil/profil.css")
	}).Methods("GET")

	router.HandleFunc("/profil/profil.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/profil/profil.js")
	}).Methods("GET")

	router.HandleFunc("/frontend/images/bunny.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/images/bunny.png")
	}).Methods("GET")

	router.HandleFunc("/api/get-profile", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Non authentifié", http.StatusUnauthorized)
			return
		}

		db, err := sql.Open("sqlite3", "../../backend/database/dancing.db")
		if err != nil {
			http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var email string
		row := db.QueryRow("SELECT email FROM sessions WHERE token = ?", cookie.Value)
		err = row.Scan(&email)
		if err != nil {
			http.Error(w, "Session invalide", http.StatusUnauthorized)
			return
		}

		var profile struct {
			Pseudo string `json:"pseudo"`
			Email  string `json:"email"`
		}

		row = db.QueryRow("SELECT pseudo, email FROM users WHERE email = ?", email)
		err = row.Scan(&profile.Pseudo, &profile.Email)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des informations du profil", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"profile": profile,
		})
	}).Methods("GET")

}

func LogoutHandler(router *mux.Router) {
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err == nil {
			db, err := sql.Open("sqlite3", "../../backend/database/dancing.db")
			if err == nil {
				_, err := db.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value)
				if err != nil {
					http.Error(w, "Erreur lors de la suppression de la session", http.StatusInternalServerError)
					log.Println("Erreur lors de la suppression de la session:", err)
				}
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

		http.Redirect(w, r, "/main-menu/menu.html", http.StatusSeeOther)
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
