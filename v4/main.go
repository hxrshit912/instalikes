
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	likesKey = "likestorm:likes"
)

func likeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentLikes, err := rdb.Incr(ctx, likesKey).Result()

	if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
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

	currentLikes, err := rdb.Get(ctx, likesKey).Int64()

	if err == redis.Nil {
		currentLikes = 0
	} else if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"likes": currentLikes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	err := rdb.Ping(ctx).Err()

	if err != nil {
		fmt.Println("Redis connection failed:", err)
		return
	}

	http.HandleFunc("/like", likeHandler)
	http.HandleFunc("/likes", likesHandler)

	fmt.Println("LikeStorm V4 running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
