package Mod

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

/////////TO DO LIST/////////////////////
// merge
// search bar
// bouton filters

type Artist struct {
	Id             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
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

func GetOneArtistInfo(name string) Artist {
	art, _ := GetData()
	var artInfo Artist
	for i := range art {
		if art[i].Name == name {
			artInfo.Image = art[i].Image
			artInfo.Name = art[i].Name
			artInfo.Members = art[i].Members
			artInfo.CreationDate = art[i].CreationDate
			artInfo.FirstAlbum = art[i].FirstAlbum
			artInfo.DatesLocations = art[i].DatesLocations
		}
	}
	return artInfo
}

func RightFormForDate(date string) string {
	date = strings.ReplaceAll(date, "/", "-")
	return date
}

func memberMatch(query string, members []string) (bool, string) {
	memberName := ""
	for _, member := range members {
		if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
			memberName = member
			return true, memberName
		}
	}
	return false, memberName
}

func locationMatch(query string, datesLocations map[string][]string) (bool, string) {
	locate := ""
	for location := range datesLocations {
		if strings.Contains(strings.ToLower(string(location)), strings.ToLower(query)) {
			locate = string(location)
			return true, locate
		}
	}
	return false, locate
}

func SearchOptions(query string, data []Artist) []string {
	var optionsSearchBar []string

	contains := func(array []string, item string) bool {
		for _, element := range array {
			if element == item {
				return true
			}
		}
		return false
	}

	if query != "" {
		for _, artist := range data {
			mbrBool, memberName := memberMatch(query, artist.Members)
			locatbool, locate := locationMatch(query, artist.DatesLocations.DatesLocations)

			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) && !contains(optionsSearchBar, artist.Name) {
				optionsSearchBar = append(optionsSearchBar, artist.Name)
			} else if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(RightFormForDate(query))) && !contains(optionsSearchBar, artist.Name) {
				optionsSearchBar = append(optionsSearchBar, artist.Name)
			} else if strings.Contains(fmt.Sprintf("%d", artist.CreationDate), query) && !contains(optionsSearchBar, artist.Name) {
				optionsSearchBar = append(optionsSearchBar, artist.Name)
			} else if locatbool && !contains(optionsSearchBar, locate) {
				optionsSearchBar = append(optionsSearchBar, locate)
			} else if mbrBool && !contains(optionsSearchBar, memberName+" - Member") {
				optionsSearchBar = append(optionsSearchBar, memberName+" - Member")
			}
		}
	}
	return optionsSearchBar
}

func SearchBar(query string, data []Artist) []Artist {
	var filteredArtists []Artist
	if query != "" {
		for _, artist := range data {

			mbrBool, _ := memberMatch(query, artist.Members)
			locatbool, _ := locationMatch(query, artist.DatesLocations.DatesLocations)

			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
				filteredArtists = append([]Artist{artist}, filteredArtists...)
			} else if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(RightFormForDate(query))) ||
				strings.Contains(fmt.Sprintf("%d", artist.CreationDate), query) ||
				locatbool ||
				mbrBool {
				filteredArtists = append(filteredArtists, artist)

			}
		}
	} else {
		filteredArtists = data
	}
	return filteredArtists
}
