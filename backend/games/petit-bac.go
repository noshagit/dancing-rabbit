package main

// package games

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"slices"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	tmpl     *template.Template
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	pastLetters   []rune
	currentLetter rune
	categories    = []string{}
	timer         *time.Timer
	timerDuration = 30
	currentRound  int
	maxRounds     int
)

type player struct {
	name    string
	score   int
	answers []string
}

func main() {
	tmpl = template.Must(template.ParseFiles(
		"../../frontend/petit-bac/petit-bac-menu.html",
		"../../frontend/petit-bac/game/petit-bac.html"))

	r := mux.NewRouter()
	r.HandleFunc("/", menu).Methods("GET")
	r.HandleFunc("/ws", handleWS).Methods("GET")
	r.HandleFunc("/petit-bac/game", game).Methods("GET")

	r.HandleFunc("/petit-bac-menu.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/petit-bac/petit-bac-menu.js")
	}).Methods("GET")
	r.HandleFunc("/petit-bac-menu.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/petit-bac/petit-bac-menu.css")
	}).Methods("GET")
	r.HandleFunc("/petit-bac.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/petit-bac/game/petit-bac.js")
	}).Methods("GET")
	r.HandleFunc("/petit-bac.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/petit-bac/game/petit-bac.css")
	}).Methods("GET")

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}

func menu(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "petit-bac-menu.html", nil)
}

func game(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "petit-bac.html", nil)
}

func handleWS(w http.ResponseWriter, r *http.Request) {}

func randomLetter() {
	i := rand.IntN(26)
	currentLetter = rune('A' + i)
	if slices.Contains(pastLetters, currentLetter) {
		randomLetter()
	} else {
		pastLetters = append(pastLetters, currentLetter)
	}
}
