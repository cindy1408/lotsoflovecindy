import React, { useState, useEffect } from "react"; // Add useEffect
import './App.css';

function App() {
    const [file, setFile] = useState(null); // Track selected file
    const [images, setImages] = useState([]); // Store image URLs

    // Fetch images when the component mounts
    useEffect(() => {
        fetch("http://localhost:8080/list-files")
            .then((response) => response.json())
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

    return (
        <div className="App">
            <header className="App-header">
                <img className="profile" src="/data/img/profile.jpg" alt="Profile"/>
                <div>
                    <h2>Lots of Love, Cindy</h2>
                    <h3>bio: hello this is a test haha </h3>
                </div>
                <input type="file" id="fileInput" onChange={handleFileChange} />
                <button onClick={uploadFile}>Upload</button> {/* Call uploadFile on button click */}

                {/* Image Gallery */}
                <div className="image-gallery">
                    {images.length > 0 ? (
                        images.map((url, index) => (
                            <img key={index} src={url.trim()} alt={`Image ${index}`} className="img" />
                        ))
                    ) : (
                        <p>No images found</p>
                    )}
                </div>
            </header>
        </div>
    );
}

export default App;
