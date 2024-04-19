package playlist

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/zmb3/spotify/v2"
)


// caches contents of user playlists to allow for future checking of updates
func cacheUserPlaylists(client *spotify.Client, ctx context.Context, spp *spotify.SimplePlaylistPage, names []string) error {

    // find playlist matching specified names
    matches := findPlaylistsByName(spp, names)
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
        // TODO: write to file with name, uri, length, and tracks
        // TODO: alternatively use sqlite database
        // TODO: check this is an ok assignment, check json can be read correctly (might need [])
        m.Tracks = t
    }

    return nil
}


// find matching playlists based on name
func findPlaylistsByName(spp *spotify.SimplePlaylistPage, names []string) []*SpotifyPlaylist {
    // collect names, URIs, # tracks for specified playlists
    var matches []*SpotifyPlaylist
    // TODO: this could be parallelized
    for _, sp := range spp.Playlists {
        for _, name := range names {
            if sp.Name == name {
                m := &SpotifyPlaylist {
                    Name: name,
                    URI: sp.URI,
                    Total_Tracks: sp.Tracks.Total,
                    Tracks: "",
                }
                matches = append(matches, m)
                fmt.Println(m)
                break
            }
        }
    }
    return matches
}


// get tracks from playlist
func getPlaylistTracks(client *spotify.Client, ctx context.Context, p *SpotifyPlaylist) (string, error) {
    // get tracks using URI
    tracks, err := client.GetPlaylistItems(ctx, spotify.ID(p.URI))
    if err != nil {
        fmt.Println("Error getting tracks for playlist:", err)
        return "", err
    }

    // encode as JSON
    t, err := json.Marshal(tracks)
    if err != nil {
        fmt.Println("JSON encoding failed:", err)
        return "", err
    }

    return string(t), nil
}

