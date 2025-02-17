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
	<div>
		<div class="header">
			<h1>Groupie Tracker</h1>
		</div>
		<div class="box">
  		 	 <form name="search">
        		<input type="text" class="input" name="txt" onmouseout="this.value = ''; this.blur();">
    		</form>
    		<i class="image.png"></i>
		</div>
		<div class="container">
			{{range .}}
			<div class="card">
				<div class="image">
					<img src="{{.Image}}" alt="Image" width="200" height="200">
				</div>
				<div class="content">
					<h2>{{.Name}}</h2>
				</div>
			</div>
			{{end}}
		</div>
	</div>
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
