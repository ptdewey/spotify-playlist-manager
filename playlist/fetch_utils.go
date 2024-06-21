package playlist

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strings"

    "github.com/zmb3/spotify/v2"
)

type SimplifiedTrack struct {
    ID spotify.ID `json:"id"`
    URI spotify.URI `json:"uri"`
    Name string `json:"name"`
    Artists []SimplifiedArtist `json:"artists"`
    Album SimplifiedAlbum `json:"album"`
    TrackNumber int `json:"track_number"`
    Duration int `json:"duration_ms"`
    Popularity int `json:"popularity"`
    AddedBy spotify.User `json:"added_by"`
    AddedAt string `json:"added_at"`
    Explicit bool `json:"explicit"`
}

type SimplifiedArtist struct {
    ID spotify.ID `json:"id"`
    URI spotify.URI `json:"uri"`
    Name string `json:"name"`
}

type SimplifiedAlbum struct {
    ID spotify.ID `json:"id"`
    URI spotify.URI `json:"uri"`
    Name string `json:"name"`
    ReleaseDate string `json:"release_date"`
    AlbumGroup string `json:"album_group"`
}

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
                    Name: name,
                    URI: sp.URI,
                    Total_Tracks: sp.Tracks.Total,
                }

                tracks, err := getPlaylistTracks(client, ctx, m)
                if err != nil {
                    fmt.Println("Error fetching tracks for playlist:", err)
                    return nil, err
                }
                m.Tracks = tracks

                matches = append(matches, m)
                break // NOTE: this could possibly cause issues with duplicate name playlists (maybe not possible?)
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
            Name: sp.Name,
            URI: sp.URI,
            Total_Tracks: sp.Tracks.Total,
        }

        tracks, err := getPlaylistTracks(client, ctx, p)
        if err != nil {
            fmt.Println("Error fetching tracks for playlist:", err)
            return nil, err
        }
        p.Tracks = tracks

        playlists = append(playlists, p)
    }
    return playlists, nil
}


// get tracks from playlist
func getPlaylistTracks(client *spotify.Client, ctx context.Context, p *SpotifyPlaylist) (string, error) {
    var alltracks []SimplifiedTrack
    playlistID := extractPlaylistID(p.URI)

    limit := 100
    offset := 0

    // iterate through playlist pages (required for playlists with length > limit)
    for {
        // fetch the playlist items using the playlist ID
        tracks, err := client.GetPlaylistItems(ctx, spotify.ID(playlistID), spotify.Limit(limit), spotify.Offset(offset))
        if err != nil {
            fmt.Println("Error getting tracks for playlist:", err)
            return "", err
        }

        // remove unwanted fields from track (i.e. locales)
        for _, item := range tracks.Items {
            if item.Track.Track == nil {
                continue
            }
            it := item.Track.Track

            // extract simplified artist info
            var artists []SimplifiedArtist
            for _, artist := range it.Artists {
                a := SimplifiedArtist {
                    ID: artist.ID,
                    URI: artist.URI,
                    Name: artist.Name,
                }
                artists = append(artists, a)
            }

            // extract simplified album info
            album := SimplifiedAlbum {
                ID: it.Album.ID,
                URI: it.Album.URI,
                Name: it.Album.Name,
                ReleaseDate: it.Album.ReleaseDate,
                AlbumGroup: it.Album.AlbumGroup,
            }

            // store simplified track information
            simplifiedtrack := SimplifiedTrack{
                ID: it.ID,
                URI: it.URI,
                Name: it.Name,
                Artists: artists,
                Album: album,
                TrackNumber: it.TrackNumber,
                Duration: it.Duration,
                Popularity: it.Popularity,
                Explicit: it.Explicit,
                AddedBy: item.AddedBy,
                AddedAt: item.AddedAt,
            }

            alltracks = append(alltracks, simplifiedtrack)
        }

        if len(tracks.Items) < limit {
            break;
        }

        offset += limit
    }
    fmt.Printf("Number of tracks in %s: %d\n", playlistID, len(alltracks))

    // encode the tracks as JSON
    t, err := json.Marshal(alltracks)
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

