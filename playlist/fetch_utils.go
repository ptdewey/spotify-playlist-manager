package playlist

import (
    "context"
    "encoding/json"
    "fmt"

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

// TODO: figure out if owner name is necessary
type userPlaylist struct {
    OwnerURI spotify.URI
    OwnerName string
    PlaylistURI spotify.URI
    PlaylistName string
}

// getPlaylists retrieves playlists and their tracks using the Spotify client.
//
// Params:
//   - client: Spotify API client.
//   - ctx: Context for request management.
//   - spp: SimplePlaylistPage with playlists.
//
// Returns:
//   - map[spotify.URI]*SpotifyPlaylist: Map of playlist URIs to SpotifyPlaylist structs.
//   - error: Error, if any occurred.
func getPlaylists(client *spotify.Client, ctx context.Context, spp *spotify.SimplePlaylistPage) (map[spotify.URI]*SpotifyPlaylist, error) {
    playlists := make(map[spotify.URI]*SpotifyPlaylist)
    // TODO: should this be parallelized?
    for _, sp := range spp.Playlists {
        p := &SpotifyPlaylist{
            Name: sp.Name,
            URI: sp.URI,
            Total_Tracks: sp.Tracks.Total,
            Owner: sp.Owner,
        }

        // TODO: store struct of owner and playlist names/uris?
        // - possibly some way that would allow immediate lookup in findPlaylistsByName

        tracks, err := getPlaylistTracks(client, ctx, p)
        if err != nil {
            fmt.Println("Error fetching tracks for playlist:", err)
            return nil, err
        }
        p.Tracks = tracks

        playlists[p.URI] = p
    }

    return playlists, nil
}


// convert all playlists in the page to SpotifyPlaylist type and fetch their tracks
func convertPlaylistsToSpotifyPlaylists(client *spotify.Client, ctx context.Context, spp *spotify.SimplePlaylistPage) (map[spotify.URI]*SpotifyPlaylist, error) {
    playlists, err := getPlaylists(client, ctx, spp)
    if err != nil {
        return nil, err
    }

    return playlists, nil
}

// find matching playlists based on name
func findPlaylistsByName(client *spotify.Client, ctx context.Context, spp *spotify.SimplePlaylistPage, names []string) (map[spotify.URI]*SpotifyPlaylist, error) {
    matches := make(map[spotify.URI]*SpotifyPlaylist)
    playlists, err := getPlaylists(client, ctx, spp)
    if err != nil {
        return nil, err
    }

    // TODO: replace nested loops with direct lookup
    matches = playlists
    // FIX: direct lookup requires playlist uris, not names
    // - names are unique per user, but not unique across users

    return matches, nil
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


