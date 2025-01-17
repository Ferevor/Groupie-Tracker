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
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Groupie Tracker</title>
</head>
<body>
    <h1>Groupie Tracker</h1>
    <div id="content">
        <pre>{{.}}</pre>
    </div>
</body>
</html>
`

func handler(w http.ResponseWriter, r *http.Request) {
	// Appeler la fonction GetArtist depuis GroupieTracker.go
	data, err := Mod.GetArtist()
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}

	// Charger le template HTML
	tmpl, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		log.Fatalf("Erreur lors du chargement du template: %v", err)
	}

	// Rendre le template avec les données
	err = tmpl.Execute(w, data)
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
	http.HandleFunc("/", handler)
	fmt.Println("Starting server at port 8080")
	go openBrowser("http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
