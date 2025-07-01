package main

import (
	"fmt"
	"log"
	"os"
	//"path/filepath"
)

const debug = false
const estMapEntries = 100

const envJournalBooks = "JOURNAL_DIR_BOOKS"
const envJournalArticles = "JOURNAL_DIR_ARTICLES"
const envJournalResearch = "JOURNAL_DIR_RESEARCH"

const envBlogNotes = "DIR_NOTES"
const envBlogResearch = "DIR_RESEARCH"

func syncDir(journalDirs []string, blogDir string) {
	fmt.Println("Syncing", journalDirs, "to", blogDir)

	blogMap := make(map[string]*CustomFileInfo, estMapEntries)
	err := RecursiveReadDir(blogDir, blogMap)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}
	if debug {
		fmt.Println("blogMap:")
		PrettyPrint(blogMap)
	}

	journalMap := make(map[string]*CustomFileInfo, estMapEntries)
	for _, dir := range journalDirs {
		err = RecursiveReadDir(dir, journalMap)
		if err != nil {
			log.Fatalf("Failed to read directory: %v", err)
		}
	}
	if debug {
		fmt.Println("journalMap:")
		PrettyPrint(journalMap)
	}

	for _, jfi := range journalMap {
		if bfi, exists := blogMap[jfi.StandardizedFileName]; exists {
			err = SyncFiles(jfi, bfi)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = CreatePost(jfi, blogDir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	err := LoadEnvFile(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Journal Dir
	jdBooks := os.Getenv(envJournalBooks)
	jdArticles := os.Getenv(envJournalArticles)
	jdResearch := os.Getenv(envJournalResearch)

	// Blog Dir
	dNotes := os.Getenv(envBlogNotes)
	dResearch := os.Getenv(envBlogResearch)

	syncDir([]string{jdBooks, jdArticles}, dNotes)
	syncDir([]string{jdResearch}, dResearch)
}
