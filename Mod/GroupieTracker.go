package Mod

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

func RightFormForDate(date string) string {
	date = strings.ReplaceAll(date, "/", "-")
	return date
}

func GetBool(query string, datesLocations map[string][]string, members []string) (bool, string) {
	value := ""
	if members != nil {
		for _, member := range members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
				value = member
				return true, value
			}
		}
	} else {
		for location := range datesLocations {
			if strings.Contains(strings.ToLower(string(location)), strings.ToLower(query)) {
				value = string(location)
				return true, value
			}
		}

	}
	return false, value
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
			mbrBool, memberName := GetBool(query, nil, artist.Members)
			locatbool, locate := GetBool(query, artist.DatesLocations.DatesLocations, nil)

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

func SearchBarCheckBox(values []string, query string, data []Artist) []Artist {
	var filteredArtists []Artist
	if len(values) == 0 {
		return SearchBar(query, data)
	}

	if len(values) == 1 {
		switch {
		case values[0] == "location":
			for _, artist := range data {
				locatbool, _ := GetBool(query, artist.DatesLocations.DatesLocations, nil)
				if locatbool {
					filteredArtists = append(filteredArtists, artist)
				}

			}
			return filteredArtists
		case values[0] == "members":
			for _, artist := range data {
				if query == string(strconv.Itoa(len(artist.Members))) {
					filteredArtists = append(filteredArtists, artist)
				}
			}
			return filteredArtists
		case values[0] == "first_album_year":
			for _, artist := range data {
				if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(RightFormForDate(query))) {
					filteredArtists = append(filteredArtists, artist)
				}
			}
			return filteredArtists
		case values[0] == "creation_date":
			for _, artist := range data {
				if strings.Contains(fmt.Sprintf("%d", artist.CreationDate), query) {
					filteredArtists = append(filteredArtists, artist)
				}
			}
			return filteredArtists
		}

	}

	if len(values) == 2 && Contains(values, "creation_date") && Contains(values, "first_album_year") {
		for _, artist := range data {
			if strings.Contains(fmt.Sprintf("%d", artist.CreationDate), query) ||
				strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(RightFormForDate(query))) {
				filteredArtists = append(filteredArtists, artist)
			}
		}
		return filteredArtists
	}
	return nil

}

func SearchBar(query string, data []Artist) []Artist {
	var filteredArtists []Artist
	if query != "" {
		for _, artist := range data {

			mbrBool, _ := GetBool(query, nil, artist.Members)
			locatbool, _ := GetBool(query, artist.DatesLocations.DatesLocations, nil)

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

func RemoveFromCheckedOptions(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
