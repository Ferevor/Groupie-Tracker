package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Location     string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Id             int      `json:"id"`
	DatesLocations []string `json:"datesLocations"`
}

func GetInfo(url string) []byte {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func GetArtist() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/artists")

	var artists []Artist
	err := json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	for _, artist := range artists {
		fmt.Println(artist.Name)
	}
}

func GetLocations() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/locations")

	var locations []Locations
	err := json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	for _, location := range locations {
		fmt.Println(location.Dates)
	}
}

func GetDates() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/dates")

	var dates []Dates
	err := json.Unmarshal(body, &dates)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	for _, date := range dates {
		fmt.Println(date.Dates)
	}
}

func GetRelation() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/dates")

	var relation []Relation
	err := json.Unmarshal(body, &relation)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	for _, relate := range relation {
		fmt.Println(relate.DatesLocations)
	}
}

func main() {
	GetArtist()
}
