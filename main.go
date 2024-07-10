package main

import (
	"fmt"
	"log"
	"net/http"

	"tracker/getapi"
)

func main() {
	http.HandleFunc("/", getapi.HomeHandler)
	http.HandleFunc("/individual", getapi.IndividualHandler)
	http.HandleFunc("/Search", getapi.Search)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	fmt.Printf("Starting server for testing HTTP POST on http://localhost:8080 ...\n")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}
