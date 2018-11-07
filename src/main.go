package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type myData struct {
	Msg         string
	Status      int
	Description string
}

func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/message", handleQryMessage).Methods("GET") //ex http://localhost:3000/message?msg=Hello%20World
	router.HandleFunc("/m/{msg}", handleURLMessage).Methods("GET")
	router.HandleFunc("/monster", handleURLMonster).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func handleQryMessage(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	message := vars.Get("msg")

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func handleURLMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	message := vars["msg"]

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}

func handleURLMonster(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var data myData
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	log.Println(data.Msg)
	log.Println(data.Status)
	// message := data.Msg
	// json.NewEncoder(w).Encode(map[string]string{"name": message})
	data.Description = "What the duck"

	dataJSON, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataJSON)
}
