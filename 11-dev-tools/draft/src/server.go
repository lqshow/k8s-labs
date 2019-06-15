package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Skaffold!!\n")
}

func handleRequests() {
	log.Print("Service is listening on port 30081")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	handleRequests()
}
