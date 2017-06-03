// Usage: go run cmd/importcorpora/main.go
//
// Updates the at the bottom of this file specified data files with content from dariusk/corpora.
package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// URL to zip file
const corporaArchive = "https://github.com/dariusk/corpora/archive/master.zip"

// Path to write Go file to
const targetDir = "pkg/data"

// Content of a Go file
const fileTemplate = `package data

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

	dir, err := downloadCorpora()
	defer func() {
		if dir == "" {
			return
		}
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range data {
		// Get data from JSON file
		f, err := ioutil.ReadFile(filepath.Join(dir, "corpora-master/data", d.From))
		if err != nil {
			log.Fatal(err)
		}
		var jsonData map[string]interface{}
		if err := json.Unmarshal(f, &jsonData); err != nil {
			log.Fatal(err)
		}

		// Fix formatting
		value := strings.Replace(fmt.Sprintf("%#v\n", jsonData[d.Key]), "interface {}", "string", 1)
		content := fmt.Sprintf(fileTemplate, d.Var, value)

		// Create required directories
		to := filepath.Join("pkg/data", d.To)
		if err := os.MkdirAll(filepath.Dir(to), 0777); err != nil {
			log.Fatal(err)
		}

		// Write to Go file
		if err := ioutil.WriteFile(to, []byte(content), 0644); err != nil {
			log.Fatal(err)
		}
	}
}

// Download corpora repo into a tmp dir and return path to dir
func downloadCorpora() (string, error) {
	// Create tmp dir
	dir, err := ioutil.TempDir("", "fakedata-import-")
	if err != nil {
		return dir, err
	}

	// Create tmp file to write zip file to
	tmpArchive, err := ioutil.TempFile("", "fakedata-import-archive")
	if err != nil {
		return dir, err
	}
	defer func() {
		if err := os.Remove(tmpArchive.Name()); err != nil {
			panic(err)
		}
	}()

	// Download zip file
	resp, err := http.Get(corporaArchive)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	_, err = io.Copy(tmpArchive, resp.Body)
	if err != nil {
		return dir, err
	}
	if err := tmpArchive.Close(); err != nil {
		return dir, err
	}

	// Unzip to tmp dir (adopted from https://stackoverflow.com/a/24792688/986455)
	r, err := zip.OpenReader(tmpArchive.Name())
	if err != nil {
		return dir, err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	extractAndWriteFile := func(zipFile *zip.File) error {
		zipReader, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := zipReader.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dir, zipFile.Name)
		if zipFile.FileInfo().IsDir() {
			return os.MkdirAll(path, zipFile.Mode())
		}
		if err := os.MkdirAll(filepath.Dir(path), zipFile.Mode()); err != nil {
			return err
		}

		target, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			if err := target.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(target, zipReader)
		return err
	}

	for _, f := range r.File {
		if err := extractAndWriteFile(f); err != nil {
			return dir, err
		}
	}

	return dir, nil
}

// Data to import is specified
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
