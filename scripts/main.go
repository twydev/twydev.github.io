package main

import (
	"fmt"
	"log"
	"os"
	//"path/filepath"
)

const debug = false
const booksEnv = "BOOKS_DIR"
const articlesEnv = "ARTICLES_DIR"
const researchEnv = "RESEARCH_DIR"
const postsEnv = "POSTS_DIR"
const estMapEntries = 100

func main() {
	err := LoadEnvFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	// Journal Dir
	booksDir := os.Getenv(booksEnv)
	articlesDir := os.Getenv(articlesEnv)
	researchDir := os.Getenv(researchEnv)
	// Blog Dir
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
	journalDir := []string{booksDir, articlesDir, researchDir}
	for _, dir := range journalDir {
		err = RecursiveReadDir(dir, journalMap)
		if err != nil {
			log.Fatalf("Failed to read journal directory: %v", err)
		}
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
