import React, { useState, useEffect } from "react";
import './App.css';

function App() {
    const [file, setFile] = useState(null); // Track selected file
    const [images, setImages] = useState([]); // Store image URLs
    const [selectedImage, setSelectedImage] = useState(null); // Track selected image for modal
    const [selectedDescription, setSelectedDescription] = useState(""); // Track selected image description

    // Fetch images when the component mounts
    useEffect(() => {
        fetch("http://localhost:8080/list-files")
            .then(async (response) => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const text = await response.text();
                if (!text) {
                    console.warn("Empty response body");
                    return [];
                }
                return JSON.parse(text);
            })
            .then((data) => {
                console.log("Fetched images:", data);
                setImages(data);
            })
            .catch((error) => console.error("Error fetching images:", error));
    }, []);

    // Handle file input change
    const handleFileChange = (e) => {
        setFile(e.target.files[0]);
    };

    // Function to handle file upload
    const uploadFile = () => {
        if (!file) {
            alert("Please select a file to upload");
            return;
        }

        const formData = new FormData();
        formData.append("file", file);

        fetch("http://localhost:8080/upload", {
            method: "POST",
            body: formData,
        })
            .then((response) => response.text())
            .then((data) => {
                alert(data); // Display the response from the server

                // Re-fetch images after upload
                fetch("http://localhost:8080/list-files")
                    .then((response) => response.json())
                    .then((data) => {
                        console.log("Updated images:", data);
                        setImages(data);
                    });
            })
            .catch((error) => {
                console.error("Error uploading file:", error);
                alert("Failed to upload file");
            });
    };

    // Handle image click to open modal with larger image and description
    const handleImageClick = (url, description) => {
        setSelectedImage(url);
        setSelectedDescription(description);
    };

    // Handle closing the modal
    const handleCloseModal = () => {
        setSelectedImage(null);
        setSelectedDescription("");
    };

    return (
        <div className="App">
            <header className="App-header">
                <div className="header">
                    <img className="profile" src="/data/img/profile.jpg" alt="Profile" />
                    <div>
                        <h2>Lots of Love, Cindy</h2>
                        <h3>bio: hello this is a test haha </h3>
                    </div>
                </div>

                <div className="upload-section">
                    <input type="file" id="fileInput" onChange={handleFileChange} />
                    <button onClick={uploadFile}>Upload</button> {/* Call uploadFile on button click */}
                </div>

                <div className="image-gallery">
                    {images.length > 0 ? (
                        images.map((url, index) => (
                            <img
                                key={index}
                                src={url.trim()}
                                alt={`Image ${index}`}
                                className="img"
                                onClick={() => handleImageClick(url, `Description for image ${index}`)} // Update description dynamically
                            />
                        ))
                    ) : (
                        <p>No images found</p>
                    )}
                </div>
            </header>

            {/* Modal to display selected image */}
            {selectedImage && (
                <div className="modal">
                    <div className="modal-content">
                        <span className="close" onClick={handleCloseModal}>&times;</span>
                        <img src={selectedImage} alt="Enlarged view" className="enlarged-img" />
                        <p>{selectedDescription}</p>
                    </div>
                </div>
            )}
        </div>
    );
}

export default App;
