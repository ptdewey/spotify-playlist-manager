package playlist

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/zmb3/spotify/v2"
)

// write playlists to a JSON file
func WritePlaylistsToFile(playlists map[spotify.URI]*SpotifyPlaylist, filename string) error {
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

// TODO:
func ReadPlaylistFile() ([]string, error) {
    return nil, nil
}

// Extract playlist ID from URI
func extractPlaylistID(uri spotify.URI) string {
    parts := strings.Split(string(uri), ":")
    if len(parts) == 3 {
        return parts[2]
    }
    return ""
}


