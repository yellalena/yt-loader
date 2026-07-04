package downloader

import (
	"fmt"
	"strings"
)

func stderrSuffix(stderr string) string {
	stderr = strings.TrimSpace(stderr)
	if stderr == "" {
		return ""
	}

	return ": " + stderr
}

type downloadLogFilter struct {
	buffer strings.Builder
}

func (f *downloadLogFilter) Write(p []byte) (int, error) {
	for _, b := range p {
		if b == '\n' || b == '\r' {
			f.flush()
			continue
		}

		f.buffer.WriteByte(b)
	}

	return len(p), nil
}

func (f *downloadLogFilter) flush() {
	line := strings.TrimSpace(f.buffer.String())
	f.buffer.Reset()

	if strings.HasPrefix(line, "[download]") {
		fmt.Println(line)
	}
}
