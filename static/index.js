// Fetch the list of uploaded images and repeatedly show it on the page.
document.addEventListener("DOMContentLoaded", function() {
    // Fetch the list of uploaded images
    fetch("/images")
        .then(response => response.json())
        .then(data => {
            const imageList = document.getElementById("imageList");
            data.images.forEach(image => {
                const listItem = document.createElement("li");
                listItem.textContent = image;
                imageList.appendChild(listItem);
            });
        })
        .catch(error => console.error("Error fetching image list:", error));
});
