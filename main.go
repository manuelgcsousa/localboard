package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type ClipboardData struct {
	Data string `json:"data"`
}

const (
	ServerPort    string = ":27049"
	ClipboardFile string = "./data/clipboard.txt"
)

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/script.js", serveScript)
	http.HandleFunc("/read", handleRead)
	http.HandleFunc("/save", handleSave)

	handler := Logger(http.DefaultServeMux)

	fmt.Printf("Server listening on port %s...\n", ServerPort)
	log.Fatal(http.ListenAndServe(ServerPort, handler))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func serveScript(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/script.js")
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// check if clipboard file exists
	// if not, create an empty file
	if _, err := os.Stat(ClipboardFile); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(ClipboardFile, []byte{}, 0644)
		if err != nil {
			http.Error(w, "Error while creating file", http.StatusInternalServerError)
			return
		}
	}

	// read clipboard text file
	content, err := os.ReadFile(ClipboardFile)
	if err != nil {
		http.Error(w, "Error while reading file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	var clipboardData ClipboardData
	err := json.NewDecoder(r.Body).Decode(&clipboardData)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	err = os.WriteFile(ClipboardFile, []byte(clipboardData.Data), 0644)
	if err != nil {
		http.Error(w, "Error while writing to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Clipboard saved"))
}
