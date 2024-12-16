package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type artists struct {
	id            int
	image         string
	name          string
	members       []string
	creationDates int
	firstAlbum    string
	location      locations `json: "location"`
	concertDates  dates     `json: "concertDates"`
	relations     relation  `json: "relations"`
}

type locations struct {
	id        int
	locations []string
	dates     dates `json: "dates"`
}

type dates struct {
	id    int
	dates []string
}

type relation struct {
	id             int
	datesLocations []string `json: "datesLocations"`
}

func GetLocations() {
	res := GetAPI("https://groupietrackers.herokuapp.com/api/locations")
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var location locations
	json.Unmarshal(body, &location)

}

func GetDates() {
	res := GetAPI("https://groupietrackers.herokuapp.com/api/dates")
	fmt.Println(res)
}

func GetArtist() {
	res := GetAPI("https://groupietrackers.herokuapp.com/api/artists")
	fmt.Println(res)
}

func GetRelation() {
	res := GetAPI("https://groupietrackers.herokuapp.com/api/dates")
	fmt.Println(res)
}

func GetAPI(url string) *http.Response {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	return res
}

func main() {

}
