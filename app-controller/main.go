package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/bearatol/lg"
)

func main() {
	modelPort, ok := os.LookupEnv("MODEL_PORT")
	if !ok {
		lg.Fatal("MODEL_PORT environment variable is not set")
	}
	controllerPort, ok := os.LookupEnv("CONTROLLER_PORT")
	if !ok {
		lg.Fatal("CONTROLLER_PORT environment variable is not set")
	}

	model := NewModel(modelPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/add" {
			http.NotFound(w, req)
			lg.Error("404")
			return
		}

		num, err := model.getRedis()
		if err != nil {
			http.Error(w, "some problem, see in logs", http.StatusBadRequest)
			lg.Error(err)
			return
		}
		lg.Info(num)
		model.setNumber = strconv.Itoa(num + 1)
		if err := model.setRedis(); err != nil {
			http.Error(w, "some problem, see in logs", http.StatusBadRequest)
			lg.Error(err)
			return
		}

		fmt.Fprintf(w, "number in redis: %s", model.setNumber)
	})

	lg.Info("start...")
	if err := http.ListenAndServe(":"+controllerPort, mux); err != nil {
		lg.Fatal(err)
	}
}

type Model struct {
	port,
	setNumber string
}

func NewModel(port string) *Model {
	return &Model{port: port}
}

func (m *Model) setRedis() error {
	responseBody := bytes.NewBuffer([]byte("{\"number\": " + m.setNumber + "}"))
	_, err := http.Post(fmt.Sprintf("http://0.0.0.0:%s/set", m.port), "application/json", responseBody)
	if err != nil {
		return err
	}
	return nil
}
func (m *Model) getRedis() (int, error) {
	res, err := http.Get(fmt.Sprintf("http://0.0.0.0:%s/get", m.port))
	if err != nil {
		return 0, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	bodyMap := make(map[string]int, 1)
	if err := json.Unmarshal(b, &bodyMap); err != nil {
		return 0, err
	}
	num, exist := bodyMap["number"]
	if !exist {
		return 0, fmt.Errorf("number is not exist")
	}
	return num, nil
}
