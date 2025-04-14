package main

// package games

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"math/rand/v2"
	"net/http"
	"slices"
	"sync"
	"text/template"
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
	categories    = []string{}
	timer         *time.Timer
	timerEndTime  int64
	timerDuration = 30
	currentRound  int
	maxRounds     int
	clientsMutex  sync.Mutex
	clients       = make(map[*websocket.Conn]player)
)

type player struct {
	ID      string
	Name    string
	Score   int
	Answers []string
}

type Message struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

func main() {
	tmpl = template.Must(template.ParseFiles(
		"../../frontend/petit-bac/petit-bac-menu.html",
		"../../frontend/petit-bac/game/petit-bac.html"))

	r := mux.NewRouter()
	r.HandleFunc("/", menu).Methods("GET")
	r.HandleFunc("/ws", handleWS).Methods("GET")
	r.HandleFunc("/petit-bac/game", game).Methods("GET")

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

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}

func menu(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "petit-bac-menu.html", nil)
}

func game(w http.ResponseWriter, r *http.Request) {
	clientsMutex.Lock()
	clientsCount := len(clients)
	clientsMutex.Unlock()

	if clientsCount == 1 && timer == nil {
		startTimer()
	}
	tmpl.ExecuteTemplate(w, "petit-bac.html", nil)
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
		guessHandler(conn, message.Content.(string))
	case "register_player":
		registerPlayer(message, conn)
	}
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
	fmt.Println(maxRounds)
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
	startTimer()
}

func endGame() {}

func randomLetter() {
	i := rand.IntN(26)
	currentLetter = rune('A' + i)
	if slices.Contains(pastLetters, currentLetter) {
		randomLetter()
	} else {
		pastLetters = append(pastLetters, currentLetter)
	}
}

func guessHandler(conn *websocket.Conn, guess string) {}
