package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

// declare fn type for comparators
type fn func(time.Time, time.Time) bool

// store CLI args
var verbose bool

func myFlagUsage() {
	fmt.Printf("Usage: %s [OPTIONS] [path to search]\n", path.Clean(os.Args[0]))
	flag.PrintDefaults()
	fmt.Println("Note: Search is not recursive. Default path is CWD")
}

// getArgs - read & validate flags to cli
func getArgs() (string, bool) {
	flag.Usage = myFlagUsage
	oldestArgPtr := flag.Bool("oldest", false, "Search for oldest")
	newestArgPtr := flag.Bool("newest", false, "Search for newest")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	var dirPath string

	if flag.NArg() == 1 {
		// dirPath specified
		dirPath = flag.Args()[0]
	}

	findOldest := true

	if *oldestArgPtr {
		findOldest = true
	} else if *newestArgPtr {
		findOldest = false
	} else if path.Clean(os.Args[0]) == "newest" {
		// if called as newest, search for newest
		findOldest = false
	}

	// validate dirPath - default to CWD if empty
	var err error
	if dirPath == "" {
		dirPath, err = os.Getwd()
		if err != nil {
			quit(err.Error())
		}
	} else {
		if !isDirectory(dirPath) {
			quit("Not a directory")
		}
	}

	return dirPath, findOldest
}

func main() {
	dirPath, findOldest := getArgs()

	var fileResult string
	var fileResultErr error

	if findOldest {
		fileResult, fileResultErr = getOldest(dirPath)
	} else {
		fileResult, fileResultErr = getNewest(dirPath)
	}

	if fileResultErr != nil {
		quit(fileResultErr.Error())
	} else {
		fmt.Println(fileResult)
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
			if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				// file is symlink, skip
				continue
			}
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

// quit - exit in error, print quitMsg if verbose
func quit(quitMsg string) {
	if verbose {
		log.Fatal(quitMsg)
	} else {
		os.Exit(1)
	}
}
