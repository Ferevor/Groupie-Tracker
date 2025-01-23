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

	fmt.Println("DatesLocations:")
	for key, value := range relation.DatesLocations {
		fmt.Println("     " + string(key))
		for i := 0; i < len(value); i++ {
			fmt.Println("          " + string(value[i]))
		}
		fmt.Println()
	}

}

func main() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/artists")

	var artists []Artist
	err := json.Unmarshal(body, &artists)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	for _, artist := range artists {
		fmt.Println("----------------------------")
		fmt.Println("ID:")
		fmt.Println("     ", artist.Id, "\n")
		fmt.Println("Name:")
		fmt.Println("     ", artist.Name, "\n")
		fmt.Println("Image:")
		fmt.Println("     ", artist.Image, "\n")
		fmt.Println("Members:")
		for i := 0; i < len(artist.Members); i++ {
			fmt.Println("     ", artist.Members[i])
		}
		fmt.Println()
		fmt.Println("Creation date:")
		fmt.Println("     ", artist.CreationDate, "\n")
		fmt.Println("First album:")
		fmt.Println("     ", artist.FirstAlbum, "\n")
		getRelation(artist.Relations)
		fmt.Println()
	}
}
