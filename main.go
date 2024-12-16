package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// URL de l'API
	url := "https://groupietrackers.herokuapp.com/api"

	// Faire une requête GET
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erreur lors de la requête GET: %v", err)
	}
	defer resp.Body.Close()

	// Lire le corps de la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du corps de la réponse: %v", err)
	}

	// Afficher la réponse
	fmt.Println(string(body))
}
