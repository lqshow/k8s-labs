package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func init() {
	for _, v := range os.Environ() {
		fmt.Println(v)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	check(err)
	fmt.Fprintf(w, "Welcome to the HomePage! \n\nHostname is %s", host)

	redisHost := os.Getenv("REDIS_SVC_SERVICE_HOST")
	redisPort := os.Getenv("REDIS_SVC_SERVICE_PORT")

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	pong, err := client.Ping().Result()

	fmt.Println(pong, err)
	fmt.Println("Endpoint Hit: homePage")
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}
