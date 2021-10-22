package fsread

import (
	"io/ioutil"
	"log"
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

// getKV read file tree in KVPair fileName->Content
// basedir starting directory
func getKV(basedir string) []KVPair {
	var output []KVPair
	err := filepath.Walk(basedir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if err != nil {
				log.Fatal(err)
			}
			if !info.IsDir() {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
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
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return output
}
