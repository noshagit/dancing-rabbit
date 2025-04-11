package main

//package games

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"slices"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/tsirysndr/go-deezer"
)

var (
	playlist      *deezer.Tracks
	previousSongs []string
	currentSong   deezer.Track
	playerGuesses map[string]string
)

func main() {
	getTrack()
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")

	r.HandleFunc("/blind-test.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/blind-test/game/blind-test.css")
	})
	r.HandleFunc("/blind-test.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/blind-test/game/blind-test.js")
	})

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HANDLER")

	tmpl := template.Must(template.ParseFiles("../../frontend/blind-test/game/blind-test.html"))
	err2 := tmpl.Execute(w, currentSong.Preview)
	if err2 != nil {
		log.Fatalf("Failed to execute template: %v", err2)
	}
}

func getTrack() {
	fmt.Println("GET TRACK")

	client := deezer.NewClient()
	playlist, err1 := client.Playlist.GetTracks("53362031")
	if err1 != nil {
		log.Fatalf("Failed to get playlist: %v", err1)
	}

	randomTrackIndex := rand.IntN(playlist.Total)
	randomTrack := playlist.Data[randomTrackIndex]
	if slices.Contains(previousSongs, randomTrack.Title) {
		fmt.Println("song already played")
		getTrack()
		return
	}
	currentSong = randomTrack
	previousSongs = append(previousSongs, randomTrack.Title)
}
