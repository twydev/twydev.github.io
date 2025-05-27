package main

import (
	"fmt"
	"log"
	"os"
	//"path/filepath"
)

const debug = true
const booksEnv = "BOOKS_DIR"
const articlesEnv = "ARTICLES_DIR"
const postsEnv = "POSTS_DIR"
const estMapEntries = 100

func main() {
	err := LoadEnvFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	booksDir := os.Getenv(booksEnv)
	articlesDir := os.Getenv(articlesEnv)
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

	err = RecursiveReadDir(booksDir, journalMap)
	if err != nil {
		log.Fatalf("Failed to read journal directory: %v", err)
	}

	err = RecursiveReadDir(articlesDir, journalMap)
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
		} else {
			err = CreatePost(jfi, postsDir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
