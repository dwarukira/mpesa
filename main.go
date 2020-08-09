package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dwarukira/mpesa/transactions"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Transaction).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type Message struct {
	Message string `json:"message"`
}

func Transaction(w http.ResponseWriter, r *http.Request) {
	var message Message
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	if err := r.Body.Close(); err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	if err := json.Unmarshal(body, &message); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			// panic(err)
			json.NewEncoder(w).Encode(err)
			return
		}
	}
	log.Println(message.Message)
	t, err := transactions.NewFromMessage(message.Message)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(t); err != nil {
		// panic(	err)
		json.NewEncoder(w).Encode(err)
		return
	}
}
