package Mod

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Artist struct {
	Id             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	Location       string   `json:"location"`
	ConcertDates   string   `json:"concertDates"`
	Relations      string   `json:"relations"`
	DatesLocations Relation
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func GetInfo(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func GetArtist() ([]Artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la requête GET: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture du corps de la réponse: %v", err)
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors du déchiffrement du JSON: %v", err)
	}

	return artists, nil
}

func ArtistInfo(name string) {

}
