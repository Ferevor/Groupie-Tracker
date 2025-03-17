// Exécute le code une fois que le html est complètement chargé
document.addEventListener("DOMContentLoaded", function () {
    var optionsList = document.getElementById('optionsList'); // Récupère la barre de recherche
    var searchQuery = document.getElementById('searchQuery'); // Endroit où la barre de recherche nettoyée sera stockée
  
    optionsList.addEventListener('input', function () {
        var query = this.value; // Obtient la valeur tapée par l'utilisateur
        if (query.length > 0) {
            fetch('/search?q=' + query) // Effectue une requête GET avec le terme de recherche
                .then(response => response.text()) // Convertit la réponse en texte brut
                .then(data => {
                    var dataList = document.getElementById('suggestionsquery'); // Liste des suggestions
                    dataList.innerHTML = data; // Met à jour le contenu avec les suggestions reçues
                });
        }
    });

    // Traite la sélection d'une option dans la liste des suggestions de la barre de recherche
    optionsList.addEventListener('change', function () {
        cleanQuery(this.value, searchQuery); // Nettoie la valeur sélectionnée avant traitement avec la fonction cleanQuery()
        this.form.submit(); // Soumet la requête automatiquement après la sélection
    });

    // Envoie une requête GET pour récupérer les options cochées
    fetch('/get-checked-options')
        .then(response => response.json())  // Convertit la réponse en JSON 
        .then(data => {
            if (data.checkedOptions) {  // Vérifie si la réponse contient des options cochées
                data.checkedOptions.forEach(option => {  // Parcourt chaque option cochée pour mettre à jour l'état des cases à cocher dans l'interface
                    const checkbox = document.querySelector(`input[name="filter"][value="${option}"]`);
                    if (checkbox) {
                        checkbox.checked = true; // Coche automatiquement la case si elle existe dans le DOM
                    }
                });
            }
        })
        .catch(err => console.error('Error fetching checked options:', err)); // Affiche toutes erreur survenus lors de la récupération ou du traitement des données
  });
  
  // Fonction appelée lorsqu'une case est cochée ou décochée
  function updateCheckboxState(checkbox) {
    const option = checkbox.value; // Récupère la valeur de la case à cocher
    const isChecked = checkbox.checked; // Regarde si elle est cochée ou pas
    
    // Envoie une requête POST pour mettre à jour l'état de l'option côté serveur
    fetch('/update-checked-options', { // Utilise la méthode POST pour envoyer des données
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }, // Indique que le corps de la requête est en JSON
        body: JSON.stringify({ option, isChecked }) // Convertit les données en JSON
    })
  }

  // Fonction qui enlève "- Member" de la suggestion choisie par l'utilisateur s'il est présent
  function cleanQuery(value, searchQueryElement) {
    var cleanedQuery = value.replace(/ - Member$/, "");
    searchQueryElement.value = cleanedQuery;
  }


