package main

import (
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
	location      locations
	concertDates  dates
	relations     relation
}

type locations struct {
	id        int
	locations []string
	dates     dates
}

type dates struct {
	id    int
	dates []string
}

type relation struct {
	id             int
	datesLocations []string
}

func main() {
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(body))
}
