import React from "react";

function ImageModal({ image, description, onClose }) {
    if (!image) return null;

    const handleBackgroundClick = (event) => {
        if (event.target.className === "modal") {
            onClose();
        }
    };

    return (
        <div className="modal" onClick={handleBackgroundClick}>
            <div className="modal-content">
                <span className="close" onClick={onClose}>&times;</span>
                <img src={image} alt="Enlarged view" className="enlarged-img" />
                <p>{description}</p>
            </div>
        </div>
    );
}

export default ImageModal;
