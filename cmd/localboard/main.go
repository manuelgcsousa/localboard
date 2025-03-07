package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

const (
	ServerPort    string = ":27049"
	ClipboardFile string = "./data/clipboard.txt"
)

var (
	mu        sync.Mutex
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan string)
)

func main() {
	// Handle all public files with FileServer
	fs := http.FileServer(http.Dir("./public/"))
	http.Handle("/", fs)

	http.Handle("/ws", websocket.Handler(handleWebSocket))

	// Start goroutine to handle message broadcasting to all connected clients
	go handleBroadcast()

	// Wrap default serve mux with logger middleware
	handler := logger(http.DefaultServeMux)

	// Start server
	fmt.Printf("Server listening on port %s...\n", ServerPort)
	log.Fatal(http.ListenAndServe(ServerPort, handler))
}

func handleWebSocket(ws *websocket.Conn) {
	defer ws.Close()

	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	if content, err := os.ReadFile(ClipboardFile); err == nil {
		if _, err = ws.Write(content); err != nil {
			log.Println("Error sending initial clipboard data: ", err)
			return
		}
	}

	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			removeClient(ws)
			break
		}

		// Save new clipboard data
		if err := saveClipboardData(message); err != nil {
			log.Println("Error saving clipboard data: ", err)
			continue
		}

		// Broadcast message to other clients
		broadcast <- message
	}
}

// handleBroadcast receives messages from the 'broadcast' channel
// and sends it to all connected clients.
func handleBroadcast() {
	for {
		message := <-broadcast

		mu.Lock()
		for client := range clients {
			websocket.Message.Send(client, message)
		}
		mu.Unlock()
	}
}

// saveClipboardContent writes the clipboard content to the file.
func saveClipboardData(content string) error {
	return os.WriteFile(ClipboardFile, []byte(content), 0644)
}

// removeClient removes a client from the clients map and closes the connection.
func removeClient(ws *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, ws)
	ws.Close()
}

// logger is a simple middleware that logs HTTP requests.
// Logs the current timestamp, HTTP method, request path and client IP address.
func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// clear default log flags
		log.SetFlags(0)

		log.Printf(
			"[%s] %s %s %s",
			time.Now().Format("2006/01/02 15:04:05"),
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		next.ServeHTTP(w, r)
	})
}
