document.addEventListener("DOMContentLoaded", function() { 
    var optionsList = document.getElementById('optionsList'); 
    var searchQuery = document.getElementById('searchQuery');
    
    optionsList.addEventListener('input', function() { 
        var query = this.value; 
        if (query.length > 0) { 
            fetch('/search?q=' + query) 
                .then(response => response.text()) 
                .then(data => { 
                    var dataList = document.getElementById('suggestionsquery'); 
                    dataList.innerHTML = data; 
                }); 
        } 
    }); 

    optionsList.addEventListener('change', function() { 
        var cleanedQuery = this.value.replace(/ - Member$/, ""); 
        searchQuery.value = cleanedQuery;
        this.form.submit(); 
    }); 

    optionsList.addEventListener('blur', function() { 
        var cleanedQuery = this.value.replace(/ - Member$/, ""); 
        searchQuery.value = cleanedQuery;
    });
});