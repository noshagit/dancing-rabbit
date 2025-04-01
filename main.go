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
	Responses map[string]string
}

// Struct of the Game
type Game struct {
	RandomLetter        string   // Random letter choose by the computer
	Round               int      // Round of the game
	LetterAlreadyChoose []string // Letter already choose by the computer
	Timer               int      // Timer of the game
	Point               int
	Players             []Player
    AllResponses        []string
	WordInput           string
}

// Function main
func main() {
	game := Game{}

	game.Menu()
}

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

func (g *Game) StartGame() {
    g.SetPlayers()
	g.SetTimer()
	g.SetRound()
	g.RoundPlayer()
	
}

func (g *Game) SetPlayers() {
    println("Enter a number of players")
    var NumberOfPlayers int
    _, _ = fmt.Scanln(&NumberOfPlayers)

    for i := 0; i < NumberOfPlayers; i++ {
        println("Enter the name of the player")
        var PlayerName string
        _, _ = fmt.Scanln(&PlayerName)

        g.Players = append(g.Players, Player{Name: PlayerName, Responses: make(map[string]string)})
    }
}

func (g *Game) RoundPlayer() string {
    cat := []string{"Album", "Artist", "MusicGroup", "Song", "MusicalGenre"}
	
	Time := make(chan bool)

    for g.Round > 0{

		letter := g.RandomLetterFunc()
    	println("Letter chosen:", letter)

		go func() {
			time.Sleep(time.Duration(g.Timer) * time.Second)
			Time <- true
			g.Round--
		}()

			go func() {
         	for i := 0; i < len(g.Players); i++ {
            fmt.Println("---------------------")
            fmt.Println("Player", i+1, ":", g.Players[i].Name)

            for _, category := range cat {
                fmt.Println("---------------------")
                fmt.Println("Category:", category)
                fmt.Println("Enter a word that starts with the letter", g.RandomLetter)
                fmt.Println("Enter your word:")

        
                _, _ = fmt.Scanln(&g.WordInput)
                g.Players[i].Responses[category] = g.WordInput
                g.AllResponses = append(g.AllResponses, g.WordInput)
				g.CheckFirstLetter()
				
            }
			if g.Round >= 0 {
				g.Round--
				Time <- true
				}
        }
    }()

	if g.Round == 0 {
		g.Round--
		Time <- true
		}
	<- Time
	println("Finish")

}
return "finish"
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
func (g *Game)CheckFirstLetter() bool {
	
	if strings.ToUpper(string(rune(g.WordInput[0]))) == strings.ToUpper(g.RandomLetter) {
		fmt.Println("bravo")
		return true
	}else {
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

func (g *Game) SetRound(){
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