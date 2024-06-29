# Spotify Playlist Manager

A Go tool enabling better management of Spotify playlists.

## Motivation

With UI and interactivity changes that Spotify has made to their client apps over the last couple of years, I felt the need to implement a better way of managing my playlists by using their API.
In particular, the change from the 'like song button' to the 'add to playlist button' severely disrupted my way of using the app, and I wanted to make a system that would allow me to regain that functionality plus some new features.

## Features

- **Fetch and Display Playlists**: Retrieve and display your Spotify playlists along with their metadata.
- **Cache Playlists**: Cache playlist data locally for quick access and comparison with current states.
- **Track Management**: Fetch and display tracks within playlists, allowing for detailed inspection and updates.
- **Future Enhancements (TODO)**: Integration with HTMX (or some JS stuff) to allow non-cli usage.

## Installation

1. **Clone the Repository**
    ```sh
    git clone https://github.com/yourusername/spotify-playlist-manager.git
    cd spotify-playlist-manager
    ```

2. **Set Up Environment Variables**
    Set your Spotify API credentials as environment variables:
    ```sh
    export SPOTIFY_ID="your_spotify_client_id"
    export SPOTIFY_SECRET="your_spotify_client_secret"
    ```

3. **Install Dependencies**
    Make sure you have Go installed. Then, install the necessary Go packages:
    ```sh
    go mod tidy
    ```

4. **Run the Application**
    ```sh
    go run main.go
    ```

## Usage

### Fetching and Caching Playlists

1. **Retrieve Playlists**
    The application will fetch your Spotify playlists and display them in the terminal. You can customize which playlists to fetch by modifying the `names` slice in the `cacheUserPlaylists` function.

2. **Cache Playlists**
    Playlist data, including tracks, will be cached to a JSON file (`playlists_cache.json`) for easy access and future comparison.

### Example Code

See [main.go](main.go) for example usage (broader functionality is WIP).

### Future Enhancements

- **Web Interface with HTMX**: Integrate HTMX to create a dynamic and interactive web interface for managing playlists.
- **Playlist Comparison**: Add functionality to compare cached playlists with the current state to identify changes.
- **Advanced Filtering**: Enable advanced filtering and sorting options for playlists and tracks.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an Issue to discuss improvements or new features.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
