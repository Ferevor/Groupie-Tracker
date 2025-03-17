package Mod

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Definie la structure Artist qui représente un artiste avec ses informations, présentes dans l'API, en passant par du JSON
type Artist struct {
	Id             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	Relations      string   `json:"relations"`
	DatesLocations Relation
}

// Definie la structure Relation qui prend les informations de l'API relation, ayant pour but de relier les API date et lieu des concerts
type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// Fonction qui récupère toutes les informations depuis une URL
func GetInfo(url string) []byte {
	res, err := http.Get(url) // Effectue une requête HTTP GET
	if err != nil {           // Gère les erreurs de manière critique
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body) // Lis le contenu
	if err != nil {
		log.Fatal(err)
	}
	return body // Renvoie le contenu
}

// Fonction qui récupère les données de tous les artistes depuis l'API et renvoie un tableau avec toutes leurs informations
func GetData() ([]Artist, error) {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/artists") // Ouvre l'API artists

	var artists []Artist // Initialise le tableau

	err := json.Unmarshal(body, &artists) // Décodage des données JSON en tableau d'objets Artist
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %v", err)
	}

	// Parcourt les artistes pour récupérer leurs relations et les associer dans artist.DatesLocations
	for i := range artists {
		relationBody := GetInfo(artists[i].Relations)
		var relation Relation
		err := json.Unmarshal(relationBody, &relation)
		if err != nil {
			log.Fatalf("Error unmarshaling relation data: %v", err)
		}
		artists[i].DatesLocations = relation
	}
	return artists, nil
}

// Fonction pour remplacer les slashs par des tirets dans les dates entrées par les utilisateurs dans la bar de recherche
// pour faire ensuite la comparaison avec les informations sur les artistes (où les dates sont écritent jj-mm-dddd)
func RightFormForDate(date string) string {
	date = strings.ReplaceAll(date, "/", "-")
	return date
}

// Fonction qui renvoie vrai si la recherche du client est présente parmis les membres ou les lieux et dates des concerts
func GetBool(query string, datesLocations map[string][]string, members []string) (bool, string) {
	value := ""
	if members != nil { // Cas des membres
		for _, member := range members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
				value = member
				return true, value
			}
		}
	} else { // Cas des lieux et dates des concerts
		for location := range datesLocations {
			if strings.Contains(strings.ToLower(string(location)), strings.ToLower(query)) {
				value = string(location)
				return true, value
			}
		}

	}
	return false, value // Aucun résultat a été trouvé
}

// Fonction qui propose un tableau de suggestions par rapport à la recherche du client
func SearchOptions(query string, data []Artist) []string {
	var optionsSearchBar []string // Initialise le tableau qui va stocker le résultat de la recherche

	if query != "" {
		for _, artist := range data { // Parcourt tous les artists
			// Utilise la fonction GetBool() pour savoir si la recherche est parmis les membres
			mbrBool, memberName := GetBool(query, nil, artist.Members)
			// Pareil mais pour les lieux et dates de concerts
			locatbool, locate := GetBool(query, artist.DatesLocations.DatesLocations, nil)

			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) && !Contains(optionsSearchBar, artist.Name) {
				optionsSearchBar = append(optionsSearchBar, artist.Name)
			} else if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(RightFormForDate(query))) && !Contains(optionsSearchBar, artist.Name) {
				optionsSearchBar = append(optionsSearchBar, artist.Name)
			} else if strings.Contains(fmt.Sprintf("%d", artist.CreationDate), query) && !Contains(optionsSearchBar, artist.Name) {
				optionsSearchBar = append(optionsSearchBar, artist.Name)
			} else if locatbool && !Contains(optionsSearchBar, locate) {
				optionsSearchBar = append(optionsSearchBar, locate)
			} else if mbrBool && !Contains(optionsSearchBar, memberName+" - Member") {
				// On ajoute - Member pour distinguer le membre du nom du groupe
				optionsSearchBar = append(optionsSearchBar, memberName+" - Member")
			}
		}
	}
	return optionsSearchBar
}

// Fonction qui filtre les artistes selon les critères dans valeurs et la requête
func SearchBarCheckBox(values []string, query string, data []Artist) []Artist {
	var filteredArtists []Artist // Tableau pour stocker les artists filtrés
	if len(values) == 0 {        // S'il n'y a aucune condition alors on renvoye les résultats de la recherche normale
		return SearchBar(query, data) // Pour cela on fait appel à la fonction SearchBar()
	}

	if len(values) == 1 { // S'il n'y a qu'une valeur, plusieurs cas sont possibles
		switch {
		case values[0] == "location": // Cas où l'utilisateur veux filtrer les lieux des concerts
			for _, artist := range data {
				locatbool, _ := GetBool(query, artist.DatesLocations.DatesLocations, nil)
				if locatbool {
					filteredArtists = append(filteredArtists, artist)
				}

			}
			return filteredArtists
		case values[0] == "members": // Cas où où l'utilisateur veux filtrer le nombre de membres
			for _, artist := range data {
				if query == string(strconv.Itoa(len(artist.Members))) {
					filteredArtists = append(filteredArtists, artist)
				}
			}
			return filteredArtists
		case values[0] == "first_album_year": // Cas où l'utilisateur veux filtrer l'année du premier album
			for _, artist := range data {
				if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(RightFormForDate(query))) {
					filteredArtists = append(filteredArtists, artist)
				}
			}
			return filteredArtists
		case values[0] == "creation_date": // Cas où l'utilisateur veux filtrer la date de création
			for _, artist := range data {
				if strings.Contains(fmt.Sprintf("%d", artist.CreationDate), query) {
					filteredArtists = append(filteredArtists, artist)
				}
			}
			return filteredArtists // renvoie le tableau complété
		}

	}

	// Si deux valeurs sont spécifiées, on applique la seule possiblité car seule la date de création et la date du premier album sont compatibles
	if len(values) == 2 && Contains(values, "creation_date") && Contains(values, "first_album_year") {
		for _, artist := range data {
			if strings.Contains(fmt.Sprintf("%d", artist.CreationDate), query) ||
				strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(RightFormForDate(query))) {
				// Si la requête est présente dans une des deux alors l'artist est ajouté à filteredArtists
				filteredArtists = append(filteredArtists, artist)
			}
		}
		return filteredArtists
	}
	return nil

}

// Fonction qui filtre les artistes en fonction de la requête
func SearchBar(query string, data []Artist) []Artist {
	var filteredArtists []Artist
	if query != "" { // Vérifie si la requête n'est pas vide
		for _, artist := range data {

			mbrBool, _ := GetBool(query, nil, artist.Members)                         // Recherche parmi les membres du groupe
			locatbool, _ := GetBool(query, artist.DatesLocations.DatesLocations, nil) // / Recherche parmi les lieux des concerts

			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
				// Si le nom de l'artiste contient la requête alors l'artist est ajouté au début de la liste
				filteredArtists = append([]Artist{artist}, filteredArtists...)
			} else if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(RightFormForDate(query))) ||
				strings.Contains(fmt.Sprintf("%d", artist.CreationDate), query) ||
				locatbool ||
				mbrBool {
				// Ajoute l'artiste à la liste si une des conditions est vrai
				filteredArtists = append(filteredArtists, artist)

			}
		}
	} else { // Si la requête est vide tous les artistes sont renvoyés
		filteredArtists = data
	}
	return filteredArtists
}

// Fonction qui supprime une option spécifique du tableau des options de filtres
func RemoveFromCheckedOptions(checkboxarray []string, value string) []string {
	for i, v := range checkboxarray { // Parcourt tous les éléments du tableau
		if v == value { // Vérifie si l'élément correspond à la valeur à supprimer
			return append(checkboxarray[:i], checkboxarray[i+1:]...) // Supprime l'élément en réassemblant le tableau sans lui
		}
	}
	return checkboxarray
}

// Fonction qui vérifie si une valeur spécifique est présente dans un tableau
func Contains(slice []string, value string) bool {
	for _, item := range slice { // Parcourt tous les éléments du tableau
		if item == value { // si l'élément correspond à la valeur recherchée on renvoie vrai
			return true
		}
	}
	return false
}
