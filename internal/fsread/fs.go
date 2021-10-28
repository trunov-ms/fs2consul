package fsread

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type KVPair struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type FS struct {
	basedir string
	KV      []KVPair
}

func NewFSRead(basedir string) *FS {
	return &FS{
		basedir: basedir,
		KV:      getKV(basedir),
	}
}

// isTxtFile checking content file type
func isTxtFile(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buffer := make([]byte, 512)
	n, err := f.Read(buffer)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	contentType := http.DetectContentType(buffer[:n])
	return strings.HasPrefix(contentType, "text/")
}

// getKV read file tree in KVPair fileName->Content
//
// basedir - starting directory
func getKV(basedir string) []KVPair {
	var output []KVPair
	err := filepath.Walk(basedir,
		func(path string, info os.FileInfo, err error) error {
			// exclude directories and files from .git folder (if exist)
			if !info.IsDir() && !strings.Contains(path, "/.git/") {
				// filtering by content type
				if isTxtFile(path) {
					content, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					if strings.HasSuffix(string(content), "\n") {
						content = []byte(strings.TrimSuffix(string(content), "\n"))
					}
					kvp := KVPair{
						Key:   strings.TrimPrefix(path, basedir),
						Value: string(content),
					}
					output = append(output, kvp)
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return output
}
