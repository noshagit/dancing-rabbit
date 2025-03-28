package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// Todo: Add a timer
// Todo: Add a score table
// Todo: Name condition
// Todo: Add a condition to check if the word is already used
// Todo: Add score system

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
	game.Menu()
}

// Menu
func (g *Game) Menu() {
	fmt.Println("--- Welcome to the Game ---")
	fmt.Println("1. Start the game")
	fmt.Println("2. Exit")
	fmt.Println("---------------------------")

	var choice int
	_, _ = fmt.Scanln(&choice)
	switch choice {
	case 1:
		g.StartGame()
	case 2:
		fmt.Println("Goodbye!")
		return
	default:
		fmt.Println("Invalid choice")
		g.Menu()
	}
}

// Start the game
func (g *Game) StartGame() {

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
func (g *Game) CheckFirstLetter() bool {

	if strings.ToUpper(string(rune(g.WordInput[0]))) == strings.ToUpper(g.RandomLetter) {
		fmt.Println("bravo")
		return true
	} else {
		fmt.Println("ca commence pas part ca")
		return false
	}

}

func (g *Game) SetTimer() {
	fmt.Println("how long did the tour last (in secondes)?")
	_, _ = fmt.Scanln(&g.Timer)

	if g.Timer <= 0 {
		fmt.Println("The time must be over 0 seconds, otherwise the game will not start.")
		g.SetTimer()
	}

	if g.Timer > 0 {
		g.Timer = g.Timer
	}
}

func (g *Game) SetRound() {
	fmt.Println("how many rounds do you want")
	_, _ = fmt.Scanln(&g.Round)

	if g.Round <= 0 {
		fmt.Println("not enough rounds")
		g.SetRound()
	}

	if g.Round > 0 {
		g.Round = g.Round
	}
}
