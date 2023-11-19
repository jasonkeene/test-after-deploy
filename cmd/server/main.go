package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	log.Fatal(http.ListenAndServe(":12345", nil))
}
