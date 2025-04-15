package main

import (
	"math/rand"
)

type Game struct {
	Categories          []string
	RandomLetter        string
	LetterAlreadyChoose []string
	Players             []Player
	CurrentPlayer       Player
	CurrentTurn         int
	CurrentRound        int
	MaxRounds           int
}

type Player struct {
	Name    string
	Answers map[string]string
	Vote    map[string]bool
}

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

func addPlayer(players *[]Player, name string) {
	*players = append(*players, Player{
		Name:    name,
		Answers: make(map[string]string),
		Vote:    make(map[string]bool),
	})
}

func removePlayer(players *[]Player, name string) {
	for i, player := range *players {
		if player.Name == name {
			*players = append((*players)[:i], (*players)[i+1:]...)
			break
		}
	}
}

func (p *Player) CalculateScore() int {
	score := 0
	for _, valid := range p.Vote {
		if valid {
			score++
		}
	}
	return score
}
