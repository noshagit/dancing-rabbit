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

	Lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var playlist *spotify.FullPlaylist
var previousSongs []string
var tmpl *template.Template
var currentSong GameState

type GameState struct {
	songName   string
	artistName string
	lyrics     []string
}

func main() {
	tmpl = template.Must(template.ParseFiles("deaf_rhythm.html"))

	getPlaylist()
	http.HandleFunc("/", start)
	http.HandleFunc("/guess", guessHandler)

	http.ListenAndServe(":8080", nil)
}

func start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("START")

	if currentSong.songName == "" || r.URL.Query().Get("new") == "true" {
		getTrack()
		getLyrics()
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
	if guess == strings.ToLower(currentSong.songName) {
		fmt.Println("correct")
		http.Redirect(w, r, "/?new=true", http.StatusSeeOther)
		return
	}
	fmt.Println("wrong")
	err := tmpl.Execute(w, currentSong.lyrics)
	if err != nil {
		log.Fatal(err)
	}
}

func getLyrics() {
	fmt.Println("GET LYRICS")

	l := Lyrics.New()
	lyricsStr, err := l.Search(currentSong.songName, currentSong.artistName)

	for i := 0; (err != nil || lyricsStr == "") && i < 10; i++ {
		fmt.Println("Error getting lyrics, retrying... ", err)
		getTrack()
		lyricsStr, err = l.Search(currentSong.songName, currentSong.artistName)
	}
	if err != nil || lyricsStr == "" {
		log.Fatal("Failed to get lyrics")
	}
	if strings.Contains(lyricsStr, "We do not have the lyrics for") {
		fmt.Println("Error getting lyrics, retrying...")
		getTrack()
		getLyrics()
		return
	}

	lyrics := strings.Split(lyricsStr, "\n")
	size := 10
	if len(lyrics) < size {
		currentSong.lyrics = lyrics
		return
	}
	start := rand.IntN(len(lyrics) - size)
	currentSong.lyrics = lyrics[start : start+size]
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

// get a random track
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
