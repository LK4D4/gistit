package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const apiUrl = "https://api.github.com/gists"

type File struct {
	Content string `json:"content"`
}

type Data struct {
	Description string          `json:"description"`
	Files       map[string]File `json:"files"`
}

type Answer struct {
	URL string `json:"html_url"`
}

func main() {
	data := Data{
		Files: make(map[string]File),
	}
	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		data.Description = os.Args[1]
	}
	if len(os.Args) > 2 {
		data.Files[os.Args[2]] = File{Content: string(content)}
	} else {
		data.Files[filepath.Base(os.Args[0])] = File{Content: string(content)}
	}
	raw, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(apiUrl, "application/json", bytes.NewReader(raw))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Wrong status from GitHub %d, should be %d, body: %q", resp.StatusCode, http.StatusCreated, body)
	}
	ans := Answer{}
	if err := json.NewDecoder(resp.Body).Decode(&ans); err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Unable to unmarshal result: %q", body)
	}
	println(ans.URL)
}
