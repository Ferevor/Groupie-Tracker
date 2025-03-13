package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Struct to capture user selection
type UserSelection struct {
	Option string `json:"option"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test.tmpl")
	})

	http.HandleFunc("/option", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var selection UserSelection
			err := json.NewDecoder(r.Body).Decode(&selection)
			if err != nil {
				http.Error(w, "Invalid data", http.StatusBadRequest)
				return
			}

			fmt.Printf(selection.Option) // c pour afficher les options qui sont touchees

			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
