package main

import (
	"Backend/database"
	"Backend/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// ====== ROUTES ====== //

func landingPageHandler(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, router *http.Request) {
		fmt.Println("Landing page")
	}).Methods("POST")
}

func presentationHandler(router *mux.Router) {
	router.HandleFunc("/game", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "pres")
	}).Methods("POST")
}

// ====== MAIN ====== //

func main() {
	database.ConnectToDatabase()

	router := mux.NewRouter()

	landingPageHandler(router)
	presentationHandler(router)

	handlers.RegisterHandler(router)
	handlers.LoginHandler(router)
	handlers.ProfileHandler(router)
	handlers.PetitBacHandler(router)
	handlers.BlindTestHandler(router)
	handlers.DeafRhythmHandler(router)

	fmt.Println("Le serveur est lancer sous le port 8080 : http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
