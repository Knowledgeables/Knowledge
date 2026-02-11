package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Create a new server instance
	http.HandleFunc("/", func(writer http.ResponseWriter, response *http.Request) {
		// Handle the request and respond with knowledge data
		fmt.Println(writer, "Hello Http!")
		writer.Write([]byte("Welcome to the Knowledgeable server!"))
	})

	// Start the server on port 8080
	http.ListenAndServe(":8080", nil)
}