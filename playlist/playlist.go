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
    Owner spotify.User `json:"owner"`
}

type SimplifiedUser struct {
    URI spotify.URI `json:"uri"`
    Name string `json:"display_name"`
}

// caches contents of user playlists to allow for future checking of updates
func cacheUserPlaylists(client *spotify.Client, ctx context.Context, spp *spotify.SimplePlaylistPage, names []string) error {
    matches := make(map[spotify.URI]*SpotifyPlaylist)
    var err error

    // check if names array is empty
    if names != nil {
        // find playlist matching specified names
        // matches, err = findPlaylistsByName(client, ctx, spp, names)
    } else {
        // default to using all found playlists
        matches, err = convertPlaylistsToSpotifyPlaylists(client, ctx, spp)
    }
    if len(matches) == 0 {
        fmt.Println("No matches found")
        return nil
    }

    // iterate through matched playlists, getting current state
    for _, m := range matches {
        t, err := getPlaylistTracks(client, ctx, m)
        if err != nil {
            fmt.Println("Error fetching tracks for playlist:", err)
            break
        }
        m.Tracks = t
    }

    // write the matched playlists to a JSON file
    err = WritePlaylistsToFile(matches, "playlists_cache.json")
    if err != nil {
        fmt.Println("Error writing playlists to file:", err)
        return err
    }

    return nil
}

// get list of user playlists (names and IDs), output as json object
// func FetchUserPlaylistsList(client *spotify.Client, ctx context.Context) (string, error) {
func fetchUserPlaylistsList(client *spotify.Client, ctx context.Context) (*spotify.SimplePlaylistPage, error) {
    // get playlists from spotify api
    playlists, err := client.CurrentUsersPlaylists(ctx)
    // playlists, err := client.GetPlaylistsForUser(ctx, user.ID)
    if err != nil {
        fmt.Println("Error getting user playlists:", err)
        log.Println("Error getting user playlists:", err)
        return nil, err
    }

    return playlists, nil
}

// TODO: function for checking (probably specified) playlists for changes
// - probably return a list of changed playlists and/or json structure of changes (probably tuple return both)
// - would need to cache playlist contents
//  - # of tracks can be obtained from tracks -> total: in json schema (total songs is imperfect metric actually)
// - needs helper function to cache playlists
func checkPlaylistUpdated(client *spotify.Client, ctx context.Context) (bool, string, error) {
    // TODO: recache playlists
    return false, "", nil
}
