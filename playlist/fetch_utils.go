package playlist

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strings"

    "github.com/zmb3/spotify/v2"
)

// caches contents of user playlists to allow for future checking of updates
func cacheUserPlaylists(client *spotify.Client, ctx context.Context, spp *spotify.SimplePlaylistPage, names []string) error {
    var matches []*SpotifyPlaylist
    var err error

    // check if names array is empty
    if names != nil {
        // find playlist matching specified names
        matches, err = findPlaylistsByName(client, ctx, spp, names)
    } else {
        // default to using all found playlists
        matches, err = ConvertPlaylistsToSpotifyPlaylists(client, ctx, spp)
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


// find matching playlists based on name
func findPlaylistsByName(client *spotify.Client, ctx context.Context, spp *spotify.SimplePlaylistPage, names []string) ([]*SpotifyPlaylist, error) {
    var matches []*SpotifyPlaylist
    for _, sp := range spp.Playlists {
        for _, name := range names {
            if sp.Name == name {
                m := &SpotifyPlaylist{
                    Name:        name,
                    URI:         sp.URI,
                    Total_Tracks: sp.Tracks.Total,
                }

                tracks, err := getPlaylistTracks(client, ctx, m)
                if err != nil {
                    fmt.Println("Error fetching tracks for playlist:", err)
                    return nil, err
                }
                m.Tracks = tracks
                matches = append(matches, m)
                break
            }
        }
    }
    return matches, nil
}


// convert all playlists in the page to SpotifyPlaylist type and fetch their tracks
func ConvertPlaylistsToSpotifyPlaylists(client *spotify.Client, ctx context.Context, spp *spotify.SimplePlaylistPage) ([]*SpotifyPlaylist, error) {
    var playlists []*SpotifyPlaylist
    for _, sp := range spp.Playlists {
        p := &SpotifyPlaylist{
            Name:        sp.Name,
            URI:         sp.URI,
            Total_Tracks: sp.Tracks.Total,
        }

        tracks, err := getPlaylistTracks(client, ctx, p)
        if err != nil {
            fmt.Println("Error fetching tracks for playlist:", err)
            return nil, err
        }
        p.Tracks = tracks

        playlists = append(playlists, p)
        // fmt.Println(p)
    }
    return playlists, nil
}


// get tracks from playlist
func getPlaylistTracks(client *spotify.Client, ctx context.Context, p *SpotifyPlaylist) (string, error) {
    playlistID := extractPlaylistID(p.URI)

    // Fetch the playlist items using the playlist ID
    tracks, err := client.GetPlaylistItems(ctx, spotify.ID(playlistID))
    if err != nil {
        fmt.Println("Error getting tracks for playlist:", err)
        return "", err
    }

    // Encode the tracks as JSON
    t, err := json.Marshal(tracks)
    if err != nil {
        fmt.Println("JSON encoding failed:", err)
        return "", err
    }

    return string(t), nil
}


// Extract playlist ID from URI
func extractPlaylistID(uri spotify.URI) string {
    parts := strings.Split(string(uri), ":")
    if len(parts) == 3 {
        return parts[2]
    }
    return ""
}


// write playlists to a JSON file
func WritePlaylistsToFile(playlists []*SpotifyPlaylist, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")

    err = encoder.Encode(playlists)
    if err != nil {
        return err
    }

    fmt.Println("Playlists cached successfully to", filename)
    return nil
}

