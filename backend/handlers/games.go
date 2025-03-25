package handlers

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

// TODO : Routes vers le petit back

func PetitBacHandler(router *mux.Router) {
	router.HandleFunc("/game/petitback/{id}", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "petitback")
	}).Methods("POST")
}

// TODO : Routes vers le blind test

func BlindTestHandler(router *mux.Router) {
	router.HandleFunc("/game/blindtest/{id}", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "blindtest")
	}).Methods("POST")
}

// TODO : Routes vers death rhythm

func DeafRhythmHandler(router *mux.Router) {
	router.HandleFunc("/game/deathrhythm/{id}", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintln(w, "")
	}).Methods("POST")
}