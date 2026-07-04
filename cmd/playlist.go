package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yellalena/yt-loader/internal/downloader"
)

var playlistOptions downloader.Options

var playlistCmd = &cobra.Command{
	Use:   "playlist URL [URL ...]",
	Short: "Download one or more YouTube playlists as MP3 files",
	Long: `Download one or more YouTube playlists as MP3 files.

Each argument must be a YouTube playlist URL or a channel tab URL that yt-dlp can
read as a playlist.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return downloader.Download(context.Background(), downloader.Playlist, args, playlistOptions)
	},
}

func init() {
	playlistCmd.Flags().
		StringVarP(&playlistOptions.OutputDir, "output", "o", "", "Output directory. Defaults to ./yt_loader_output under the location where the command is run.")
	playlistCmd.Flags().
		BoolVarP(&playlistOptions.Verbose, "verbose", "v", false, "Show raw yt-dlp output.")

	if err := playlistCmd.MarkFlagDirname("output"); err != nil {
		panic(fmt.Sprintf("mark output flag as directory: %v", err))
	}

	rootCmd.AddCommand(playlistCmd)
}
