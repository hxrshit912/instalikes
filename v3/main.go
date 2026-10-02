package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	db *sql.DB
	mu sync.Mutex
)

func likeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()

	_, err := db.Exec(`
		UPDATE likes
		SET count = count + 1
		WHERE id = 1
	`)

	if err != nil {
		mu.Unlock()
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var currentLikes int

	err = db.QueryRow(`
		SELECT count
		FROM likes
		WHERE id = 1
	`).Scan(&currentLikes)

	mu.Unlock()

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

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

	var currentLikes int

	err := db.QueryRow(`
		SELECT count
		FROM likes
		WHERE id = 1
	`).Scan(&currentLikes)

	mu.Unlock()

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"likes": currentLikes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func initDatabase() error {
	var err error

	db, err = sql.Open("sqlite", "likestorm.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY,
			count INTEGER NOT NULL
		)
	`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT OR IGNORE INTO likes (id, count)
		VALUES (1, 0)
	`)

	return err
}

func main() {
	err := initDatabase()

	if err != nil {
		fmt.Println("Database initialization failed:", err)
		return
	}

	defer db.Close()

	http.HandleFunc("/like", likeHandler)
	http.HandleFunc("/likes", likesHandler)

	fmt.Println("LikeStorm V3 running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}

