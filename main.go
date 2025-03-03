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
	Query   string
	Artists []Mod.Artist
}

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
				<div class="dropdown">
                <button class="dropbtn">Sorts</button>
                    <div class="dropdown-content">
                         <a href="?sort=asc">Sort Name Ascending</a>
                          <a href="?sort=desc">Sort Name Descending</a>
                     </div>
                </div>
			</div>
			<div class="box">
    			<form method="GET" action="/">
        			<input type="text" class="input" name="query" value="{{.Query}}" onmouseout="this.value = ''; this.blur();">
    			</form>
    			<i class="image.png"></i>
			</div>
			<div class="container">
				{{range .Artists}}
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
                			<center><img src="{{.Image}}" alt="Image" width="300" height="300"></center>
						</div>
						<div class="invisbox">
							<p>Members: {{range .Members}}{{.}}, {{end}}</p>
							<p>Creation Date: {{.CreationDate}}</p>
							<p>First Album: {{.FirstAlbum}}</p>
						</div>
						<div class="invisbox">
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
        		</div>
        		{{end}}
			</div>
		</div>
	</body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
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

	pageData := PageData{
		Query:   query,
		Artists: filteredArtists,
	}

	t, err := template.New("webpage").Parse(tmpl)
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
