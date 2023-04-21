package main

import (
	"encoding/json"
	"fmt"
	"io"
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

var keyExtractor = func(key string) func(io.ReadCloser) []string {
	return func(body io.ReadCloser) []string {
		var jsonData map[string]json.RawMessage
		if err := json.NewDecoder(body).Decode(&jsonData); err != nil {
			log.Fatal(err)
		}

		var data []string
		if err := json.Unmarshal(jsonData[key], &data); err != nil {
			log.Fatal(err)
		}
		return data
	}
}

var tasks = []struct {
	URL       string
	Extractor func(io.ReadCloser) []string
	Var       string
}{
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/words/nouns.json",
		keyExtractor("nouns"),
		"Nouns",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/common.json",
		keyExtractor("animals"),
		"Animals",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/dinosaurs.json",
		keyExtractor("dinosaurs"),
		"Dinosaurs",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/dogs.json",
		keyExtractor("dogs"),
		"Dogs",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/cats.json",
		keyExtractor("cats"),
		"Cats",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/words/emoji/emoji.json",
		keyExtractor("emoji"),
		"Emoji",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/words/adjs.json",
		keyExtractor("adjs"),
		"Adjectives",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/words/harvard_sentences.json",
		keyExtractor("data"),
		"Sentences",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/corporations/industries.json",
		keyExtractor("industries"),
		"Industries",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/humans/occupations.json",
		keyExtractor("occupations"),
		"Occupations",
	},
	{
		"https://raw.githubusercontent.com/dariusk/corpora/master/data/geography/us_cities.json",
		func(body io.ReadCloser) []string {
			var jsonData struct {
				Cities []struct {
					City string `json:"city"`
				} `json:"cities"`
			}
			if err := json.NewDecoder(body).Decode(&jsonData); err != nil {
				log.Fatal(err)
			}

			data := make([]string, len(jsonData.Cities))
			for i, city := range jsonData.Cities {
				data[i] = city.City
			}

			return data
		},
		"Cities",
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

		file := filepath.Join(targetDir, strings.ToLower(task.Var)+".go")
		if err := os.MkdirAll(filepath.Dir(file), 0777); err != nil {
			log.Fatal(err)
		}

		data := task.Extractor(resp.Body)

		content := fmt.Sprintf(fileTemplate, task.Var, data)

		// Write to Go file
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
