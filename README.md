# yt-loader

A tiny tool to download youtube videos in mp3 format.

Only download videos you own or have permission to save.

## Requirements

Install the external tools used by the CLI:

```sh
brew install yt-dlp ffmpeg
```

## Build

```sh
mkdir -p bin
go build -o bin/yt-loader .
```

## Usage

Download one or more individual videos:

```sh
./bin/yt-loader video "https://www.youtube.com/watch?v=VIDEO_ID"
```

Download one or more playlists:

```sh
./bin/yt-loader playlist "https://www.youtube.com/playlist?list=PLAYLIST_ID"
```

By default, the CLI shows concise progress messages and `yt-dlp` `[download]` progress lines while hiding other raw `yt-dlp` logs.
Use `--verbose` to stream raw `yt-dlp` output:

```sh
./bin/yt-loader playlist --verbose "https://www.youtube.com/playlist?list=PLAYLIST_ID"
```

Specify output location:

```sh
./bin/yt-loader playlist \
  --output ~/Downloads/my-mp3s \
  "https://www.youtube.com/playlist?list=PLAYLIST_ID"
```

If `--output` is omitted, files are written under `./yt_loader_output` in the directory where the command was run.

Each URL gets its own media folder:

- videos: `<video title>_video`
- playlists: `<playlist title>_playlist`

Playlist downloads are parallelized automatically by sharding each playlist into four `yt-dlp` workers, so a single large playlist can download multiple entries at the same time.
