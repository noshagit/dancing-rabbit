package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"text/template"
)

var lobbyPlayer []Player
var currentGame *Game

func fileServer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/start", http.StatusSeeOther)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		player := &currentGame.Players[currentGame.CurrentTurn]

		for _, category := range currentGame.Categories {
			answer := r.FormValue(category)
			player.Answers[category] = answer
		}

		log.Println("Réponses du joueur", player.Name, ":", player.Answers)

		currentGame.CurrentTurn++

		if currentGame.CurrentTurn >= len(currentGame.Players) {
			http.Redirect(w, r, "/vote", http.StatusSeeOther)
			return
		}
	}

	tmpl, err := template.ParseFiles("static/game.html")
	if err != nil {
		log.Println("Erreur lors du chargement du template game.html:", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, currentGame)
	if err != nil {
		log.Println("Erreur lors de l'exécution du template:", err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		action := r.FormValue("action")

		switch action {
		case "add":
			name := r.FormValue("name")
			if name != "" {
				addPlayer(&lobbyPlayer, name)
				log.Println("Joueur ajouté :", name)
			}
		case "remove":
			name := r.FormValue("name")
			if name != "" {
				removePlayer(&lobbyPlayer, name)
				log.Println("Joueur retiré :", name)
			}
		case "start":
			if len(lobbyPlayer) == 0 {
				log.Println("Aucun joueur n'a été ajouté")
				http.Error(w, "Aucun joueur n'a été ajouté", http.StatusBadRequest)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			roundStr := r.FormValue("rounds")
			maxRounds := 3
			if roundStr != "" {
				if parsed, err := strconv.Atoi(roundStr); err == nil && parsed > 0 {
					maxRounds = parsed
				}
			}

			game := &Game{
				Categories:   []string{"Animaux", "Ville", "Plantes", "Pays"},
				Players:      lobbyPlayer,
				RandomLetter: "",
				CurrentTurn:  0,
				CurrentRound: 1,
				MaxRounds:    maxRounds,
			}

			game.RandomLetter = game.RandomLetterFunc()
			currentGame = game
			lobbyPlayer = []Player{}
			log.Println("Jeu commencé, redirection vers /play")
			http.Redirect(w, r, "/play", http.StatusSeeOther)
			return
		}
	}

	tmpl, err := template.ParseFiles("static/start.html")
	if err != nil {
		log.Println("Erreur lors du chargement du template start.html:", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Players []Player
	}{
		Players: lobbyPlayer,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Erreur lors de l'exécution du template start.html:", err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		for i := range currentGame.Players {
			for category := range currentGame.Players[i].Answers {
				key := currentGame.Players[i].Name + "_" + category
				if r.FormValue(key) == "on" {
					currentGame.Players[i].Vote[category] = true
				} else {
					currentGame.Players[i].Vote[category] = false
				}
			}
		}

		currentGame.CurrentRound++
		if currentGame.CurrentRound > currentGame.MaxRounds {
			log.Println("Fin du jeu, redirection vers la page d'accueil")
			http.Redirect(w, r, "/end", http.StatusSeeOther)
		} else {
			currentGame.CurrentTurn = 0
			currentGame.RandomLetter = currentGame.RandomLetterFunc()

			for i := range currentGame.Players {
				currentGame.Players[i].Answers = make(map[string]string)
				currentGame.Players[i].Vote = make(map[string]bool)
			}
			log.Println("Début du tour", currentGame.CurrentRound, "avec la lettre", currentGame.RandomLetter)
			http.Redirect(w, r, "/play", http.StatusSeeOther)
		}
		return
	}

	tmpl, err := template.ParseFiles("static/vote.html")
	if err != nil {
		log.Println("Erreur lors du chargement du template vote.html:", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, currentGame)
	if err != nil {
		log.Println("Erreur lors de l'exécution du template vote.html:", err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}
}

func handleEnd(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sortedPlayers := make([]Player, len(currentGame.Players))
	copy(sortedPlayers, currentGame.Players)

	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i].CalculateScore() > sortedPlayers[j].CalculateScore()
	})

	tmpl, err := template.ParseFiles("static/end.html")
	if err != nil {
		log.Println("Erreur lors du chargement du template end.html:", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, sortedPlayers)
	if err != nil {
		log.Println("Erreur lors de l'exécution du template end.html:", err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}

}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", fileServer)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/play", handleGame)
	http.HandleFunc("/vote", handleVote)
	http.HandleFunc("/end", handleEnd)
	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
