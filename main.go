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
	"strconv"
	"strings"
)

const tmpl = `
<!DOCTYPE html>
<html>
	<head>
    	<title>Artist Info</title>
		<link rel="stylesheet" type="text/css" href="/Styles/style.css">
	</head>
	<body>
		<div>
			<div class="header">
				<h1>Groupie Tracker</h1>
				<form method="GET" action="/">
					<select id="sort-select" name="sort">
                		<option value="">Sort By</option>
                		<option value="asc">Sort Ascending</option>
                	    <option value="desc">Sort Descending</option>
                	</select>
					<select id="filter-select" name="filter">
						<option value="">Filter By</option>
						<option value="CreationDate">Creation Date</option>
						<option value="FirstAlbum">First Album</option>
					</select>
					<input type="number" name="start" placeholder="De ">
					<input type="number" name="end" placeholder="à ">
					<button type="Submit">GO</button>

				</form>
			</div>
			<div class="box">
    			<form method="GET" action="/" id="searchForm">
        			<input type="text" class="input" list="suggestionsquery" id="optionsList" name="displayQuery" value="{{.Query}}" autocomplete="off" placeholder=" "/>
					<input type="hidden" id="searchQuery" name="query">
					<datalist id="suggestionsquery">
						{{range .OptionsSearchBar}}
                        	<option value="{{.}}">{{.}}</option>
                    {{end}}
					</datalist>
    			</form>
			</div>
			<div class="container">
				{{if .No_results}}
					<div class="no-results">
						No results found
					</div>
				{{end}}
				{{range .Data}}
            		<label for="modal-{{.Name}}" class="button">
						<div>
                			<img src="{{.Image}}" alt="Image" width="200" height="200">
						</div>
						<div>
                			<h2>{{.Name}}</h2>
						</div>
            		</label>
            		<input type="checkbox" id="modal-{{.Name}}" class="modal-toggle">
            		<div class="modal">
                		<div class="modal-content">
                    		<label for="modal-{{.Name}}" class="close">&times;</label>
							<div>
                   				<center> <h2>{{.Name}}</h2></center>
   							</div>
   							<div class ="image">
                				<center><img src="{{.Image}}" alt="Image" width="192" height="192"></center>
							</div>
							<p>Members: {{range .Members}}{{.}}, {{end}}</p>
							<p>Creation Date: {{.CreationDate}}</p>
							<p>First Album: {{.FirstAlbum}}</p>
							<p>Concert Location and Dates:</p>
        					{{range $key, $value := .DatesLocations.DatesLocations}}
							<p>{{$key}}</p>
          					<ul>
                				{{range $value}}
                  				<li>{{.}}</li>
                   				{{end}}
            				</ul>
        					{{end}}
                		</div>
        			</div>
				{{end}}
			</div>
		</div>
	</body>
</html>
`

type PageData struct {
	Query            string
	Data             []Mod.Artist
	OptionsSearchBar []string
	No_results       bool
}

func handler(w http.ResponseWriter, r *http.Request) {
	displayQuery := r.URL.Query().Get("displayQuery")
	query := r.URL.Query().Get("query")
	data, err := Mod.GetData()
	if err != nil {
		log.Printf("Erreur lors de la récupération des données: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		return
	}

	//log.Printf("Données récupérées: %v", data)

	sortOrder := r.URL.Query().Get("sort")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	// Filtrer les artistes par date de création si les paramètres sont présents
	if start != "" && end != "" { // Si les paramètres sont présents
		startYear, err := strconv.Atoi(start) // Convertir les paramètres en entiers
		if err != nil {
			log.Printf("Erreur lors de la conversion de l'année de début: %v", err)
			http.Error(w, "Erreur lors de la conversion de l'année de début", http.StatusInternalServerError)
			return
		}
		endYear, err := strconv.Atoi(end)
		if err != nil {
			log.Printf("Erreur lors de la conversion de l'année de fin: %v", err)
			http.Error(w, "Erreur lors de la conversion de l'année de fin", http.StatusInternalServerError)
			return
		}

		var filteredData []Mod.Artist
		switch r.URL.Query().Get("filter") {
		case "CreationDate":
			for _, artist := range data {
				if artist.CreationDate >= startYear && artist.CreationDate <= endYear { //si l'année de création est comprise entre les années de début et de fin
					filteredData = append(filteredData, artist) // Ajouter l'artiste à la liste des artistes filtrés
				}
			}
			data = filteredData // Remplacer les données par les données filtrées
			log.Printf("Données filtrées: %v", filteredData)
		case "FirstAlbum":
			for _, artist := range data {
				firstAlbumYear := strings.Split(artist.FirstAlbum, "-")   // Diviser la chaîne de l'année du premier album pour avoir seulement l'année
				firstAlbumYearInt, err := strconv.Atoi(firstAlbumYear[2]) //convertir l'année en entier
				if err != nil {
					log.Printf("Erreur lors de la conversion de l'année du premier album: %v", err)
					http.Error(w, "Erreur lors de la conversion de l'année du premier album", http.StatusInternalServerError)
					return
				}
				if firstAlbumYearInt >= startYear && firstAlbumYearInt <= endYear { //si l'année du premier album est comprise entre les années de début et de fin
					filteredData = append(filteredData, artist) // Ajouter l'artiste à la liste des artistes filtrés
				}
			}
			data = filteredData // Remplacer les données par les données filtrées
			log.Printf("Données filtrées: %v", filteredData)
		}
	}

	// Trier les artistes si le paramètre de tri est présent
	if sortOrder == "asc" {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Name < data[j].Name
		})
	} else if sortOrder == "desc" {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Name > data[j].Name
		})
	}

	//filteredArtists := Mod.SearchBar(query, data)
	optionsSearchBar := Mod.SearchOptions(query, data)

	var no_results bool
	if len(data) == 0 {
		no_results = true
		log.Printf("Aucun résultat trouvé")
	} else {
		log.Printf("Des résultats ont été trouvés")
	}

	// Charger le template HTML
	tmpl, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		log.Printf("Erreur lors du chargement du template: %v", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	log.Printf("Valeur de no_results: %v", no_results)

	pageData := PageData{
		Query:            displayQuery,
		Data:             data,
		OptionsSearchBar: optionsSearchBar,
		No_results:       no_results,
	}

	// Rendre le template avec les données
	err = tmpl.Execute(w, pageData)

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
		return
	}
}
