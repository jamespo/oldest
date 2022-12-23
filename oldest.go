package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// declare fn type for comparators
type fn func(time.Time, time.Time) bool

func main() {
	// get CLI args path or default to cwd
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
		fmt.Println(oldestFile)
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

// fileIsOlder checks if fileModTime passed is older than oldestFileTime
func fileIsOlder(fileModTime time.Time, oldestFileTime time.Time) bool {
	return fileModTime.Before(oldestFileTime)
}

// fileIsNewer checks if fileModTime passed is newer than newestFileTime
func fileIsNewer(fileModTime time.Time, newestFileTime time.Time) bool {
	return fileModTime.After(newestFileTime)
}

// getOldest returns oldest file in path
func getOldest(path string) (string, error) {
	return findFile(path, fileIsOlder)
}

// getNewest returns newest file in path
func getNewest(path string) (string, error) {
	return findFile(path, fileIsNewer)
}

// findFile - gets the oldest or newest file from supplied path and comparator
func findFile(path string, comparator fn) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	var oldestOrNewestFile string
	var oldestOrNewestFileTime time.Time
	unset := true

	// go through files & get oldest
	for _, file := range files {
		if !file.IsDir() {
			fileInfo, _ := file.Info()
			fileModTime := fileInfo.ModTime()
			if unset || comparator(fileModTime, oldestOrNewestFileTime) {
				oldestOrNewestFile = file.Name()
				oldestOrNewestFileTime = fileModTime
				unset = false
			}
		}
	}

	// no files found
	if unset {
		return "", errors.New("no files found")
	} else {
		return oldestOrNewestFile, nil
	}
}
