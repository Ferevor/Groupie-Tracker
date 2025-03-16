pack age main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ///////////////////DEBUT/////////////////////////////////////////
// Struct pour capturer les options de l'utilisateur
type UserSelection struct {
	Option string `json:"option"`
}

var checkedOptions = []string{}

///////////////////////FIN///////////////////////////////////////

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test.tmpl")
	})

	http.HandleFunc("/option", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			//////////////////////////DEBUT////////////////////////////////////
			var selection UserSelection
			err := json.NewDecoder(r.Body).Decode(&selection)
			if err != nil {
				http.Error(w, "Données invalides", http.StatusBadRequest)
				return
			}

			if selection.Option != "" { // Vérifie si l'option n'est pas vide
				if contains(checkedOptions, selection.Option) {
					checkedOptions = remove(checkedOptions, selection.Option)
				} else {
					checkedOptions = append(checkedOptions, selection.Option)
				}
			}

			// Affiche toutes les options sélectionnées
			fmt.Println("Options cochées :", checkedOptions)
			//// qui devo fare in modo di verificare prima di chiamare la funzione per la searchbar
			//// se la richiesta non fa parte (tipo sbagliato o altro) allora a array data è vuota
			//// anche la searchbar suggestions sono toccate
			////
			//////////////////////////FIN////////////////////////////////////

		} else {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Serveur en cours d'exécution sur http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur : %v\n", err)
	}
}

// ////////////////////////DEBUT////////////////////////////////////
// Fonction auxiliaire pour vérifier la présence d'une valeur dans une slice
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Fonction auxiliaire pour supprimer une valeur d'une slice
func remove(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

////////////////////////////FIN//////////////////////////////////
