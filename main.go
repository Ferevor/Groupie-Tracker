package main

import (
	"encoding/json"
	"fmt"
	"groupie/Mod"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type PageData struct {
	Query            string
	Data             []Mod.Artist
	OptionsSearchBar []string
	CheckedOptions   []string
	No_results       bool
}

var checkedOptions = []string{}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	displayQuery := r.URL.Query().Get("displayQuery")

	data, err := Mod.GetData()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	if start != "" && end != "" {
		startYear, err := strconv.Atoi(start)
		if err != nil {
			log.Printf("Error converting start year to int: %v", err)
			http.Error(w, "Error converting start year to int", http.StatusBadRequest)
			return
		}
		endYear, err := strconv.Atoi(end)
		if err != nil {
			log.Printf("Error converting end year to int: %v", err)
			http.Error(w, "Error converting end year to int", http.StatusBadRequest)
			return
		}

		var filteredArtists []Mod.Artist
		switch r.URL.Query().Get("filter") {
		case "CreationDate":
			for _, artist := range data {
				if artist.CreationDate >= startYear && artist.CreationDate <= endYear {
					filteredArtists = append(filteredArtists, artist)
				}
			}
			data = filteredArtists
		case "FirstAlbum":
			for _, artist := range data {
				firstAlbumYear := strings.Split(artist.FirstAlbum, "-")
				firstAlbumYearInt, err := strconv.Atoi(firstAlbumYear[2])
				if err != nil {
					log.Printf("Error converting first album year to int: %v", err)
					http.Error(w, "Error converting first album year to int", http.StatusBadRequest)
					return
				}
				if firstAlbumYearInt >= startYear && firstAlbumYearInt <= endYear {
					filteredArtists = append(filteredArtists, artist)
				}
			}
			data = filteredArtists
		}
	}
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "asc" {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Name < data[j].Name
		})
	} else if sortOrder == "desc" {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Name > data[j].Name
		})
	}

	filteredArtists := Mod.SearchBarCheckBox(checkedOptions, query, data)
	optionsSearchBar := Mod.SearchOptions(query, data)

	var no_results bool
	if len(filteredArtists) == 0 {
		no_results = true
	}

	pageData := PageData{
		Query:            displayQuery,
		Data:             filteredArtists,
		OptionsSearchBar: optionsSearchBar,
		CheckedOptions:   checkedOptions,
		No_results:       no_results,
	}

	t := template.New("GroupTra.tmpl").Funcs(template.FuncMap{
		"contains": Mod.Contains,
	})
	t, err = t.ParseFiles("Templates/GroupTra.tmpl")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, pageData)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func getCheckedOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"checkedOptions": checkedOptions})
}

func updateCheckedOptionsHandler(w http.ResponseWriter, r *http.Request) {
	var selection struct {
		Option    string `json:"option"`
		IsChecked bool   `json:"isChecked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&selection); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	if selection.IsChecked {
		if !Mod.Contains(checkedOptions, selection.Option) {
			checkedOptions = append(checkedOptions, selection.Option)
		}
	} else {
		checkedOptions = Mod.RemoveFromCheckedOptions(checkedOptions, selection.Option)
	}

	w.WriteHeader(http.StatusOK)
}

func searchOptionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	data, err := Mod.GetData()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	optionsSearchBar := Mod.SearchOptions(query, data)

	if len(optionsSearchBar) > 5 {
		optionsSearchBar = optionsSearchBar[:5]
	}

	w.Header().Set("Content-Type", "text/html")
	for _, option := range optionsSearchBar {
		fmt.Fprintf(w, "<option value=\"%s\">%s</option>", option, option)
	}
}

func openBrowser(url string) {
	var err error
	switch os := runtime.GOOS; os {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

func main() {
	http.Handle("/Styles/", http.StripPrefix("/Styles/", http.FileServer(http.Dir("Styles"))))
	http.Handle("/Scripts/", http.StripPrefix("/Scripts/", http.FileServer(http.Dir("Scripts"))))
	http.Handle("/Images/", http.StripPrefix("/Images/", http.FileServer(http.Dir("Images"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/get-checked-options", getCheckedOptionsHandler)
	http.HandleFunc("/update-checked-options", updateCheckedOptionsHandler)
	http.HandleFunc("/search", searchOptionsHandler)
	fmt.Println("Starting server at port 8080")
	go openBrowser("http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
