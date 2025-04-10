import React from "react";

function Modal({ image, description, onClose }) {
    if (!image) return null;

    return (
        <div className="modal">
            <div className="modal-content">
                <span className="close" onClick={onClose}>&times;</span>
                <img src={image} alt="Enlarged view" className="enlarged-img" />
                <p>{description}</p>
            </div>
        </div>
    );
}

export default Modal;