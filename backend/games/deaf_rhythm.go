package main

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
	playlist      *spotify.FullPlaylist
	previousSongs []string
	playedSongs    []song
	tmpl          *template.Template
	currentSong   song
	playerGuesses map[string]string
	songTimer     *time.Timer
)

type song struct {
	songName   string
	artistName string
	lyrics     []string
}

func main() {
	tmpl = template.Must(template.ParseFiles("deaf_rhythm.html"))
	playerGuesses = make(map[string]string)

	getPlaylist()

	r := mux.NewRouter()

	r.HandleFunc("/", start).Methods("GET")
	r.HandleFunc("/guess", guessHandler).Methods("POST")

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

func start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("START")

	if len(playedSongs) == 10 {
		// TODO end game and vote
	}

	if currentSong.songName == "" {
		getTrack()
		getLyrics()
		resetSongTimer(w, r)
	}

	for _, line := range currentSong.lyrics {
		fmt.Println(line)
	}
	fmt.Printf("\n%s ; by %s\n\n", currentSong.songName, currentSong.artistName)
	println("tmpl execute")

	err := tmpl.Execute(w, currentSong.lyrics)
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
	playerGuesses[currentSong.songName] = guess

	currentSong.songName = ""

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetSongTimer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RESET TIMER")

	if songTimer != nil {
		songTimer.Stop()
	}

	songTimer = time.AfterFunc(10*time.Second, func() {
		fmt.Println("timer expired")
		playerGuesses[currentSong.songName] = ""
		getTrack()
		getLyrics()
		// TODO refresh page to show the new lyrics
	})
}

func getLyrics() {
	fmt.Println("GET LYRICS")

	for i := 0; i < 10; i++ {
		l := Lyrics.New()
		lyricsStr, err := l.Search(currentSong.artistName, currentSong.songName)

		if err == nil && lyricsStr != "" && !strings.Contains(lyricsStr, "We do not have the lyrics for") {
			lyrics := strings.Split(lyricsStr, "\n")
			size := 10
			if len(lyrics) < size {
				currentSong.lyrics = lyrics
				return
			}
			start := rand.IntN(len(lyrics) - size)
			currentSong.lyrics = lyrics[start : start+size]
			playedSongs = append(playedSongs, currentSong)
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
	if slices.Contains(previousSongs, randomTrack.Name) {
		fmt.Println("song already played")
		getTrack()
		return
	}
	previousSongs = append(previousSongs, randomTrack.Name)
	currentSong.songName = randomTrack.Name
	currentSong.artistName = randomTrack.Artists[0].Name
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
