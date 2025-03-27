package main

import (
	"fmt"
	"net/http"
	"Backend/handlers"
    "Backend/database"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// ====== ROUTES ====== //

// TODO : Routes vers la racine ( landing page )
func landingPageHandler(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "welcome")
	}).Methods("POST")
}

// TODO : Routes vers page de présentation des jeux

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
