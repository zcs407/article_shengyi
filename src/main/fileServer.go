package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("fileServer running, access address http://0.0.0.0:8888")
	err := http.ListenAndServe(":8888", http.FileServer(http.Dir("/Users/ander/go/src/articlebk/articleData")))
	if err != nil {
		log.Println("[fileServer running error:]", err)
	}
}
