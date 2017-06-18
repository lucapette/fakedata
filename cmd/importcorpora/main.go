package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// relative package data dir
const targetDir = "pkg/data"

// Content of a Go file
const fileTemplate = `package data

// %s is an array of %s
var %s = %#v
`

var tasks = []struct {
	URL, Key, Var, File string
}{
	{
		URL:  "https://raw.githubusercontent.com/dariusk/corpora/master/data/words/nouns.json",
		Key:  "nouns",
		Var:  "Nouns",
		File: "nouns.go",
	},
	{
		URL:  "https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/common.json",
		Key:  "animals",
		Var:  "Animals",
		File: "animals.go",
	},
	{
		URL:  "https://raw.githubusercontent.com/dariusk/corpora/master/data/animals/cats.json",
		Key:  "cats",
		Var:  "Cats",
		File: "cats.go",
	},
	{
		URL:  "https://raw.githubusercontent.com/dariusk/corpora/master/data/words/emoji/emoji.json",
		Key:  "emoji",
		Var:  "Emoji",
		File: "emoji.go",
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

		file := filepath.Join(targetDir, task.File)
		if err := os.MkdirAll(filepath.Dir(file), 0777); err != nil {
			log.Fatal(err)
		}

		content := fmt.Sprintf(fileTemplate, task.Var, task.Key, task.Var, data)

		// Write to Go file
		if err := ioutil.WriteFile(file, []byte(content), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
