package main

//package games

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var (
	playlist      *spotify.FullPlaylist
	previousSongs []string
	playedSongs   []song
	currentSong   song
	playerGuesses map[string]string
)

type song struct {
	songName   string
	artistName string
	lyrics     []string
	id         spotify.ID
	previewURL string
}

func main() {
	getPlaylist()
	getTrack()
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
	currentSong = song{
		songName:   randomTrack.Name,
		artistName: randomTrack.Artists[0].Name,
		id:         randomTrack.ID,
		previewURL: randomTrack.PreviewURL,
	}
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
