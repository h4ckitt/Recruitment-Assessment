package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Listening On Port localhost:9943....")

	if err := http.ListenAndServe(":9943", nil); err != nil {
		log.Fatalln(err)
	}
}
