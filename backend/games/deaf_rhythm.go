package games

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"maps"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	Lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/tsirysndr/go-deezer"
)

var (
	playlist      []deezer.Track
	previousSongs []string
	playedSongs   []song
	tmpl          *template.Template
	currentSong   song
	songTimer     *time.Timer
	clients       = make(map[*websocket.Conn]player)
	clientsMutex  sync.Mutex
	upgrader      = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	currentRound  int
	maxRounds     = 2
	timerDuration = 30
	timerEndTime  int64
	haveGuessed   []string
	guessedMutex  sync.Mutex
	finalResults  []map[string]any
	gameEnded     bool
)

type player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Message struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type song struct {
	songName   string
	artistName string
	lyrics     []string
}

func DeafRhythmHandler(r *mux.Router) {
	r.HandleFunc("/deaf/ws", handleWS)
	r.HandleFunc("/scoreboard", showScoreboard).Methods("GET")

	r.HandleFunc("/deaf-rhythm/game/deaf-rhythm.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "../../frontend/deaf-rhythm/game/deaf-rhythm.html")
		start(w, r)
	})
	r.HandleFunc("/deaf-rhythm/game/deaf-rhythm.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "../../frontend/deaf-rhythm/game/deaf-rhythm.js")
	})
	r.HandleFunc("/deaf-rhythm/game/deaf-rhythm.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "../../frontend/deaf-rhythm/game/deaf-rhythm.css")
	})

	r.HandleFunc("/score.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "../../frontend/deaf-rhythm/game/deaf-rhythm.html")
	})
	r.HandleFunc("/score.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./score.css")
	})
	r.HandleFunc("/score.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "./score.js")
	})
}

func start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("START")

	gameEnded = false
	getPlaylist()
	if currentRound >= maxRounds {
		http.Redirect(w, r, "/scoreboard", http.StatusSeeOther)
		return
	}
	if currentSong.songName == "" || len(currentSong.lyrics) == 0 {
		getTrack()
		getLyrics()
		if len(clients) > 0 {
			startTimer()
		} else {
			fmt.Println("No clients connected yet, deferring timer start")
		}
	}
}

func showScoreboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SCOREBOARD")
	playersList := make([]map[string]any, 0)

	if gameEnded && len(finalResults) > 0 {
		playersList = finalResults
	} else {
		clientsMutex.Lock()
		for _, player := range clients {
			playersList = append(playersList, map[string]any{
				"id":    player.ID,
				"name":  player.Name,
				"score": player.Score,
			})
		}
		clientsMutex.Unlock()
	}
	data := map[string]any{
		"players": playersList,
		"game":    "Deaf Rhythm",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshalling data:", err)
		http.Error(w, "Error marshalling data", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "score.html", template.JS(jsonData))
	if err != nil {
		log.Println("Error executing scoreboard template:", err)
		http.Error(w, "Error rendering scoreboard", http.StatusInternalServerError)
	}
	fmt.Println("Scoreboard rendered")
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
	clientsCount := len(clients)
	clientsMutex.Unlock()

	sendToClient(conn, Message{
		Type: "player_assigned",
		Content: map[string]string{
			"id":   playerID,
			"name": "Anonymous Player",
		},
	})
	sendPlayerList()

	if currentSong.songName == "" {
		getTrack()
		getLyrics()
	}
	sendToClient(conn, Message{
		Type:    "lyrics",
		Content: currentSong.lyrics,
	})
	if clientsCount == 1 && songTimer == nil {
		startTimer()
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

	// fmt.Printf("Processing message of type '%s' from client %s\n", message.Type, conn.RemoteAddr())
	// fmt.Printf("Message content raw: %+v\n", message.Content)

	switch message.Type {
	case "guess":
		guess, ok := message.Content.(string)
		if !ok {
			log.Printf("Invalid content type for 'guess': %T\n", message.Content)
			return
		}
		guessHandler(conn, guess)
	case "register_player":
		registerPlayer(message, conn)
	case "request_lyrics":
		requestLyrics(conn)
	default:
		fmt.Printf("Unknown message type: %s\n", message.Type)
	}
}

func requestLyrics(conn *websocket.Conn) {
	fmt.Println("REQUEST LYRICS")

	fmt.Println("Lyrics request received, preparing response")
	fmt.Printf("Current song state - Name: '%s', Artist: '%s', Lyrics length: %d\n",
		currentSong.songName, currentSong.artistName, len(currentSong.lyrics))

	if len(currentSong.lyrics) > 0 {
		fmt.Println("Sending lyrics to client:", conn.RemoteAddr())

		fmt.Println("Lyrics content being sent:")
		for i, line := range currentSong.lyrics {
			fmt.Printf("  [%d]: %s\n", i, line)
		}

		sendToClient(conn, Message{
			Type:    "lyrics",
			Content: currentSong.lyrics,
		})
	} else {
		fmt.Println("No lyrics available yet")
		sendToClient(conn, Message{
			Type:    "lyrics",
			Content: []string{"No lyrics available yet. Waiting for next song..."},
		})
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

			if len(currentSong.lyrics) > 0 {
				fmt.Println("Sending initial lyrics to new client")
				sendToClient(conn, Message{
					Type:    "lyrics",
					Content: currentSong.lyrics,
				})
				sendToClient(conn, Message{
					Type: "timer_start",
					Content: map[string]any{
						"timerEnd":  timerEndTime,
						"round":     currentRound,
						"maxRounds": maxRounds,
					},
				})
			} else {
				sendToClient(conn, Message{
					Type:    "lyrics",
					Content: []string{"No lyrics available yet. Waiting for next song..."},
				})
			}
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

		if msg.Type == "lyrics" {
			lyrics, ok := msg.Content.([]string)
			if ok {
				fmt.Printf("Sent %d lyrics lines to %s: %v\n", len(lyrics), client.RemoteAddr(), lyrics)
			}
		}
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

	fmt.Println("exiting broadcastMessage")
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
	fmt.Println("Exiting sendPlayerList function")
}

func guessHandler(conn *websocket.Conn, guess string) {
	fmt.Println("GUESS\n", guess)

	clientsMutex.Lock()
	player, exists := clients[conn]
	clientsMutex.Unlock()

	if !exists {
		log.Println("Player not found for connection")
		return
	}
	points := calculatePoints()
	isCorrect := cleanText(currentSong.songName) == cleanText(guess)
	if isCorrect {
		guessedMutex.Lock()
		haveGuessed = append(haveGuessed, player.ID)
		guessedMutex.Unlock()
		player.Score += points
		clientsMutex.Lock()
		clients[conn] = player
		clientsMutex.Unlock()
		fmt.Printf("Player %s score increased to %d\n", player.Name, player.Score)
		broadcastMessage(Message{
			Type: "player_guessed",
			Content: map[string]any{
				"playerID":   player.ID,
				"playerName": player.Name,
				"score":      player.Score,
			},
		})
	}
	sendToClient(conn, Message{
		Type: "guess_result",
		Content: map[string]any{
			"correct": isCorrect,
			"score":   player.Score,
			"points":  points,
		},
	})
	if checkIfAllGuessed() {
		skipSong()
	}
}

func cleanText(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	text = strings.NewReplacer(
		".", "", ",", "", "!", "", "?", "",
		"'", "", "\"", "", "-", " ",
		"_", " ", "(", "", ")", "",
	).Replace(text)
	return strings.Join(strings.Fields(text), " ")
}

func calculatePoints() int {
	currentTime := time.Now().UnixMilli()

	if timerEndTime == 0 {
		return 10
	}
	timeRemaining := timerEndTime - currentTime
	if timeRemaining <= 0 {
		return 1
	}
	msDuration := timerDuration * 1000
	points := 1 + int((float64(timeRemaining)/float64(msDuration))*9)

	return min(10, points)
}

func checkIfAllGuessed() bool {
	clientsMutex.Lock()
	guessedMutex.Lock()
	defer clientsMutex.Unlock()
	defer guessedMutex.Unlock()

	if len(haveGuessed) == 0 || len(clients) == 0 {
		return false
	}
	for _, client := range clients {
		if !slices.Contains(haveGuessed, client.ID) {
			return false
		}
	}
	return true
}

func skipSong() {
	fmt.Println("SKIP SONG")

	if songTimer != nil {
		songTimer.Stop()
	}
	broadcastMessage(Message{
		Type: "timer_end",
		Content: map[string]any{
			"songName":   currentSong.songName,
			"artistName": currentSong.artistName,
			"skipped":    true,
		},
	})
	if currentRound >= maxRounds {
		endGame()
		return
	}
	getTrack()
	getLyrics()
	startTimer()
}

func startTimer() {
	fmt.Println("START TIMER")

	if songTimer != nil {
		songTimer.Stop()
	}
	currentRound++
	guessedMutex.Lock()
	haveGuessed = []string{}
	guessedMutex.Unlock()
	timerEndTime = time.Now().Add(time.Duration(timerDuration) * time.Second).UnixMilli()
	songTimer = time.AfterFunc(time.Duration(timerDuration)*time.Second, timerEnd)
	fmt.Println(maxRounds)
	broadcastMessage(Message{
		Type: "timer_start",
		Content: map[string]any{
			"timerEnd":  timerEndTime,
			"round":     currentRound,
			"maxRounds": maxRounds,
		},
	})
	broadcastMessage(Message{
		Type:    "lyrics",
		Content: currentSong.lyrics,
	})
}

func timerEnd() {
	fmt.Println("TIMER END")

	broadcastMessage(Message{
		Type: "timer_end",
		Content: map[string]any{
			"songName":   currentSong.songName,
			"artistName": currentSong.artistName,
		},
	})
	if currentRound >= maxRounds {
		endGame()
		return
	}
	getTrack()
	getLyrics()
	startTimer()
}

func endGame() {
	gameEnded = true
	finalResults = make([]map[string]any, 0)
	clientsMutex.Lock()
	for _, player := range clients {
		finalResults = append(finalResults, map[string]any{
			"id":    player.ID,
			"name":  player.Name,
			"score": player.Score,
		})
	}
	clientsMutex.Unlock()

	currentRound = 0
	previousSongs = []string{}
	guessedMutex.Lock()
	haveGuessed = []string{}
	guessedMutex.Unlock()
	broadcastMessage(Message{
		Type:    "game_over",
		Content: "/scoreboard",
	})
	fmt.Println("Game ended")
	for i, player := range finalResults {
		fmt.Printf("Player %d: , Name: %s, Score: %d\n", i+1, player["name"], player["score"])
	}
}

func getPlaylist() {
	fmt.Println("GET PLAYLIST")

	client := deezer.NewClient()
	tracks, err := client.Playlist.GetTracks("1565553361")
	if err != nil {
		log.Printf("Failed to get playlist: %v", err)
		return
	}
	playlist = tracks.Data
	rand.Shuffle(len(playlist), func(i, j int) {
		playlist[i], playlist[j] = playlist[j], playlist[i]
	})
	fmt.Println("Playlist loaded successfully")
}

func getTrack() {
	fmt.Println("GET TRACK")

	if len(playlist) == 0 {
		getPlaylist()
	}
	randomTrackIndex := rand.IntN(len(playlist))
	randomTrack := playlist[randomTrackIndex]
	if slices.Contains(previousSongs, randomTrack.Title) {
		fmt.Println("song already played")
		getTrack()
		return
	}
	currentSong.songName = randomTrack.TitleShort
	currentSong.artistName = randomTrack.Artist.Name

	previousSongs = append(previousSongs, randomTrack.TitleShort)
	playlist = slices.Delete(playlist, randomTrackIndex, randomTrackIndex+1)
}

func getLyrics() {
	fmt.Println("GET LYRICS")

	if currentSong.songName == "" {
		log.Println("Cannot get lyrics - song name is empty")
		return
	}
	l := Lyrics.New()
	var lyricsStr string
	var err error
	for i := 0; i < 10; i++ {
		lyricsStr, err = l.Search(currentSong.artistName, currentSong.songName)

		if err != nil {
			fmt.Printf("Error searching for lyrics: %v\n", err)
			getTrack()
			continue
		} else if lyricsStr == "" || strings.Contains(lyricsStr, "We do not have the lyrics for") {
			fmt.Println("No lyrics found for:", currentSong.songName)
			getTrack()
			continue
		} else {
			break
		}
	}
	if err != nil || lyricsStr == "" {
		log.Println("Failed to get lyrics after multiple attempts")
		currentSong.lyrics = []string{"Error: no lyrics found"}
		return
	}

	lyrics := strings.Split(lyricsStr, "\n")

	var filteredLyrics []string
	for _, line := range lyrics {
		if strings.TrimSpace(line) != "" {
			filteredLyrics = append(filteredLyrics, strings.TrimSpace(line))
		}
	}
	size := 10
	if len(filteredLyrics) < size {
		currentSong.lyrics = filteredLyrics
	} else {
		start := rand.IntN(len(filteredLyrics) - size)
		currentSong.lyrics = filteredLyrics[start : start+size]
	}
	playedSongs = append(playedSongs, currentSong)

	fmt.Println("Found lyrics for:", currentSong.songName, "by", currentSong.artistName)
	for _, line := range currentSong.lyrics {
		fmt.Println("  ", line)
	}
}
