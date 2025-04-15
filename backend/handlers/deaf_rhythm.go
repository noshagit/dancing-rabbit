package handlers

// package games

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/mux"
	Lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var (
	playlist                *spotify.FullPlaylist
	playedSongNames         []string
	playedSongs             []song
	tmpl                    *template.Template
	currentDeafRhythmSong   song
	deafRhythmPlayerGuesses map[string]string
	songTimer               *time.Timer
)

type song struct {
	songName   string
	artistName string
	lyrics     []string
}

func DeafRhythmHandler(r *mux.Router) {
	r.HandleFunc("/deaf-rhythm/game/deaf-rhythm.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/deaf-rhythm/game/deaf-rhythm.html")
	}).Methods("GET")

	r.HandleFunc("/deaf-rhythm/game/deaf-rhythm.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/deaf-rhythm/game/deaf-rhythm.css")
	}).Methods("GET")

	r.HandleFunc("/deaf-rhythm/game/deaf-rhythm.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "/home/ilian/dancing-rabbit/frontend/deaf-rhythm/game/deaf-rhythm.js")
	}).Methods("GET")

	r.HandleFunc("/start", start).Methods("POST")
	r.HandleFunc("/guess", guessHandler).Methods("POST")

	deafRhythmPlayerGuesses = make(map[string]string)

	getPlaylist()
}

func start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("START")

	if len(playedSongs) == 10 {
		// TODO end game and vote
	}

	if currentDeafRhythmSong.songName == "" {
		getTrack()
		getLyrics()
		//resetSongTimer(w, r)
	}

	for _, line := range currentDeafRhythmSong.lyrics {
		fmt.Println(line)
	}
	fmt.Printf("\n%s ; by %s\n\n", currentDeafRhythmSong.songName, currentDeafRhythmSong.artistName)
	println("tmpl execute")

	err := tmpl.Execute(w, currentDeafRhythmSong.lyrics)
	if err != nil {
		log.Fatal(err)
	}
}

func guessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GUESS")

	r.ParseForm()
	guess := r.Form.Get("user-input")
	guess = strings.TrimSpace(strings.ToLower(guess))
	fmt.Println(guess)
	deafRhythmPlayerGuesses[currentDeafRhythmSong.songName] = guess

	currentDeafRhythmSong.songName = ""

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*func resetSongTimer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RESET TIMER")

	if songTimer != nil {
		songTimer.Stop()
	}

	songTimer = time.AfterFunc(10*time.Second, func() {
		fmt.Println("timer expired")
		deafRhythmPlayerGuesses[currentDeafRhythmSong.songName] = ""
		getTrack()
		getLyrics()
		// TODO refresh page to show the new lyrics
	})
}*/

func getLyrics() {
	fmt.Println("GET LYRICS")

	for i := 0; i < 10; i++ {
		l := Lyrics.New()
		lyricsStr, err := l.Search(currentDeafRhythmSong.artistName, currentDeafRhythmSong.songName)

		if err == nil && lyricsStr != "" && !strings.Contains(lyricsStr, "We do not have the lyrics for") {
			lyrics := strings.Split(lyricsStr, "\n")
			size := 10
			if len(lyrics) < size {
				currentDeafRhythmSong.lyrics = lyrics
				return
			}
			start := rand.IntN(len(lyrics) - size)
			currentDeafRhythmSong.lyrics = lyrics[start : start+size]
			playedSongs = append(playedSongs, currentDeafRhythmSong)
			return
		}

		fmt.Println("Error getting lyrics, trying a different song...")
		getTrack()
	}

	log.Fatal("Failed to get lyrics after multiple attempts")
}

func getPlaylist() {
	fmt.Println("GET PLAYLIST")

	ctx := context.Background()
	token := &oauth2.Token{
		AccessToken: newToken(),
		TokenType:   "Bearer",
	}

	client := spotify.New(spotifyauth.New().Client(ctx, token))

	playlistID := spotify.ID("6i2Qd6OpeRBAzxfscNXeWp")                        // playlist cannot be private or made by spotify
	fields := spotify.Fields("tracks(total,items(track(name,artists(name)))") // only get nb of tracks, track name and artist name

	var err error
	playlist, err = client.GetPlaylist(ctx, playlistID, fields)
	if err != nil {
		log.Fatalf("Failed to get playlist: %v", err)
	}
}

func getTrack() {
	fmt.Println("GET TRACK")

	randomTrackIndex := rand.IntN(int(math.Min(float64(playlist.Tracks.Total), 100))) // limit to track 100 because of Spotify API limit
	randomTrack := playlist.Tracks.Tracks[randomTrackIndex].Track
	if slices.Contains(playedSongNames, randomTrack.Name) {
		fmt.Println("song already played")
		getTrack()
		return
	}
	playedSongNames = append(playedSongNames, randomTrack.Name)
	currentDeafRhythmSong.songName = randomTrack.Name
	currentDeafRhythmSong.artistName = randomTrack.Artists[0].Name
}

func newToken() string {
	fmt.Println("NEW TOKEN")

	params := url.Values{}
	params.Add("grant_type", `client_credentials`)
	params.Add("client_id", `bb69f85b5ee84285bad7f1c28cadaf14`)
	params.Add("client_secret", `0fac697e3ea9456793fb92c22fc7977d`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.AccessToken
}
