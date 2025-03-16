document.addEventListener("DOMContentLoaded", function () {
    fetch('/get-checked-options')
        .then(response => response.json())
        .then(data => {
            if (data.checkedOptions) {
                data.checkedOptions.forEach(option => {
                    const checkbox = document.querySelector(`input[name="filter"][value="${option}"]`);
                    if (checkbox) {
                        checkbox.checked = true;
                    }
                });
            }
        })
        .catch(err => console.error('Error fetching checked options:', err));
  });
  
  function updateCheckboxState(checkbox) {
    const option = checkbox.value;
    const isChecked = checkbox.checked;
  
    
    fetch('/update-checked-options', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ option, isChecked })
    })
  }
  
  document.addEventListener("DOMContentLoaded", function () {
    var optionsList = document.getElementById('optionsList');
    var searchQuery = document.getElementById('searchQuery');
  
    optionsList.addEventListener('input', function () {
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
  
    optionsList.addEventListener('change', function () {
        cleanQuery(this.value, searchQuery);
        this.form.submit();
    });
  
    optionsList.addEventListener('blur', function () {
        cleanQuery(this.value, searchQuery);
    });
  });
  
  function cleanQuery(value, searchQueryElement) {
    var cleanedQuery = value.replace(/ - Member$/, "");
    searchQueryElement.value = cleanedQuery;
  }
