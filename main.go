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
				<select id="sort-select" name="sort" onchange="location = this.value;">
                    <option value="">Sort By</option>
                    <option value="?sort=asc">Sort Ascending</option>
                    <option value="?sort=desc">Sort Descending</option>
                </select>
				<form method="GET" action="/">
					<select id="filter-select" name="filter">
						<option value="">Filter By</option>
						<option value="?filter=CreatioDate">Creation Date</option>
					</select>
					<input type="number" name="start_year" placeholder="De ">
					<input type="number" name="end_year" placeholder="à ">
					<button type="Submit">GO</button>
				</form>
			</div>
			<div class="box">
    			<form name="search">
        			<input type="text" class="input" name="txt" onmouseout="this.value = ''; this.blur();">
    			</form>
    			<i class="image.png"></i>
			</div>
			<div class="container">
				{{range .}}
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

func handler(w http.ResponseWriter, r *http.Request) {
	// Appeler la fonction GetArtist depuis GroupieTracker.go
	data, err := Mod.GetData()
	if err != nil {
		log.Printf("Erreur lors de la récupération des données: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		return
	}

	// Trier les artistes si le paramètre de tri est présent
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

	// Filtrer les artistes par date de création si les paramètres sont présents
	startYearStr := r.URL.Query().Get("start_year")
	endYearStr := r.URL.Query().Get("end_year")
	if startYearStr != "" && endYearStr != "" {
		startYear, err := strconv.Atoi(startYearStr)
		if err != nil {
			log.Printf("Erreur lors de la conversion de l'année de début: %v", err)
			http.Error(w, "Erreur lors de la conversion de l'année de début", http.StatusInternalServerError)
			return
		}
		endYear, err := strconv.Atoi(endYearStr)
		if err != nil {
			log.Printf("Erreur lors de la conversion de l'année de fin: %v", err)
			http.Error(w, "Erreur lors de la conversion de l'année de fin", http.StatusInternalServerError)
			return
		}

		var filteredData []Mod.Artist
		for _, artist := range data {
			creationYear := artist.CreationDate
			if creationYear >= startYear && creationYear <= endYear {
				filteredData = append(filteredData, artist)
			}
		}
		data = filteredData
	}

	// Charger le template HTML
	tmpl, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		log.Printf("Erreur lors du chargement du template: %v", err)
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	// Rendre le template avec les données
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		return
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
	http.HandleFunc("/", handler)
	fmt.Println("Starting server at port 8080")
	go openBrowser("http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
