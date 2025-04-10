import React from "react";

function Gallery({ images, onImageClick }) {
    return (
        <div className="image-gallery">
            {images.length > 0 ? (
                images.map((url, index) => (
                    <img
                        key={index}
                        src={url.trim()}
                        alt={`Image ${index}`}
                        className="img"
                        onClick={() => onImageClick(url, `Description for image ${index}`)} // Update description dynamically
                    />
                ))
            ) : (
                <p>No images found</p>
            )}
        </div>
    );
}

export default Gallery;
