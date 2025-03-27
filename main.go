package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Struct Categories
type Categories struct {
	Album        string
	Artist       string
	MusicGroup   string
	Song         string
	MusicalGenre string
}

// Struct Player
type Player struct {
	Name  string
	Score int
	Vote  int
}

// Struct of the Game
type Game struct {
	RandomLetter        string   // Random letter choose by the computer
	Round               int      // Round of the game
	LetterAlreadyChoose []string // Letter already choose by the computer
	Timer               int      // Timer of the game
}

// Function main
func main() {
	game := Game{}

	// Random letter
	letter := game.RandomLetterFunc()
	println(letter)

	// Check first letter
	result := game.CheckFirstLetter(game.UserInput())
	println(result)
}

// User Input (temporary)
func (g *Game) UserInput() string {
	var input string
	println("Enter a word")
	_, _ = fmt.Scanln(&input)
	return input
}

// Random letter
func (g *Game) RandomLetterFunc() string {
	rand.NewSource(time.Now().UnixNano())

	for {
		g.RandomLetter = string(rune(rand.Intn(26) + 65))

		alreadyChoosen := false
		for _, letter := range g.LetterAlreadyChoose {
			if letter == g.RandomLetter {
				alreadyChoosen = true
				break
			}
		}

		if !alreadyChoosen {
			break
		}

		if len(g.LetterAlreadyChoose) >= 26 {
			return ""
		}
	}

	g.LetterAlreadyChoose = append(g.LetterAlreadyChoose, g.RandomLetter)
	return g.RandomLetter
}

// Check first letter
func (g *Game) CheckFirstLetter(word string) bool {
	if strings.ToUpper(string(rune(word[0]))) == strings.ToUpper(g.RandomLetter) {
		return true
	}
	return false
}
