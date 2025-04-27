package main

import (
	"fmt"
	"log"
	"os"
	//"path/filepath"
)

const debug = true
const journalEnv = "JOURNAL_DIR"
const postsEnv = "POSTS_DIR"
const estMapEntries = 100

func main() {
	err := LoadEnvFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	journalDir := os.Getenv(journalEnv)
	postsDir := os.Getenv(postsEnv)

	postsMap := make(map[string]*CustomFileInfo, estMapEntries)
	err = RecursiveReadDir(postsDir, postsMap)
	if err != nil {
		log.Fatalf("Failed to read posts directory: %v", err)
	}
	if debug {
		fmt.Println("postsMap:")
		PrettyPrint(postsMap)
	}

	journalMap := make(map[string]*CustomFileInfo, estMapEntries)
	err = RecursiveReadDir(journalDir, journalMap)
	if err != nil {
		log.Fatalf("Failed to read journal directory: %v", err)
	}
	if debug {
		fmt.Println("journalMap:")
		PrettyPrint(journalMap)
	}

	for _, jfi := range journalMap {
		if pfi, exists := postsMap[jfi.StandardizedFileName]; exists {
			err = SyncFiles(jfi, pfi)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, jfi := range journalMap {
		if _, exists := postsMap[jfi.StandardizedFileName]; !exists {
			err = CreatePost(jfi, postsDir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
