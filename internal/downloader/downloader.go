package downloader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type Mode string

const (
	Video               Mode = "video"
	Playlist            Mode = "playlist"
	playlistParallelism      = 4
)

type Options struct {
	OutputDir string
	Verbose   bool
}

func Download(ctx context.Context, mode Mode, urls []string, options Options) error {
	if len(urls) == 0 {
		return errors.New("provide at least one URL")
	}

	outputDir, err := createOutputDir(options.OutputDir)
	if err != nil {
		return err
	}

	if err := requireCommand("yt-dlp"); err != nil {
		return err
	}
	if err := requireCommand("ffmpeg"); err != nil {
		return err
	}

	fmt.Printf("Output: %s\n", outputDir)

	if mode == Playlist {
		if err := createMetadataDir(outputDir); err != nil {
			return err
		}
		return downloadPlaylists(ctx, outputDir, urls, options.Verbose)
	}

	return downloadVideos(ctx, outputDir, urls, options.Verbose)
}

func requireCommand(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s is required but was not found on PATH", name)
	}

	return nil
}

func downloadVideos(ctx context.Context, outputDir string, urls []string, verbose bool) error {
	for index, url := range urls {
		fmt.Printf("Downloading video %d/%d: %s\n", index+1, len(urls), url)

		if err := runYTDLP(ctx, Video, outputDir, url, "", verbose); err != nil {
			return err
		}

		fmt.Printf("Downloaded video %d/%d\n", index+1, len(urls))
	}

	return nil
}

func downloadPlaylists(ctx context.Context, outputDir string, urls []string, verbose bool) error {
	for index, url := range urls {
		fmt.Printf("Downloading playlist %d/%d: %s\n", index+1, len(urls), url)

		if err := downloadPlaylistParallel(ctx, outputDir, url, verbose); err != nil {
			return err
		}

		fmt.Printf("Downloaded playlist %d/%d\n", index+1, len(urls))
	}

	return nil
}

func downloadPlaylistParallel(ctx context.Context, outputDir, url string, verbose bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, 1)

	var wg sync.WaitGroup
	for shard := 1; shard <= playlistParallelism; shard++ {
		wg.Add(1)
		go func(shard int) {
			defer wg.Done()

			playlistItems := fmt.Sprintf("%d::%d", shard, playlistParallelism)
			if err := runYTDLP(ctx, Playlist, outputDir, url, playlistItems, verbose); err != nil {
				select {
				case errs <- err:
					cancel()
				default:
				}
			}
		}(shard)
	}

	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}

func runYTDLP(ctx context.Context, mode Mode, outputDir, url, playlistItems string, verbose bool) error {
	args, err := buildArgs(mode, outputDir, url, playlistItems, verbose)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	var stderr bytes.Buffer
	if verbose {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		filter := downloadLogFilter{}
		cmd.Stdout = &filter
		cmd.Stderr = &stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp failed for %q: %w%s", url, err, stderrSuffix(stderr.String()))
	}

	return nil
}

func buildArgs(mode Mode, outputDir, url, playlistItems string, verbose bool) ([]string, error) {
	args := []string{
		"--ignore-config",
		"--continue",
		"--ignore-errors",
		"--no-overwrites",
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--embed-metadata",
		"--embed-thumbnail",
		"--convert-thumbnails", "jpg",
		"--restrict-filenames",
		"--windows-filenames",
		"--output", outputTemplate(outputDir, mode),
	}

	if !verbose {
		args = append([]string{"--progress", "--no-warnings"}, args...)
	}

	switch mode {
	case Video:
		args = append(args, "--no-playlist")
	case Playlist:
		args = append(
			args,
			"--yes-playlist",
			"--download-archive", archivePath(outputDir, mode, url, playlistItems),
			"--playlist-items", playlistItems,
		)
	default:
		return nil, fmt.Errorf("unknown download mode %q", mode)
	}

	return append(args, url), nil
}
