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

func getRelation(url string) {
	body := GetInfo(url)

	var relation Relation
	err := json.Unmarshal(body, &relation)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	fmt.Println(relation.DatesLocations)

}

// exemple pour avoir les info sur les API de dates et location
func getDates(url string) {
	body := GetInfo(url)

	var date Dates
	err := json.Unmarshal(body, &date)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	fmt.Println("Name:", date.Dates)
}

func main() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/artists")

	var artists []Artist
	err := json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	var nbr int
	for _, artist := range artists {
		nbr += 1
		fmt.Println("----------------------------")
		fmt.Println("ID:", artist.Id)
		fmt.Println("Name:", artist.Name)
		fmt.Println("Image:", artist.Image)
		fmt.Println("Members:", artist.Members)
		fmt.Println("Creation date:", artist.CreationDate)
		fmt.Println("First album:", artist.FirstAlbum)
		getRelation(artist.Relations)
	}
}
