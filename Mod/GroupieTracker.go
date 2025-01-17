package Mod

import (
	"encoding/json"
	"fmt"
	"io"
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
	Dates     Dates    `json:"dates"`
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
		return nil, fmt.Errorf("Erreur lors du déchiffrement des données: %v", err)
	}

	return artists, nil
}

// func GetArtist() {
// 	body := GetInfo("https://groupietrackers.herokuapp.com/api/artists")

// 	var artists []Artist
// 	err := json.Unmarshal(body, &artists)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling data: %v", err)
// 	}

// 	for _, artist := range artists {
// 		// Convert the artist struct to a JSON string with indentation
// 		artistJSON, err := json.MarshalIndent(artist, "", "  ")
// 		if err != nil {
// 			log.Fatalf("Error marshaling data: %v", err)
// 		}
// 		fmt.Println(string(artistJSON))
// 		// fmt.Println("------------------------------")
// 	}
// 	return artists, nil
// }

func GetLocations() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/locations")

	// Unmarshal into a single Locations struct to inspect the structure
	var locations LocationData
	err := json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	locationJSON, err := json.MarshalIndent(locations, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
	}
	fmt.Println(string(locationJSON))
	fmt.Println("------------------------------")
}

func GetDates() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/dates")

	fmt.Println("Raw JSON Response:")
	fmt.Println(string(body))

	var dates Dates
	err := json.Unmarshal(body, &dates)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	dateJSON, err := json.MarshalIndent(dates, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
	}
	fmt.Println(string(dateJSON))
	fmt.Println("------------------------------")
}

func GetRelation() {
	body := GetInfo("https://groupietrackers.herokuapp.com/api/relation")

	fmt.Println("Raw JSON Response:")
	fmt.Println(string(body))

	var relation Relation
	err := json.Unmarshal(body, &relation)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	// Convert the relation struct to a JSON string with indentation
	relationJSON, err := json.MarshalIndent(relation, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
	}
	fmt.Println(string(relationJSON))
	fmt.Println("------------------------------")
}

//	func main() {
//		GetArtist()
//		//GetLocations()
//		//GetDates()
//		//GetRelation()
//	}
