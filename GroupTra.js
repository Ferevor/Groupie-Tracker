document.addEventListener("DOMContentLoaded", function() { 
    var optionsList = document.getElementById('optionsList'); 
    optionsList.addEventListener('input', function() { 
        var query = this.value; 
        if (query.length > 0) { 
            fetch('/search?q=' + query) 
                .then(response => response.text()) 
                .then(data => { 
                    var dataList = document.getElementById('fontstyle'); 
                    dataList.innerHTML = data; 
                }); 
        } 
    }); 
    optionsList.addEventListener('change', function() { 
        this.form.submit(); 
    }); 
});
