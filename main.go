package main

import (
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s", "world");
	})

	http.ListenAndServe(":8000", nil)
}