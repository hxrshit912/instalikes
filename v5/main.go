package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	shardCounter uint64

	numberOfShards = 10
)

func getNextShard() uint64 {
	return atomic.AddUint64(&shardCounter, 1) % uint64(numberOfShards)
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shard := getNextShard()

	key := fmt.Sprintf("likestorm:likes:%d", shard)

	currentLikes, err := rdb.Incr(ctx, key).Result()

	if err != nil {
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "liked",
		"shard":   shard,
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

	var totalLikes int64

	for shard := 0; shard < numberOfShards; shard++ {

		key := fmt.Sprintf("likestorm:likes:%d", shard)

		currentLikes, err := rdb.Get(ctx, key).Int64()

		if err == redis.Nil {
			continue
		}

		if err != nil {
			http.Error(w, "Redis error", http.StatusInternalServerError)
			return
		}

		totalLikes += currentLikes
	}

	response := map[string]interface{}{
		"likes": totalLikes,
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

	fmt.Println("LikeStorm V5 running on http://localhost:8080")
	fmt.Println("Using", numberOfShards, "Redis shards")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
