package playlist

import (
    "context"
    "fmt"
    "log"

    "github.com/zmb3/spotify/v2"
)


// Custom playlist struct type for easier and more memory efficient storage
type SpotifyPlaylist struct {
    Name string `json:"name"`
    URI spotify.URI `json:"uri"`
    Total_Tracks uint `json:"total_tracks"`
    Tracks string `json:"tracks"`
}


// get list of user playlists (names and IDs), output as json object
// func FetchUserPlaylistsList(client *spotify.Client, ctx context.Context) (string, error) {
func FetchUserPlaylistsList(client *spotify.Client, ctx context.Context) (*spotify.SimplePlaylistPage, error) {
    // get playlists from spotify api
    playlists, err := client.CurrentUsersPlaylists(ctx)
    // playlists, err := client.GetPlaylistsForUser(ctx, user.ID)
    if err != nil {
        fmt.Println("Error getting user playlists:", err)
        log.Println("Error getting user playlists:", err)
        // return "", err
        return nil, err
    }

    // convert playlists to json object
    // playlistsJSON, err := json.Marshal(playlists)
    // if err != nil {
    //     fmt.Println("Error converting to JSON:", err)
    //     return "", err
    // }

    return playlists, nil

    // return string(playlistsJSON), nil
}


// TODO: function for checking (probably specified) playlists for changes
// - probably return a list of changed playlists and/or json structure of changes (probably tuple return both)
// - would need to cache playlist contents
//  - # of tracks can be obtained from tracks -> total: in json schema (total songs is imperfect metric actually)
// - needs helper function to cache playlists
func CheckPlaylistUpdated(client *spotify.Client, ctx context.Context) (bool, string, error) {
    // TODO: recache playlists
    return false, "", nil
}
