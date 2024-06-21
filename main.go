package main

import (
    "context"
    "fmt"
    "log"
    // "os"
    "net/http"

    "spotify-playlist-manager/playlist"

    "github.com/zmb3/spotify/v2/auth"
    "github.com/zmb3/spotify/v2"
)

const redirectURI = "http://localhost:8080/callback"

var (
    // TODO: pass token/secret into here as opts
    auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI),
        spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
    ch = make(chan *spotify.Client)
    state = "abc123"
)

func main() {
    // initialize http server
    http.HandleFunc("/callback", completeAuth)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Got request for:", r.URL.String())
    })
    go func() {
        err := http.ListenAndServe(":8080", nil)
        if err != nil {
            log.Fatal(err)
        }
    }()

    url := auth.AuthURL(state)
    fmt.Println("Please log into to Spotify by visiting the following page in your browser:", url)

    // wait for auth to complete
    client := <-ch

    ctx := context.Background()

    // use the client to make calls that require authorization
    user, err := client.CurrentUser(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("You are logged in as:", user.ID)

    // get user playlists
    playlists_page, err := playlist.FetchUserPlaylistsList(client, ctx)
    if err != nil {
        fmt.Println("Error fetching json playlist:", err)
        return
    }
    // fmt.Println(playlists)
    playlists, err := playlist.ConvertPlaylistsToSpotifyPlaylists(client, ctx, playlists_page)
    if err != nil {
        fmt.Println("Error converting playlists to type SpotifyPlaylist:", err)
        return
    }

    // write playlist data to json
    playlist.WritePlaylistsToFile(playlists, "output/user_playlists.json")
}

// auth token generation and checking
func completeAuth(w http.ResponseWriter, r *http.Request) {

    tok, err := auth.Token(r.Context(), state, r)
    if err != nil {
        fmt.Println("Couldn't get token:", http.StatusForbidden)
        log.Fatal(err)
    }

    if st := r.FormValue("state"); st != state {
        http.NotFound(w, r)
        log.Fatalf("State mismatch: %s != %s\n", st, state)
    }

    // use token to authorize client
    client := spotify.New(auth.Client(r.Context(), tok))
    fmt.Println(w, "Login Completed!")
    ch <- client
}

