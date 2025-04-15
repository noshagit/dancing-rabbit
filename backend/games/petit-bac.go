package main

// package games

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"maps"
	"math/rand/v2"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	tmpl     *template.Template
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	pastLetters   []rune
	currentLetter rune
	timer         *time.Timer
	timerEndTime  int64
	timerDuration = 30
	currentRound  int
	maxRounds     int
	clientsMutex  sync.Mutex
	clients       = make(map[*websocket.Conn]player)
	categories    = []string{
		"Artist",
		"Genre",
		"Song",
	}
	voteMode    = false
	votingMutex sync.Mutex
	playerVotes = make(map[string]map[string]map[string]bool) // playerID -> targetPlayerID -> category -> vote
	gameData    = make(map[string]player)
)

type player struct {
	ID     string
	Name   string
	Score  int
	Rounds []round
}

type Message struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type round struct {
	Letter  string
	Answers map[string]string
}

func main() {
	tmpl = template.Must(template.ParseFiles(
		"../../frontend/petit-bac/petit-bac-menu.html",
		"../../frontend/petit-bac/game/petit-bac.html",
		"../../frontend/vote/vote.html",
		"../../frontend/score/score.html",
		"../../frontend/main-menu/menu.html"))

	r := mux.NewRouter()
	r.HandleFunc("/", menu).Methods("GET")
	r.HandleFunc("/ws", handleWS).Methods("GET")
	r.HandleFunc("/game", game).Methods("GET")
	r.HandleFunc("/vote", vote).Methods("GET")
	r.HandleFunc("/scoreboard", scoreboard).Methods("GET")

	r.HandleFunc("/main-menu", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "menu.html", nil)
	}).Methods("GET")
	r.HandleFunc("/menu.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/main-menu/menu.js")
	}).Methods("GET")
	r.HandleFunc("/menu.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/main-menu/menu.css")
	}).Methods("GET")

	r.HandleFunc("/petit-bac-menu.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/petit-bac/petit-bac-menu.js")
	}).Methods("GET")
	r.HandleFunc("/petit-bac-menu.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/petit-bac/petit-bac-menu.css")
	}).Methods("GET")

	r.HandleFunc("/petit-bac.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/petit-bac/game/petit-bac.js")
	}).Methods("GET")
	r.HandleFunc("/petit-bac.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/petit-bac/game/petit-bac.css")
	}).Methods("GET")

	r.HandleFunc("/vote.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/vote/vote.js")
	}).Methods("GET")
	r.HandleFunc("/vote.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/vote/vote.css")
	}).Methods("GET")

	r.HandleFunc("/score.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/score/score.js")
	}).Methods("GET")
	r.HandleFunc("/score.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/score/score.css")
	}).Methods("GET")

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}

func menu(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "petit-bac-menu.html", nil)
}

func game(w http.ResponseWriter, r *http.Request) {
	rounds := r.URL.Query().Get("rounds")
	duration := r.URL.Query().Get("duration")

	fmt.Printf("Received parameters - rounds: %s, duration: %s\n", rounds, duration)

	if rounds != "" {
		var err error
		maxRounds, err = parseInt(rounds, 5)
		if err != nil {
			log.Printf("Error parsing rounds parameter '%s': %v", rounds, err)
			maxRounds = 5
		}
	} else {
		maxRounds = 5
	}
	if duration != "" {
		var err error
		timerDuration, err = parseInt(duration, 30)
		if err != nil {
			log.Printf("Error parsing duration parameter '%s': %v", duration, err)
			timerDuration = 30
		}
	} else {
		timerDuration = 30
	}
	fmt.Printf("Game started with rounds: %d, duration: %d seconds\n", maxRounds, timerDuration)

	currentRound = 0
	pastLetters = []rune{}
	currentLetter = 0

	clientsMutex.Lock()
	for conn, p := range clients {
		p.Rounds = nil
		p.Score = 0
		clients[conn] = p
	}
	clientsMutex.Unlock()

	clientsMutex.Lock()
	clientsCount := len(clients)
	clientsMutex.Unlock()

	if clientsCount > 0 && timer == nil {
		startTimer()
	}
	if currentLetter == 0 {
		randomLetter()
	}
	gameParams := map[string]interface{}{
		"letter":     string(currentLetter),
		"categories": categories,
		"rounds":     maxRounds,
		"duration":   timerDuration,
	}
	gameParamsJSON, err := json.Marshal(gameParams)
	if err != nil {
		log.Printf("Error marshalling game params: %v", err)
		gameParamsJSON = []byte(`{"letter":"A","categories":["Artist","Genre","Song"],"rounds":5,"duration":30}`)
	}
	data := map[string]interface{}{
		"params":     gameParams,
		"paramsJSON": template.JS(gameParamsJSON),
	}
	tmpl.ExecuteTemplate(w, "petit-bac.html", data)
}

func parseInt(value string, defaultVal int) (int, error) {
	if value == "" {
		return defaultVal, nil
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal, err
	}

	return val, nil
}

func sendGameParams(conn *websocket.Conn) {
	sendToClient(conn, Message{
		Type: "game_params",
		Content: map[string]any{
			"letter":     string(currentLetter),
			"categories": categories,
			"rounds":     maxRounds,
			"duration":   timerDuration,
		},
	})
}

func vote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("VOTE")
	voteMode = true
	clientsMutex.Lock()
	for _, p := range clients {
		fmt.Printf("Player %s has %d rounds\n", p.Name, len(p.Rounds))
		if len(p.Rounds) > 0 {
			fmt.Printf("Last round answers: %v\n", p.Rounds[len(p.Rounds)-1].Answers)
		}
	}
	clientsMutex.Unlock()
	tmpl.ExecuteTemplate(w, "vote.html", nil)
}

func scoreboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SCOREBOARD")

	finalScores := calculateFinalScores()
	tmpl.ExecuteTemplate(w, "scoreboard.html", finalScores)
}

func calculateFinalScores() []map[string]interface{} {
	scores := []map[string]interface{}{}

	for playerID, player := range gameData {
		playerScore := map[string]interface{}{
			"id":    playerID,
			"name":  player.Name,
			"score": player.Score,
		}
		scores = append(scores, playerScore)
	}
	return scores
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WS connection request from:", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	fmt.Println("WebSocket connection esplaylistlished with:", conn.RemoteAddr())
	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()

	playerID := fmt.Sprintf("player-%d", time.Now().UnixNano())

	clientsMutex.Lock()
	clients[conn] = player{
		ID:   playerID,
		Name: "Anonymous Player",
	}
	clientsMutex.Unlock()

	sendToClient(conn, Message{
		Type: "player_assigned",
		Content: map[string]string{
			"id":   playerID,
			"name": "Anonymous Player",
		},
	})
	sendPlayerList()
	if !voteMode {
		sendGameParams(conn)
	}
	if voteMode {
		sendVoteData(conn)
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", conn.RemoteAddr(), err)
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			sendPlayerList()
			break
		}
		processMessage(conn, msg)
	}
}

func processMessage(conn *websocket.Conn, msg []byte) {
	fmt.Printf("Raw message received from %s: %s\n", conn.RemoteAddr(), string(msg))

	var message Message
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Printf("Error unmarshalling message from %s: %v", conn.RemoteAddr(), err)
		return
	}

	fmt.Printf("Processing message of type '%s' from client %s\n", message.Type, conn.RemoteAddr())
	fmt.Printf("Message content raw: %+v\n", message.Content)

	switch message.Type {
	case "guess":
		fmt.Println("Received guess message")
		fmt.Printf("Guess content type: %T, value: %v\n", message.Content, message.Content)

		if contentMap, ok := message.Content.(map[string]any); ok {
			answers := make(map[string]string)
			for k, v := range contentMap {
				if strVal, ok := v.(string); ok {
					answers[k] = strVal
				} else {
					answers[k] = fmt.Sprintf("%v", v)
				}
			}
			submit(conn, answers)
		} else {
			log.Println("Invalid guess content format")
		}
	case "register_player":
		registerPlayer(message, conn)
	case "request_vote_data":
		sendVoteData(conn)
	case "voted":
		fmt.Println("Received vote message")
		voted(message)
	case "request_game_params":
		fmt.Println("Received game params request")
		sendGameParams(conn)
	case "start_game":
		fmt.Println("Received start game request")
		if content, ok := message.Content.(map[string]any); ok {
			rounds, _ := content["rounds"].(float64)
			duration, _ := content["duration"].(float64)
			broadcastGameStart(int(rounds), int(duration))
		}
	}
}

func broadcastGameStart(rounds, duration int) {
	fmt.Printf("Broadcasting game start with rounds: %d, duration: %d\n", rounds, duration)
	maxRounds = rounds
	timerDuration = duration
	currentRound = 0
	pastLetters = []rune{}
	currentLetter = 0
	voteMode = false

	clientsMutex.Lock()
	for conn, p := range clients {
		p.Rounds = nil
		p.Score = 0
		clients[conn] = p
	}
	clientsMutex.Unlock()

	broadcastMessage(Message{
		Type: "game_start",
		Content: map[string]any{
			"rounds":   rounds,
			"duration": duration,
		},
	})

	time.AfterFunc(500*time.Millisecond, func() {
		if currentLetter == 0 {
			randomLetter()
		}
		startTimer()
	})
}

func sendVoteData(conn *websocket.Conn) {
	voteData := make(map[string][]round)

	if len(gameData) > 0 {
		fmt.Println("Using stored game data for vote page")
		for id, p := range gameData {
			if len(p.Rounds) > 0 {
				voteData[id] = p.Rounds
			} else {
				voteData[id] = []round{}
			}
		}
	} else {
		clientsMutex.Lock()
		for _, p := range clients {
			if p.Rounds == nil {
				voteData[p.ID] = []round{}
			} else {
				voteData[p.ID] = p.Rounds
			}
		}
		clientsMutex.Unlock()
	}
	if len(voteData) == 0 || allEmptyRounds(voteData) {
		voteData["dummy-player"] = []round{
			{
				Letter:  "A",
				Answers: map[string]string{"Artist": "ACDC", "Genre": "Alternative", "Song": "All Star"},
			},
		}
	}

	sendToClient(conn, Message{
		Type:    "vote_data",
		Content: voteData,
	})
}

func allEmptyRounds(voteData map[string][]round) bool {
	for _, rounds := range voteData {
		if len(rounds) > 0 {
			return false
		}
	}
	return true
}

func registerPlayer(message Message, conn *websocket.Conn) {
	fmt.Println("REGISTER PLAYER")

	content, ok := message.Content.(map[string]any)
	if ok {
		name, ok2 := content["name"].(string)
		if ok2 {
			clientsMutex.Lock()
			client := clients[conn]
			client.Name = name

			clients[conn] = client
			clientsMutex.Unlock()

			sendToClient(conn, Message{
				Type:    "registered",
				Content: clients[conn],
			})
			sendPlayerList()
		}
	}
}

func sendToClient(client *websocket.Conn, msg Message) {
	fmt.Printf("Sending %s message to client %s\n", msg.Type, client.RemoteAddr())

	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("Error sending message to %s: %v", client.RemoteAddr(), err)
		client.Close()
		clientsMutex.Lock()
		delete(clients, client)
		clientsMutex.Unlock()
		sendPlayerList()
	} else {
		fmt.Printf("Successfully sent %s message to %s\n", msg.Type, client.RemoteAddr())
	}
}

func broadcastMessage(msg Message) {
	clientsMutex.Lock()
	clientsCopy := make(map[*websocket.Conn]player, len(clients))
	maps.Copy(clientsCopy, clients)
	clientsMutex.Unlock()

	fmt.Printf("Broadcasting message of type '%s' to %d clients\n", msg.Type, len(clientsCopy))

	for client := range clientsCopy {
		sendToClient(client, msg)
	}
}
func sendPlayerList() {
	playersList := make([]map[string]any, 0)

	clientsMutex.Lock()
	for _, player := range clients {
		fmt.Println(player.Name)
		playersList = append(playersList, map[string]any{
			"id":    player.ID,
			"name":  player.Name,
			"score": player.Score,
		})
	}
	clientsMutex.Unlock()

	fmt.Println("Player list created:", playersList)
	broadcastMessage(Message{
		Type:    "player_list",
		Content: playersList,
	})
}

func startTimer() {
	fmt.Println("START TIMER")

	if timer != nil {
		timer.Stop()
	}
	currentRound++
	timerEndTime = time.Now().Add(time.Duration(timerDuration) * time.Second).UnixMilli()
	timer = time.AfterFunc(time.Duration(timerDuration)*time.Second, timerEnd)

	fmt.Printf("Starting round %d/%d with duration %d seconds\n", currentRound, maxRounds, timerDuration)

	broadcastMessage(Message{
		Type: "timer_start",
		Content: map[string]any{
			"timerEnd":  timerEndTime,
			"round":     currentRound,
			"maxRounds": maxRounds,
		},
	})
}

func timerEnd() {
	fmt.Println("TIMER END")

	if currentRound >= maxRounds {
		endGame()
		return
	}
	broadcastMessage(Message{
		Type:    "round_reset",
		Content: string(currentLetter),
	})

	startTimer()
}

func endGame() {
	fmt.Println("END GAME")

	clientsMutex.Lock()
	for _, p := range clients {
		gameData[p.ID] = p
	}
	clientsMutex.Unlock()

	fmt.Println("Game data stored for voting:", len(gameData), "players")
	broadcastMessage(Message{
		Type:    "end_game",
		Content: nil,
	})
}

func voted(message Message) {
	content, ok := message.Content.(map[string]interface{})
	if !ok {
		log.Println("Invalid vote content format")
		return
	}
	playerID, okP := content["playerID"].(string)
	targetPlayerID, okT := content["targetPlayerID"].(string)
	category, okC := content["category"].(string)
	voteValue, okV := content["valid"].(bool)

	roundNum := 0
	if roundVal, ok := content["round"].(float64); ok {
		roundNum = int(roundVal) - 1
	}

	if !okP || !okT || !okC || !okV {
		log.Println("Missing required vote fields")
		return
	}

	votingMutex.Lock()
	if playerVotes[playerID] == nil {
		playerVotes[playerID] = make(map[string]map[string]bool)
	}
	roundCategoryKey := fmt.Sprintf("r%d_%s", roundNum, category)
	if playerVotes[playerID][targetPlayerID] == nil {
		playerVotes[playerID][targetPlayerID] = make(map[string]bool)
	}
	playerVotes[playerID][targetPlayerID][roundCategoryKey] = voteValue
	votingMutex.Unlock()

	calculateAndBroadcastScores()

	if isVotingComplete() {
		broadcastMessage(Message{
			Type:    "voting_complete",
			Content: nil,
		})
	}
}

func calculateAndBroadcastScores() {
	voteResults := make(map[string]map[string]map[string]int)

	clientsMutex.Lock()
	for _, p := range clients {
		voteResults[p.ID] = make(map[string]map[string]int)
		if p.Rounds != nil {
			for roundIndex, r := range p.Rounds {
				for category := range r.Answers {
					key := fmt.Sprintf("Round %d - %s", roundIndex+1, category)
					voteResults[p.ID][key] = map[string]int{
						"valid":   0,
						"invalid": 0,
					}
				}
			}
		}
	}
	clientsMutex.Unlock()
	votingMutex.Lock()
	for _, targetVotes := range playerVotes {
		for targetID, categoryVotes := range targetVotes {
			for roundCategoryKey, vote := range categoryVotes {
				parts := strings.Split(roundCategoryKey, "_")
				if len(parts) != 2 {
					continue
				}
				roundStr := strings.TrimPrefix(parts[0], "r")
				roundNum, err := strconv.Atoi(roundStr)
				if err != nil {
					continue
				}
				category := parts[1]
				readableKey := fmt.Sprintf("Round %d - %s", roundNum+1, category)
				if voteResults[targetID] == nil {
					voteResults[targetID] = make(map[string]map[string]int)
				}
				if voteResults[targetID][readableKey] == nil {
					voteResults[targetID][readableKey] = map[string]int{
						"valid":   0,
						"invalid": 0,
					}
				}
				if vote {
					voteResults[targetID][readableKey]["valid"]++
				} else {
					voteResults[targetID][readableKey]["invalid"]++
				}
			}
		}
	}
	votingMutex.Unlock()
	broadcastMessage(Message{
		Type:    "vote_results",
		Content: voteResults,
	})
	updatePlayerScores(voteResults)
}

func updatePlayerScores(voteResults map[string]map[string]map[string]int) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	scoreUpdates := make(map[string]int)
	for conn, p := range clients {
		scoreUpdates[p.ID] = 0

		if results, ok := voteResults[p.ID]; ok {
			for _, counts := range results {
				if counts["valid"] > counts["invalid"] {
					p.Score += 10
					scoreUpdates[p.ID] += 10
				}
			}
			clients[conn] = p
		}
	}
	for id, points := range scoreUpdates {
		if player, exists := gameData[id]; exists && points > 0 {
			player.Score += points
			gameData[id] = player
		}
	}
	sendPlayerList()
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

func submit(conn *websocket.Conn, answers map[string]string) {
	fmt.Println("GUESS HANDLER")
	fmt.Printf("Guess content type: %T, value: %v\n", answers, answers)

	clientsMutex.Lock()
	player := clients[conn]
	if player.Rounds == nil {
		player.Rounds = []round{}
	}

	player.Rounds = append(player.Rounds, round{
		Letter:  string(currentLetter),
		Answers: answers,
	})
	clients[conn] = player

	fmt.Printf("Player %s now has %d rounds\n", player.Name, len(player.Rounds))
	fmt.Printf("Last round answers: %v\n", player.Rounds[len(player.Rounds)-1].Answers)
	clientsMutex.Unlock()
	if timer != nil {
		timer.Stop()
	}
	broadcastMessage(Message{
		Type:    "player_submitted",
		Content: player.Name,
	})
	time.AfterFunc(2*time.Second, func() {
		randomLetter()
		timerEnd()
	})
}

func isVotingComplete() bool {
	votingMutex.Lock()
	defer votingMutex.Unlock()
	if len(gameData) == 0 || len(playerVotes) == 0 {
		return false
	}
	totalVotesCast := 0
	for _, targetVotes := range playerVotes {
		for _, categoryVotes := range targetVotes {
			totalVotesCast += len(categoryVotes)
		}
	}
	activeVoters := len(clients)
	if activeVoters == 0 {
		return false
	}
	expectedVotes := len(gameData) * len(categories) * activeVoters
	return totalVotesCast >= expectedVotes*3/4
}
