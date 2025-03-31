package handlers

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func PetitBacHandler(router *mux.Router) {
	router.HandleFunc("/game/petitback", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "petitback")
	}).Methods("POST")
}

func BlindTestHandler(router *mux.Router) {
	router.HandleFunc("/game/blindtest", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "blindtest")
	}).Methods("POST")
}

func DeafRhythmHandler(router *mux.Router) {
	router.HandleFunc("/game/deathrhythm", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "")
	}).Methods("POST")
}
