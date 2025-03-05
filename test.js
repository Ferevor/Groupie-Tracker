document.getElementById('searchInput').addEventListener('input', function (event) {
    console.log('Input event triggered');
    const searchTerm = event.target.value.toLowerCase();
    console.log('Search Term:', searchTerm);
    const listItems = document.querySelectorAll('#itemList li');
    console.log('List Items:', listItems);

    listItems.forEach(function (item) {
        const itemText = item.getAttribute('data-name').toLowerCase();
        console.log('Item Text:', itemText);

        if (itemText.includes(searchTerm)) {
            item.style.display = 'list-item';
        } else {
            item.style.display = 'none';
        }
    });
});
