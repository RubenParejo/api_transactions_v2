package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"api_transactions_v2/pkg/model"
)

var (
	transactions sync.Map
)

func LaunchAPI() {
	http.HandleFunc("/transactions", handleTransactions)
	port := ":8080"
	fmt.Printf("Server listening in port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("LaunchAPI - Error starting server: %s", err)
	}
}

func handleTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostTransaction(w, r)
	case http.MethodGet:
		handleGetTransaction(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostTransaction(w http.ResponseWriter, r *http.Request) {
	var data model.Data

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Printf("handlePostTransaction - Error reading body: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Printf("handlePostTransaction - Error unmarshaling body: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transactions.Store(data.Id, data)
	result := model.Response{Status: "Success"}

	response, err := json.Marshal(result)
	if err != nil {
		log.Printf("handlePostTransaction - Error marshaling result: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	id := query.Get("id")
	if id == "" {
		err := errors.New("param id is missing")
		log.Printf("handleGetTransaction - Error getting id from request: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, exists := transactions.Load(id)
	if !exists {
		err := errors.New("id not found")
		log.Printf("handleGetTransaction - Error getting data: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("handleGetTransaction - Error marshaling result: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
