package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "yt-loader",
	Short:         "yt-loader downloads YouTube videos and playlists as MP3 files",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `yt-loader is a CLI tool for downloading YouTube videos and playlists as MP3 files.

It uses yt-dlp for downloading and ffmpeg for MP3 conversion, so both tools must be
installed and available on PATH.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
