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

type LocationData struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
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
		if artist.Location[:len(artist.Location)-3] == "https://groupietrackers.herokuapp.com/api/locations" {
			getLocation(artist.Location)
		}
		if artist.ConcertDates[:len(artist.ConcertDates)-3] == "https://groupietrackers.herokuapp.com/api/dates" {
			getDates(artist.ConcertDates)
		}
		if artist.Relations[:len(artist.Relations)-3] == "https://groupietrackers.herokuapp.com/api/relation" {
			getRelation(artist.Relations)
		}
		artistJSON, err := json.MarshalIndent(artist, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling data: %v", err)
		}
		fmt.Println(string(artistJSON))
		fmt.Println("------------------------------")

	}
}

func getLocation(url string) {
	body := GetInfo(url)
	fmt.Println(string(body))

}

func getDates(url string) {
	body := GetInfo(url)
	fmt.Println(string(body))
}

func getRelation(url string) {
	body := GetInfo(url)
	fmt.Println(string(body))
}

func main() {
	GetArtist()
	//GetLocations()
	//GetDates()
	//GetRelation()
}
