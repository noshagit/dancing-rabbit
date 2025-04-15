package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/tsirysndr/go-deezer"
)

var (
	client         = deezer.NewClient()
	shuffledTracks []deezer.Track
	previousSongs  []string
	currentSong    deezer.Track
	playerGuesses  = make(map[string]string)
	gameState      = GameState{}
	mutex          = &sync.Mutex{}
	totalRounds    = 1
	currentRound   = 0
	gameStarted    = false
	roundEndTime   time.Time
	roundStartTime time.Time
)

type GameState struct {
	CurrentRound  int               `json:"currentRound"`
	TotalRounds   int               `json:"totalRounds"`
	TimeLeft      int               `json:"timeLeft"`
	CurrentSong   string            `json:"currentSong"`
	Players       []Player          `json:"players"`
	PlayerScores  map[string]int    `json:"playerScores"`
	PlayerAnswers map[string]string `json:"playerAnswers"`
	GameStarted   bool              `json:"gameStarted"`
	RoundActive   bool              `json:"roundActive"`
	PreviewURL    string            `json:"previewUrl"`
}

type Player struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Score  int    `json:"score"`
}

func BlindTestHandler(r *mux.Router) {
	r.HandleFunc("/blind-test/game/blind-test.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "/../../dancing-rabbit//frontend/blind-test/game/blind-test.html")
	}).Methods("GET")

	r.HandleFunc("/start", startGame).Methods("POST")
	r.HandleFunc("/guess", handleGuess).Methods("POST")
	r.HandleFunc("/state", getGameState).Methods("GET")
	r.HandleFunc("/next-round", nextRound).Methods("POST")

	r.HandleFunc("/blind-test/game/blind-test.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "/../../frontend/blind-test/game/blind-test.css")
	}).Methods("GET")

	r.HandleFunc("/blind-test/game/blind-test.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "/../../frontend/blind-test/game/blind-test.js")
	}).Methods("GET")

	initGameState()
}

func initGameState() {
	gameState = GameState{
		CurrentRound:  0,
		TotalRounds:   totalRounds,
		TimeLeft:      30,
		CurrentSong:   "",
		Players:       []Player{},
		PlayerScores:  make(map[string]int),
		PlayerAnswers: make(map[string]string),
		GameStarted:   false,
		RoundActive:   false,
		PreviewURL:    "",
	}
}

func startGame(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var players []string
	err := json.NewDecoder(r.Body).Decode(&players)
	if err != nil {
		http.Error(w, "Invalid player data", http.StatusBadRequest)
		return
	}

	gameState.Players = make([]Player, len(players))
	for i, name := range players {
		gameState.Players[i] = Player{
			Name:   name,
			Status: "waiting",
			Score:  0,
		}
		gameState.PlayerScores[name] = 0
	}

	gameStarted = true
	currentRound = 0
	previousSongs = []string{}
	shuffledTracks = loadAndShufflePlaylist()
	gameState.GameStarted = true

	startRound()

	json.NewEncoder(w).Encode(gameState)
}

func loadAndShufflePlaylist() []deezer.Track {
	playlist, err := client.Playlist.GetTracks("12739384001")
	if err != nil || len(playlist.Data) == 0 {
		log.Fatalf("Failed to load playlist: %v", err)
	}

	tracks := playlist.Data
	rand.Shuffle(len(tracks), func(i, j int) {
		tracks[i], tracks[j] = tracks[j], tracks[i]
	})
	return tracks
}

func startRound() {
	if currentRound >= totalRounds {
		endGame()
		return
	}

	currentRound++
	gameState.CurrentRound = currentRound
	getBlindTestTrack()

	gameState.CurrentSong = currentSong.Title
	gameState.PreviewURL = currentSong.Preview
	gameState.RoundActive = true
	gameState.TimeLeft = 30
	roundStartTime = time.Now()
	roundEndTime = roundStartTime.Add(30 * time.Second)

	for name := range gameState.PlayerScores {
		gameState.PlayerAnswers[name] = ""
	}

	for i := range gameState.Players {
		gameState.Players[i].Status = "playing"
	}

	go func() {
		time.Sleep(30 * time.Second)
		endRound()
	}()
}

func endRound() {
	mutex.Lock()
	defer mutex.Unlock()

	gameState.RoundActive = false

	for i, player := range gameState.Players {
		if gameState.PlayerAnswers[player.Name] == "" {
			gameState.Players[i].Status = "timeout"
		}
	}

	go func() {
		time.Sleep(10 * time.Second)
		nextRound(nil, nil)
	}()
}

func endGame() {
	gameState.GameStarted = false
	gameStarted = false
	gameState.RoundActive = false
	gameState.CurrentRound = totalRounds
}

func nextRound(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	startRound()

	if w != nil {
		json.NewEncoder(w).Encode(gameState)
	}
}

func getBlindTestTrack() {
	if len(shuffledTracks) == 0 {
		log.Println("No more songs left in shuffled list.")
		currentSong = deezer.Track{
			Title:   "No Song",
			Preview: "",
		}
		return
	}

	for len(shuffledTracks) > 0 {
		track := shuffledTracks[0]
		shuffledTracks = shuffledTracks[1:]

		if !slices.Contains(previousSongs, track.Title) {
			currentSong = track
			previousSongs = append(previousSongs, track.Title)
			log.Printf("Selected track: %+v", currentSong)
			return
		}
	}

	log.Println("Ran out of unique songs, reloading playlist.")
	shuffledTracks = loadAndShufflePlaylist()
	getBlindTestTrack()
}

func handleGuess(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	if !gameState.RoundActive {
		http.Error(w, "Round is not active", http.StatusBadRequest)
		return
	}

	type Guess struct {
		Player string `json:"player"`
		Answer string `json:"answer"`
	}

	var guess Guess
	err := json.NewDecoder(r.Body).Decode(&guess)
	if err != nil {
		http.Error(w, "Invalid guess data", http.StatusBadRequest)
		return
	}

	found := false
	for i, player := range gameState.Players {
		if player.Name == guess.Player {
			found = true
			if gameState.PlayerAnswers[player.Name] == "" {
				gameState.PlayerAnswers[player.Name] = guess.Answer

				if strings.ToLower(guess.Answer) == strings.ToLower(currentSong.Title) {
					timeElapsed := time.Since(roundStartTime)
					score := calculateScore(timeElapsed)
					gameState.PlayerScores[player.Name] += score
					gameState.Players[i].Score = gameState.PlayerScores[player.Name]
					gameState.Players[i].Status = "correct"
				} else {
					gameState.Players[i].Status = "wrong"
				}
			}
			break
		}
	}

	if !found {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(gameState)
}

func calculateScore(timeElapsed time.Duration) int {
	secondsElapsed := int(timeElapsed.Seconds())
	if secondsElapsed > 30 {
		secondsElapsed = 30
	}
	score := 100 - (secondsElapsed * 3)
	if score < 10 {
		score = 10
	}
	return score
}

func getGameState(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameState.RoundActive {
		gameState.TimeLeft = int(roundEndTime.Sub(time.Now()).Seconds())
		if gameState.TimeLeft < 0 {
			gameState.TimeLeft = 0
		}
	}

	json.NewEncoder(w).Encode(gameState)
}
