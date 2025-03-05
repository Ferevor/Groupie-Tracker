//package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var items = []string{"Paris", "Marseille", "Lyon", "Toulouse", "Nice", "Nantes", "Strasbourg", "Montpellier", "Bordeaux", "Lille"}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	var results []string
	for _, item := range items {
		if len(query) > 0 && containsIgnoreCase(item, query) {
			results = append(results, item)
		}
	}
	json.NewEncoder(w).Encode(results)
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	http.HandleFunc("/search", searchHandler)
	fmt.Println("Serveur en écoute sur le port 8080...")
	http.ListenAndServe(":8080", nil)
}
