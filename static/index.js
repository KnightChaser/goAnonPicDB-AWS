// Fetch the list of uploaded images and repeatedly show them on the page.
document.addEventListener("DOMContentLoaded", function() {
    // Fetch the list of uploaded images
    fetch("/images")
        .then(response => response.json())
        .then(data => {
            const imageList = document.getElementById("imageList");

            // Create a Bootstrap row container
            let rowDiv;

            data.images.forEach((image, index) => {
                // Create a new row for every three images
                if (index % 3 === 0) {
                    rowDiv = document.createElement("div");
                    rowDiv.classList.add("row");
                    imageList.appendChild(rowDiv);
                }

                // Create the Bootstrap card structure
                const colDiv = document.createElement("div");
                colDiv.classList.add("col-md-4", "mb-3");

                const cardDiv = document.createElement("div");
                cardDiv.classList.add("card");

                // Create an image element
                const imgElement = document.createElement("img");
                imgElement.src = image.objectAccessURL; // Set the image source from the pre-signed URL
                imgElement.classList.add("card-img-top");
                imgElement.alt = image.key; // Set alt text

                // Create the card body
                const cardBodyDiv = document.createElement("div");
                cardBodyDiv.classList.add("card-body");

                // Display the key (filename)
                const keyElement = document.createElement("p");
                keyElement.classList.add("card-text");
                keyElement.innerHTML = `<b>Name:</b> <code>${image.key}</code>`;
                cardBodyDiv.appendChild(keyElement);

                // Display the uploader
                const uploaderElement = document.createElement("p");
                uploaderElement.classList.add("card-text");
                uploaderElement.innerHTML = `<b>Uploader:</b> ${image.uploader}`;
                cardBodyDiv.appendChild(uploaderElement);

                // Display the size
                const sizeElement = document.createElement("p");
                sizeElement.classList.add("card-text");
                sizeElement.innerHTML = `<b>Size:</b> ${image.size} bytes`;
                cardBodyDiv.appendChild(sizeElement);

                // Display the date
                const lastModifiedDate = document.createElement("p");
                lastModifiedDate.classList.add("card-text");
                lastModifiedDate.innerHTML = `<b>Last modified:</b> ${image.lastModified}`;
                cardBodyDiv.appendChild(lastModifiedDate);

                // Append elements to form the card structure
                cardDiv.appendChild(imgElement);
                cardDiv.appendChild(cardBodyDiv);
                colDiv.appendChild(cardDiv);
                
                // Append the column to the current row
                rowDiv.appendChild(colDiv);
            });
        })
        .catch(error => console.error("Error fetching image list:", error));
});
