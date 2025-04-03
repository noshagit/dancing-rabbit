package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func PetitBacMenuHandler(router *mux.Router) {
	router.HandleFunc("/petit-bac/petit-bac-menu.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/petit-bac/petit-bac-menu.html")
	}).Methods("GET")

	router.HandleFunc("/petit-bac/petit-bac-menu.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/petit-bac/petit-bac-menu.js")
	}).Methods("GET")

	router.HandleFunc("/petit-bac/petit-bac-menu.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/petit-bac/petit-bac-menu.css")
	}).Methods("GET")
}

func BlindTestMenuHandler(router *mux.Router) {
	router.HandleFunc("/blind-test/blind-test-menu.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/blind-test/blind-test-menu.html")
	}).Methods("GET")

	router.HandleFunc("/blind-test/blind-test-menu.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/blind-test/blind-test-menu.js")
	}).Methods("GET")

	router.HandleFunc("/blind-test/blind-test-menu.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/blind-test/blind-test-menu.css")
	}).Methods("GET")
}

func DeafRhythmMenuHandler(router *mux.Router) {
	router.HandleFunc("/deaf-rhythm/deaf-rhythm-menu.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/deaf-rhythm/deaf-rhythm-menu.html")
	}).Methods("GET")

	router.HandleFunc("/deaf-rhythm/deaf-rhythm-menu.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/deaf-rhythm/deaf-rhythm-menu.js")
	}).Methods("GET")

	router.HandleFunc("/deaf-rhythm/deaf-rhythm-menu.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/deaf-rhythm/deaf-rhythm-menu.css")
	}).Methods("GET")
}
