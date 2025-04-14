package main

// package games

import (
	"math/rand/v2"
	"net/http"
	"slices"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/net/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	pastLetters   []rune
	currentLetter rune
	categories    = []string{}
	timer         *time.Timer
	currentRound  int
	maxRounds     int
)

type player struct {
	name    string
	score   int
	answers []string
}

func main() {

}

func randomLetter() {
	i := rand.IntN(26)
	currentLetter = rune('A' + i)
	if slices.Contains(pastLetters, currentLetter) {
		randomLetter()
	} else {
		pastLetters = append(pastLetters, currentLetter)
	}
}
