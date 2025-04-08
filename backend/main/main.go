package main

import (
	"Backend/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// ====== ROUTES ====== //

func landingPageHandler(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/index.html")
	}).Methods("GET")

	router.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/index.js")
	}).Methods("GET")

	router.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/index.css")
	}).Methods("GET")
}

func presentationHandler(router *mux.Router) {
	router.HandleFunc("/main-menu/menu.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/main-menu/menu.html")
	}).Methods("GET")

	router.HandleFunc("/main-menu/menu.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/main-menu/menu.js")
	}).Methods("GET")

	router.HandleFunc("/main-menu/menu.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/main-menu/menu.css")
	}).Methods("GET")

	router.HandleFunc("/images/petit-bac.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/images/petit-bac.png")
	}).Methods("GET")

	router.HandleFunc("/images/deaf-rhythm.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/images/deaf-rhythm.png")
	}).Methods("GET")

	router.HandleFunc("/images/blind-test.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/images/blind-test.png")
	}).Methods("GET")
}

// ====== MAIN ====== //

func main() {

	router := mux.NewRouter()

	landingPageHandler(router)
	presentationHandler(router)

	handlers.RegisterHandler(router)
	handlers.LoginHandler(router)

	handlers.ProfileHandler(router)
	handlers.LogoutHandler(router)

	handlers.PetitBacMenuHandler(router)

	handlers.BlindTestMenuHandler(router)

	handlers.DeafRhythmMenuHandler(router)

	fmt.Println("Le serveur est lancer sous le port 8080 : http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
