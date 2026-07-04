package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yellalena/yt-loader/internal/downloader"
)

var videoOptions downloader.Options

var videoCmd = &cobra.Command{
	Use:   "video URL [URL ...]",
	Short: "Download one or more YouTube videos as MP3 files",
	Long: `Download one or more YouTube videos as MP3 files.

Each argument must be a YouTube video URL. If a URL also belongs to a playlist,
only that individual video is downloaded.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return downloader.Download(context.Background(), downloader.Video, args, videoOptions)
	},
}

func init() {
	videoCmd.Flags().
		StringVarP(&videoOptions.OutputDir, "output", "o", "", "Output directory. Defaults to ./yt_loader_output under the location where the command is run.")
	videoCmd.Flags().
		BoolVarP(&videoOptions.Verbose, "verbose", "v", false, "Show raw yt-dlp output.")

	if err := videoCmd.MarkFlagDirname("output"); err != nil {
		panic(fmt.Sprintf("mark output flag as directory: %v", err))
	}

	rootCmd.AddCommand(videoCmd)
}
