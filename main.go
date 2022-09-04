package main

import (
	"errors"
	"log"
	"os"
	"time"
)

func main() {
	// get path or default to cwd
	var path string
	var err error
	if len(os.Args) < 2 {
		path, err = os.Getwd()
		if err != nil {
			_ = path // get around path not used
			log.Fatal(err)
		}
	} else {
		path = os.Args[1]
		if !isDirectory(path) {
			log.Fatal("Not a directory")
		}
	}

	oldestFile, getOldErr := getOldest(path)

	if getOldErr != nil {
		log.Fatal(getOldErr)
	} else {
		print(oldestFile)
	}
}

// isDirectory determines if a file represented
// by `path` is a directory or not
func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func getOldest(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	var oldestFile string
	var oldestFileTime time.Time
	unset := true

	// go through files & get oldest
	for _, file := range files {
		if !file.IsDir() {
			fileInfo, _ := file.Info()
			if unset || fileInfo.ModTime().Before(oldestFileTime) {
				oldestFile = file.Name()
				oldestFileTime = fileInfo.ModTime()
				unset = false
			}
		}
	}

	// no files found
	if unset {
		return "", errors.New("no files found")
	} else {
		return oldestFile, nil
	}
}
