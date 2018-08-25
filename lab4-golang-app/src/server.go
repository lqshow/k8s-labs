package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func homePage(w http.ResponseWriter, r *http.Request){
	host, err := os.Hostname()
	check(err)
    fmt.Fprintf(w, "Welcome to the HomePage! \n\nHostname is %s", host)

    fmt.Println("Endpoint Hit: homePage")
}

func check(err error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func handleRequests() {
    http.HandleFunc("/", homePage)
    log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
    handleRequests()
}