// Fetch the list of uploaded images and repeatedly show it on the page.
document.addEventListener("DOMContentLoaded", function() {
    // Fetch the list of uploaded images
    fetch("/images")
        .then(response => response.json())
        .then(data => {
            const imageList = document.getElementById("imageList");
            data.images.forEach(image => {
                const listItem = document.createElement("li");

                // Create a div to display image details
                const detailsDiv = document.createElement("div");

                // Display the key (filename)
                const keyElement = document.createElement("p");
                keyElement.textContent = `Key: ${image.key}`;
                detailsDiv.appendChild(keyElement);

                // Display the uploader
                const uploaderElement = document.createElement("p");
                uploaderElement.textContent = `Uploader: ${image.uploader}`;
                detailsDiv.appendChild(uploaderElement);

                // Display the size
                const sizeElement = document.createElement("p");
                sizeElement.textContent = `Size: ${image.size} bytes`;
                detailsDiv.appendChild(sizeElement);

                // Display the date
                const lastModifiedDate = document.createElement("p");
                lastModifiedDate.textContent = `Last modified: ${image.lastModified}`;
                detailsDiv.appendChild(lastModifiedDate);

                // Append the details div to the list item
                listItem.appendChild(detailsDiv);

                // Append the list item to the image list
                imageList.appendChild(listItem);
            });
        })
        .catch(error => console.error("Error fetching image list:", error));
});
