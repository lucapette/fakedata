// Usage: go run cmd/importcorpora/main.go
//
// Updates the at the bottom of this file specified data from dariusk/corpora.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Base URL to download JSON files from
const baseURL = "https://raw.githubusercontent.com/dariusk/corpora/master/data/"

// Path to write Go files to
const targetDir = "pkg/data"

// Content of a Go file
const fileTemplate = `package data

// %s is an array of %s
var %s = %s`

func main() {
	// Check if running in repo directory
	_, err := os.Stat(targetDir)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatalf("The data directory cannot be found at %s. Ensure the importer is running in the correct location.", targetDir)
	}

	for _, d := range data {
		// Get JSON from URL
		resp, err := http.Get(baseURL + d.From)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		// Get data from JSON
		var jsonData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&jsonData); err != nil {
			log.Fatal(err)
		}

		// Fix formatting
		value := strings.Replace(fmt.Sprintf("%#v\n", jsonData[d.Key]), "interface {}", "string", 1)
		content := fmt.Sprintf(fileTemplate, d.Var, strings.ToLower(d.Var), d.Var, value)

		// Create required directories
		to := filepath.Join(targetDir, d.To)
		if err := os.MkdirAll(filepath.Dir(to), 0777); err != nil {
			log.Fatal(err)
		}

		// Write to Go file
		if err := ioutil.WriteFile(to, []byte(content), 0644); err != nil {
			log.Fatal(err)
		}
	}
}

// Data to import is specified here
var data = []struct {
	// JSON file to read from
	From string
	// Key containing data in JSON file
	Key string
	// Go File to write to
	To string
	// Variable name in the Go File
	Var string
}{
	{
		From: "words/nouns.json",
		Key:  "nouns",
		To:   "nouns.go",
		Var:  "Nouns",
	},
	{
		From: "animals/common.json",
		Key:  "animals",
		To:   "animals.go",
		Var:  "Animals",
	},
	{
		From: "animals/cats.json",
		Key:  "cats",
		To:   "cats.go",
		Var:  "Cats",
	},
	{
		From: "words/emoji/emoji.json",
		Key:  "emoji",
		To:   "emojis.go",
		Var:  "Emojis",
	},
}
