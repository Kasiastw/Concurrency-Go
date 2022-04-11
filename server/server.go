package main

import (
	"flag"
	"log"
	"net/http"
	"rate-limiting/limiter"
	"time"
)

func main() {
	address := flag.String("address", ":8080", "Server address")
	window := flag.Duration("window", time.Second*10, "Length of observation window")
	limit := flag.Int("limit", 10, "Request limit")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", mainHandle)
	mux.HandleFunc("/limit", limiter.LimitHandler(*window, *limit, finalHandle))

	log.Printf("Listening on %s...\n", *address)
	err := http.ListenAndServe(*address, mux)
	log.Fatal(err)

}

func mainHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is an application to limit the number of requests"))
}

func finalHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("If you see this, it means that you have not yet used up your limit"))
}
