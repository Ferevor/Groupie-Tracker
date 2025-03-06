package main

import (
	"fmt"
	"groupie/Mod"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sort"
)

type PageData struct {
	Query            string
	Artists          []Mod.Artist
	OptionsSearchBar []string
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	displayQuery := r.URL.Query().Get("displayQuery")

	data, err := Mod.GetData()
	if err != nil {
		log.Printf("Erreur lors de la récupération des données: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		return
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

	filteredArtists := Mod.SearchBar(query, data)
	optionsSearchBar := Mod.SearchOptions(query, data)

	pageData := PageData{
		Query:            displayQuery,
		Artists:          filteredArtists,
		OptionsSearchBar: optionsSearchBar,
	}

	t, err := template.ParseFiles("Templates/grouptra.tmpl")
	if err != nil {
		log.Printf("Erreur lors du chargement du template: %v", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, pageData)
	if err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		return
	}
}

func searchOptionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	data, err := Mod.GetData()
	if err != nil {
		log.Printf("Erreur lors de la récupération des données: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
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
	http.HandleFunc("/", handler)
	http.HandleFunc("/search", searchOptionsHandler)
	fmt.Println("Starting server at port 8080")
	go openBrowser("http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
