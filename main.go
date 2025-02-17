package main

import (
	"fmt"
	"groupie/Mod"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

const tmpl = `
<!DOCTYPE html>
<html>
<head>
    <title>Artist Info</title>
    <link rel="stylesheet" type="text/css" href="/Styles/style.css">
</head>
<body>
    {{range .}}
    	<h1>{{.Name}}</h1>
    	<img src="{{.Image}}" alt="{{.Name}}">
    	<p>Members: {{range .Members}}{{.}}, {{end}}</p>
    	<p>Creation Date: {{.CreationDate}}</p>
    	<p>First Album: {{.FirstAlbum}}</p>
    	<p>Relations:</p>
    	{{range $key, $value := .DatesLocations.DatesLocations}}
    	    <p>{{$key}}</p>
    	    <ul>
    	    	{{range $value}}
            		<li>{{.}}</li>
    	    	{{end}}
        	</ul>
    	{{end}}
    {{end}}
</body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
	// Call GetData() from Mod package
	data, err := Mod.GetData()
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}

	// Load HTML template
	t, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		log.Fatalf("Erreur lors du chargement du template: %v", err)
	}

	// Render the template with data
	err = t.Execute(w, data)
	if err != nil {
		log.Fatalf("Erreur lors du rendu du template: %v", err)
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
