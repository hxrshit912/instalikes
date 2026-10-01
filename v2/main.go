package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	likes int
	mu    sync.Mutex
)

func likeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	likes++
	currentLikes := likes
	mu.Unlock()

	response := map[string]interface{}{
		"message": "liked",
		"likes":   currentLikes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func likesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	currentLikes := likes
	mu.Unlock()

	response := map[string]interface{}{
		"likes": currentLikes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/like", likeHandler)
	http.HandleFunc("/likes", likesHandler)

	fmt.Println("LikeStorm V2 running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
