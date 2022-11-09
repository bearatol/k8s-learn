package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bearatol/lg"
	"github.com/go-redis/redis/v9"
)

func main() {
	redisAddr, ok := os.LookupEnv("REDIS_PORT")
	if !ok {
		redisAddr = "localhost:6001"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "123", // no password set
		DB:       0,     // use default DB
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/add" {
			http.NotFound(w, req)
			lg.Error("404")
			return
		}

		num, err := getRedis(rdb)
		if err != nil {
			http.Error(w, "some problem", http.StatusBadRequest)
			lg.Error(err)
			return
		}
		lg.Infof("a guest comes, counter: %d", num)
		if err := setRedis(rdb, (num + 1)); err != nil {
			http.Error(w, "some problem", http.StatusBadRequest)
			lg.Error(err)
			return
		}

		fmt.Fprintf(w, "Herlou epti!")
	})

	lg.Info("go-go-go...")
	if err := http.ListenAndServe("0.0.0.0:6002", mux); err != nil {
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
