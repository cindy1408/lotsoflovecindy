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
    const [selectedImage, setSelectedImage] = useState(null); // Image url shown in modal

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

            console.log("Raw fetched text:", text);
            console.log("Parsed JSON data:", data);

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
        formData.append("filename", file.name);
        formData.append("contentType", file.type);

        console.log(file.name)

        try {
            // Step 1: Get signed upload URL from server
            const response = await fetch("http://localhost:8080/upload", {
                method: "POST",
                body: formData,
            });

            const { signedUrl, publicUrl } = await response.json();

            // Step 2: Upload directly to Google Cloud Storage
            const uploadResponse = await fetch(signedUrl, {
                method: "PUT",
                headers: {
                    "Content-Type": file.type,
                },
                body: file,
            });

            if (uploadResponse.ok) {
                alert("Upload successful!");
                console.log("File accessible at:", publicUrl);
                fetchImages();
            } else {
                alert("Upload failed.");
            }
        } catch (error) {
            console.error("Upload failed:", error);
            alert("Failed to upload file");
        }
    };


    const handleDescriptionUpdate = async (newDescription) => {
        console.log("selected image: ", selectedImage);

        const formData = new FormData();
        formData.append("id", selectedImage.ID);
        formData.append("url_path", selectedImage.ContentURL);
        formData.append("description", newDescription);

        try {
            const response = await fetch("http://localhost:8080/update-description", {
                method: "POST",
                body: formData,
            });

            await response.text();

            // Update the description of the selected image in the images array
            setImages((prevImages) =>
                prevImages.map((image) =>
                    image.ID === selectedImage.ID
                        ? { ...image, Description: newDescription }
                        : image
                )
            );

            // Also update the selectedImage directly to reflect changes in the modal
            setSelectedImage((prevImage) => ({
                ...prevImage,
                Description: newDescription,
            }));
        } catch (error) {
            console.error("Update failed:", error);
            alert("Failed to update description");
        }
    };

    // === Render ===
    return (
        <div className="App">
            <header className="App-header">
                <Header />
                <UploadSection setFile={setFile} handleUpload={handleUpload} />
                <Gallery
                    images={images}
                    onImageClick={(image) => {
                        setSelectedImage(image);
                    }}
                />
            </header>

            {selectedImage && (
                <Modal
                    image={selectedImage.ContentURL}
                    description={selectedImage.Description}
                    onClose={() => {
                        setSelectedImage(null);
                    }}
                    updatedDescription={handleDescriptionUpdate}
                />
            )}
        </div>
    );
}

export default App;
