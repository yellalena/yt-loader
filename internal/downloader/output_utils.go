package downloader

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultOutputDir = "yt_loader_output"
	metadataDir      = ".metadata"
)

func createOutputDir(outputDir string) (string, error) {
	var destination string
	if outputDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		destination = filepath.Join(wd, defaultOutputDir)
	} else {
		absOutputDir, err := filepath.Abs(outputDir)
		if err != nil {
			return "", fmt.Errorf("failed to resolve output directory %q: %w", outputDir, err)
		}
		destination = absOutputDir
	}

	if err := os.MkdirAll(destination, 0750); err != nil {
		return "", fmt.Errorf("failed to create output directory %q: %w", destination, err)
	}

	return destination, nil
}

func createMetadataDir(outputDir string) error {
	metadataPath := filepath.Join(outputDir, metadataDir)
	if err := os.MkdirAll(metadataPath, 0750); err != nil {
		return fmt.Errorf("failed to create metadata directory %q: %w", metadataPath, err)
	}
	return nil

}

func outputTemplate(outputDir string, mode Mode) string {
	switch mode {
	case Playlist:
		return filepath.Join(outputDir, "%(playlist_title|Playlist).180B_playlist", "%(playlist_index&{} - |)s%(title).180B [%(id)s].%(ext)s")
	default:
		return filepath.Join(outputDir, "%(title).180B_video", "%(title).180B [%(id)s].%(ext)s")
	}
}

func archivePath(outputDir string, mode Mode, url string, playlistItems string) string {
	sum := sha1.Sum([]byte(string(mode) + ":" + playlistItems + ":" + url))
	return filepath.Join(outputDir, metadataDir, "archive-"+hex.EncodeToString(sum[:8])+".txt")
}
