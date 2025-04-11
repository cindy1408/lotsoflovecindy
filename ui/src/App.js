import React, { useState, useEffect } from "react";
import './App.css';
import UploadSection from "./components/Upload";
import Header from "./components/Header";
import Gallery from "./components/gallery";
import Modal from "./components/modal";

function App() {
    // State
    const [file, setFile] = useState(null);    // Selected file to upload
    const [images, setImages] = useState([]); // Lists of images objects
    // const [posts, setPosts] = useState([]); // List of posts
    const [selectedImage, setSelectedImage] = useState(null); // Image url shown in modal
    const [selectedDescription, setSelectedDescription] = useState(""); // Image description in modal

    // Fetch images from server
    const fetchImages = async () => {
        try {
            console.log("Fetching images...")
            const response = await fetch("http://localhost:8080/list-files");
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

            const text = await response.text();
            if (!text) return setImages([]);

            const data = JSON.parse(text);
            setImages(data);
        } catch (error) {
            console.error("Error fetching images:", error);
        }
    };

    useEffect(() => {
        fetchImages();
    }, []);

    // === Upload File Handler ===
    const handleUpload = async () => {
        if (!file) {
            alert("Please select a file to upload");
            return;
        }

        const formData = new FormData();
        formData.append("file", file);

        try {
            const response = await fetch("http://localhost:8080/upload", {
                method: "POST",
                body: formData,
            });

            const message = await response.text();
            alert(message);

            // Refresh image gallery after upload
            fetchImages();
        } catch (error) {
            console.error("Upload failed:", error);
            alert("Failed to upload file");
        }
    };

    const handleDescriptionUpdate = async (newDescription) => {
        console.log("selected image: ", selectedImage)

        const formData = new FormData();
        formData.append("postUrl", selectedImage);
        formData.append("updatedDescription", newDescription);

        try {
            const response = await fetch("http://localhost:8080/update-description", {
                method: "POST",
                body: formData,
            });

            const message = await response.text();
            alert(message);

            // Refresh image gallery after upload
            fetchImages();
        } catch (error) {
            console.error("Upload failed:", error);
            alert("Failed to upload file");
        }

        setSelectedDescription(newDescription);
    };

    // === Render ===
    return (
        <div className="App">
            <header className="App-header">
                <Header />
                <UploadSection setFile={setFile} handleUpload={handleUpload} />
                <Gallery
                    images={images}
                    onImageClick={(url, desc) => {
                        setSelectedImage(url);
                        setSelectedDescription(desc);
                    }}
                />
            </header>

            <Modal
                image={selectedImage}
                description={selectedDescription}
                onClose={() => {
                    setSelectedImage(null);
                    setSelectedDescription("");
                }}
                updatedDescription={handleDescriptionUpdate}
            />
        </div>
    );
}

export default App;
