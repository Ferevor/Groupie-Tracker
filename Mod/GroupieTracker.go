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

func GetData() ([]Artist, error) {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/artists")

	var artists []Artist

	err := json.Unmarshal(body, &artists)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %v", err)
	}

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
