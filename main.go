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

// Structure représentant les données nécessaires pour la page Web
type PageData struct {
	Query            string       // La requête de l'utilisateur à afficher
	Data             []Mod.Artist // Les données des artistes filtrées
	OptionsSearchBar []string     // Le tableau de suggestions pour la recherche
	CheckedOptions   []string     // Les options cochées parmi les filtres
	No_results       bool
}

// Variable globale pour stocker les options cochées par l'utilisateur
var checkedOptions = []string{}

func handler(w http.ResponseWriter, r *http.Request) {
	// Récupère les paramètres de la requête
	query := r.URL.Query().Get("query")
	displayQuery := r.URL.Query().Get("displayQuery") // Texte affiché pour la requête

	// Récupère les données des artistes via le package "Mod"
	data, err := Mod.GetData()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}

	start := r.URL.Query().Get("start") // Année de début pour la recherche
	end := r.URL.Query().Get("end")     // Année de fin pour la recherche

	if start != "" && end != "" { // Si les années de début et de fin sont renseignées
		startYear, err := strconv.Atoi(start) // Convertir les années en entiers
		if err != nil {
			log.Printf("Error converting start year to int: %v", err)
			http.Error(w, "Error converting start year to int", http.StatusBadRequest)
			return
		}
		endYear, err := strconv.Atoi(end) // Convertir les années en entiers
		if err != nil {
			log.Printf("Error converting end year to int: %v", err)
			http.Error(w, "Error converting end year to int", http.StatusBadRequest)
			return
		}

		var filteredArtists []Mod.Artist
		switch r.URL.Query().Get("filter") {
		case "CreationDate": //dans le cas où l'utilisateur veut filtrer par date de création
			for _, artist := range data { //on parcourt les artistes
				if artist.CreationDate >= startYear && artist.CreationDate <= endYear { //on vérifie si la date de création de l'artiste est comprise entre les dates de début et de fin
					filteredArtists = append(filteredArtists, artist) //si c'est le cas, on ajoute l'artiste à la liste des artistes filtrés
				}
			}
			data = filteredArtists //on met à jour les données avec les artistes filtrés
		case "FirstAlbum": //dans le cas où l'utilisateur veut filtrer par date du premier album
			for _, artist := range data { //on parcourt les artistes
				firstAlbumYear := strings.Split(artist.FirstAlbum, "-")   //on récupère l'année du premier album de l'artiste
				firstAlbumYearInt, err := strconv.Atoi(firstAlbumYear[2]) //on convertit l'année en entier
				if err != nil {
					log.Printf("Error converting first album year to int: %v", err)
					http.Error(w, "Error converting first album year to int", http.StatusBadRequest)
					return
				}
				if firstAlbumYearInt >= startYear && firstAlbumYearInt <= endYear { //on vérifie si l'année du premier album de l'artiste est comprise entre les dates de début et de fin
					filteredArtists = append(filteredArtists, artist) //si c'est le cas, on ajoute l'artiste à la liste des artistes filtrés
				}
			}
			data = filteredArtists //on met à jour les données avec les artistes filtrés
		}
	}
	sortOrder := r.URL.Query().Get("sort") // Récupère le paramètre de tri "sort" dans la chaîne de requête URL
	if sortOrder == "asc" {                // Si le tri est ascendant
		sort.Slice(data, func(i, j int) bool {
			return data[i].Name < data[j].Name
		})
	} else if sortOrder == "desc" { // Si le tri est descendant
		sort.Slice(data, func(i, j int) bool {
			return data[i].Name > data[j].Name
		})
	}

	// Filtrage des artistes selon les options cochées et la requête utilisateur
	filteredArtists := Mod.SearchBarCheckBox(checkedOptions, query, data)
	optionsSearchBar := Mod.SearchOptions(query, data) // Suggestions pour la barre de recherche

	var no_results bool            // Variable pour indiquer si aucun résultat n'a été trouvé
	if len(filteredArtists) == 0 { // Si la liste des artistes filtrés est vide
		no_results = true
	}

	// Structure des données à transmettre à la page
	pageData := PageData{
		Query:            displayQuery,
		Data:             filteredArtists,
		OptionsSearchBar: optionsSearchBar,
		CheckedOptions:   checkedOptions,
		No_results:       no_results,
	}

	// Chargement et exécution du template
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

// Fonction qui gère une requête pour récupérer les options cochées actuellement
func getCheckedOptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Indiquer que le contenu renvoyé est au format JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode la liste globale des options cochées (checkedOptions) en JSON et l'envoie dans la réponse
	json.NewEncoder(w).Encode(map[string][]string{"checkedOptions": checkedOptions})
}

func updateCheckedOptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Structure pour recevoir les données du client au format JSON
	var selection struct {
		Option    string `json:"option"`    // Option à ajouter ou retirer
		IsChecked bool   `json:"isChecked"` // Indique si l'option est cochée ou décochée
	}

	// Tenter de décoder les données du JSON dans la structure selection
	if err := json.NewDecoder(r.Body).Decode(&selection); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	if selection.IsChecked {
		// Si l'option est cochée, on vérifie si elle n'est pas déjà dans la liste checkedOptions
		if !Mod.Contains(checkedOptions, selection.Option) {
			// Si l'option n'y est pas alors l'ajouter
			checkedOptions = append(checkedOptions, selection.Option)
		}
	} else {
		// Sinon on la retire avec la fonction RemoveFromCheckedOptions du fichier GroupieTracker.go
		checkedOptions = Mod.RemoveFromCheckedOptions(checkedOptions, selection.Option)
	}

	// Renvoyer un message de OK pour dire que la mise à jour a été faite avec succès
	w.WriteHeader(http.StatusOK)
}

func searchOptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le paramètre de recherche "q" dans la chaîne de requête URL
	query := r.URL.Query().Get("q")

	// Appeler la fonction Mod.GetData() pour récupérer les données nécessaires à la recherche
	data, err := Mod.GetData()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	// Fait une recherche dans les données en fonction de la requête utilisateur avec la fonction Mod.SearchOptions()
	optionsSearchBar := Mod.SearchOptions(query, data)

	// On limite le nombre de résultats à un maximum de 5
	if len(optionsSearchBar) > 5 {
		optionsSearchBar = optionsSearchBar[:5]
	}

	// Indiquer que le contenu renvoyé est au format JSON
	w.Header().Set("Content-Type", "text/html")

	// Boucler sur les résultats de la recherche et créer une balise <option> pour chaque option
	for _, option := range optionsSearchBar {
		// Écrire chaque option sous forme d'élément HTML dans la réponse
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

	// Gère les requêtes vers le dossier "Scripts", de manière similaire au dossier "Styles".
	http.Handle("/Scripts/", http.StripPrefix("/Scripts/", http.FileServer(http.Dir("Scripts"))))

	http.Handle("/Images/", http.StripPrefix("/Images/", http.FileServer(http.Dir("Images"))))

	// Gère la page d'accueil avec la fonction Handler()
	http.HandleFunc("/", handler)

	// Gère les cases cochées par l'utilisateur avec /get-checked-options
	http.HandleFunc("/get-checked-options", getCheckedOptionsHandler)

	// Met à jour la liste des cases cochées par l'utilisateur avec /update-checked-options
	http.HandleFunc("/update-checked-options", updateCheckedOptionsHandler)

	// Gère les requêtes pour optenir les suggestions parmi les données disponibles et selon les critères donnés dans la requête avec /search
	http.HandleFunc("/search", searchOptionsHandler)

	// Ouvre automatiquement le navigateur web à l'adresse "http://localhost:8080"
	fmt.Println("Starting server at port 8080")
	// Démarre le serveur HTTP sur le port 8080
	go openBrowser("http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
