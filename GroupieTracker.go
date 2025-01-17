package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Artist struct {
	ID            int      `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationDates int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
}

type Location struct {
	ID        int    `json:"id"`
	Locations string `json:"locations"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

func main() {
	var index int = 0
	var chemin string = "https://groupietrackers.herokuapp.com/api/artists"
	res, err := http.Get(chemin)
	if err != nil {
		log.Fatal("Erreur lors de la requête HTTP : ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erreur lors de la lecture du corps de la réponse : ", err)
	}

	var artists []Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		log.Fatal("Erreur lors du décodage JSON : ", err)
	}

	var locations []Location
	if err := json.Unmarshal(body, &locations); err != nil {
		log.Fatal("Erreur lors du décodage JSON : ", err)
	}

	var dates []Date
	if err := json.Unmarshal(body, &dates); err != nil {
		log.Fatal("Erreur lors du décodage JSON : ", err)
	}

	for _, artist := range artists {
		chemin = "https://groupietrackers.herokuapp.com/api/artists"
		fmt.Println("----------------------------")
		fmt.Println("ID:", artist.ID)
		fmt.Println("Name:", artist.Name)
		fmt.Println("Image:", artist.Image)
		fmt.Println("Members:", artist.Members)
		fmt.Println("Creation date:", artist.CreationDates)
		fmt.Println("First album:", artist.FirstAlbum)
		chemin = "https://groupietrackers.herokuapp.com/api/locations"
		for _, location := range locations {
			fmt.Println("Locations:", location.Locations[3])
		}
		chemin = "https://groupietrackers.herokuapp.com/api/dates"
		date := dates[index]
		fmt.Println("Dates", date.Dates)
		index += 1
	}
}

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// type Response struct {
// 	Artist []Artists `json:"https://groupietrackers.herokuapp.com/api/artists"`
// }

// type Artists struct {
// 	Id            int      `json:"id"`
// 	Image         string   `json:"image"`
// 	Name          string   `json:"name"`
// 	Members       []string `json:"members"`
// 	CreationDates int      `json:"creationDate"`
// 	FirstAlbum    string   `json:"firstAlbum"`
// 	Location      locations
// 	ConcertDates  dates
// 	Relations     relation
// }

// type locations struct {
// 	id        int
// 	locations []string
// 	dates     dates
// }

// type dates struct {
// 	id    int
// 	dates []string
// }

// type relation struct {
// 	id             int
// 	datesLocations []string
// }

// func main() {
// 	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var response Response
// 	json.Unmarshal(body, &response)

// 	for i, p := range response.Artist {
// 		fmt.Println("----------------------------")
// 		fmt.Println("Name : ", p.name)
// 		i += 1
// 	}
// }
