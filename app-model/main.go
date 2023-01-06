package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/bearatol/lg"
	"github.com/go-redis/redis/v9"
)

func main() {
	redisPass, ok := os.LookupEnv("REDIS_PASS")
	if !ok {
		lg.Fatal("REDIS_PASS environment variable is not set")
	}
	redisPort, ok := os.LookupEnv("REDIS_PORT")
	if !ok {
		lg.Fatal("REDIS_PORT environment variable is not set")
	}
	modelPort, ok := os.LookupEnv("MODEL_PORT")
	if !ok {
		lg.Fatal("MODEL_PORT environment variable is not set")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:" + redisPort,
		Password: redisPass,
		DB:       0,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/get":
			num, err := getRedis(rdb)
			if err != nil {
				http.Error(w, "some problem, see in logs", http.StatusBadRequest)
				lg.Error(err)
				return
			}
			fmt.Fprintf(w, "{\"number\": %d}", num)
		case "/set":
			if req.Method != "POST" {
				http.Error(w, "", http.StatusMethodNotAllowed)
				return
			}
			b, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, "some problem, see in logs", http.StatusBadRequest)
				lg.Error(err)
				return
			}
			bodyMap := make(map[string]int, 1)
			if err := json.Unmarshal(b, &bodyMap); err != nil {
				http.Error(w, "some problem, see in logs", http.StatusBadRequest)
				lg.Error(err)
				return
			}
			addToRedis, exist := bodyMap["number"]
			if !exist {
				http.Error(w, "some problem, see in logs", http.StatusBadRequest)
				lg.Error("invalid json data")
				return
			}
			if err := setRedis(rdb, addToRedis); err != nil {
				http.Error(w, "some problem, see in logs", http.StatusBadRequest)
				lg.Error(err)
				return
			}
		default:
			http.NotFound(w, req)
			lg.Error("404")
			return
		}
	})

	lg.Info("start...")
	if err := http.ListenAndServe(":"+modelPort, mux); err != nil {
		lg.Fatal(err)
	}
}

func setRedis(rdb *redis.Client, number int) error {
	return rdb.Set(context.Background(), "key", number, 0).Err()
}
func getRedis(rdb *redis.Client) (int, error) {
	val, err := rdb.Get(context.Background(), "key").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}
