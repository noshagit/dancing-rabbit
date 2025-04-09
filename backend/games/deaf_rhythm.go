package games

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"strings"

	Lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var playlist *spotify.FullPlaylist

func main() { // (w http.ResponseWriter) as input to work with html
	getPlaylist()
	songName, artistName := getTrack(playlist)

	fmt.Printf("%s ; by %s\n\n", songName, artistName)
	lyrics := strings.Split(getLyrics(songName, artistName), "\n")

	size := 10
	start := rand.IntN(len(lyrics) - size)
	lyrics = lyrics[start : start+size]
	fmt.Println(lyrics)

	/* tmpl := template.Must(template.ParseFiles("../html/deaf_rhythm.html"))
	tmpl.Execute(w, lyrics)*/
}

func getLyrics(songName string, artistName string) string {
	l := Lyrics.New()
	lyrics, err := l.Search(artistName, songName)
	for err != nil { // if no lyrics found, search for another song
		lyrics, err = l.Search(getTrack(playlist))
	}

	return lyrics
}

func getPlaylist() {
	// Spotify API setup
	ctx := context.Background()
	token := &oauth2.Token{
		AccessToken: "BQBR27gaX294-EoNUCXFNGNEI8reGdJCFO8-ZjPh-mgAjhLgvregKr792pswdfyQxtz-466gCFBGgPqsQCAxbimPBMN8AFykFKBxhn9p1xS3Fw5TRYvh1Ao8CXmnEj0GBrnE2nhC-6I",
		TokenType:   "Bearer",
	}

	// Create a Spotify client using the token
	client := spotify.New(spotifyauth.New().Client(ctx, token))

	playlistID := spotify.ID("0MSCX9tZWQmitMQsfhvZIl")                        // playlist cannot be private or made by spotify
	fields := spotify.Fields("tracks(total,items(track(name,artists(name)))") // only get nb of tracks, track name and artist name

	// Get the playlist
	var err error
	playlist, err = client.GetPlaylist(ctx, playlistID, fields)
	if err != nil {
		log.Fatalf("Failed to get playlist: %v", err)
	}
}

func getTrack(playlist *spotify.FullPlaylist) (string, string) {
	// get a random track
	randomTrackIndex := rand.IntN(int(math.Min(float64(playlist.Tracks.Total), 100))) // Limit to track 100 because of Spotify API limit
	randomTrack := playlist.Tracks.Tracks[randomTrackIndex].Track

	return randomTrack.Name, randomTrack.Artists[0].Name
}
