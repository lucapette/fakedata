package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// relative package data dir
const targetDir = "pkg/data"

// Content of a Go file
const fileTemplate = `package data

var %s = %#v
`

var tasks = []struct {
	URL, Key, Var string
}{
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/words/nouns.json",
		"nouns",
		"Nouns",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/common.json",
		"animals",
		"Animals",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/dogs.json",
		"dogs",
		"Dogs",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/cats.json",
		"cats",
		"Cats",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/words/emoji/emoji.json",
		"emoji",
		"Emoji",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/words/adjs.json",
		"adjs",
		"Adjectives",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/words/harvard_sentences.json",
		"data",
		"Sentences",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/corporations/industries.json",
		"industries",
		"Industries",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/humans/occupations.json",
		"occupations",
		"Occupations",
	},
}

func main() {
	// Check if running in repository directory
	_, err := os.Stat(targetDir)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatalf("the data directory cannot be found at %s. Ensure the importer is running in the correct location: %v", targetDir, err)
	}

	for _, task := range tasks {
		// Get JSON from URL
		resp, err := http.Get(task.URL)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		var jsonData map[string]json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
			log.Fatal(err)
		}

		var data []string
		if err := json.Unmarshal(jsonData[task.Key], &data); err != nil {
			log.Fatal(err)
		}

		file := filepath.Join(targetDir, strings.ToLower(task.Var)+".go")
		if err := os.MkdirAll(filepath.Dir(file), 0777); err != nil {
			log.Fatal(err)
		}

		content := fmt.Sprintf(fileTemplate, task.Var, data)

		// Write to Go file
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
