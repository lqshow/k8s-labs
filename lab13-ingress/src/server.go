package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var port = os.Getenv("PORT")

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This service is listening on port %s", port)
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}
