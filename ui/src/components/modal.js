import React, { useState } from "react";

function ImageModal({ image, description, onClose, updatedDescription }) {
    const [editMode, setEditMode] = useState(false);
    const [editedDescription, setEditedDescription] = useState(description);

    if (!image) return null;

    const handleBackgroundClick = (event) => {
        if (event.target.className === "modal") {
            onClose();
        }
    };

    const handleBlur = () => {
        if (editedDescription !== description) {
            updatedDescription(editedDescription);
        }
        setEditMode(false);
    };

    return (
        <div className="modal" onClick={handleBackgroundClick}>
            <div className="modal-content">
                <span className="close" onClick={onClose}>&times;</span>
                <img src={image} alt="Enlarged view" className="enlarged-img" />

                {editMode ? (
                    <textarea
                        className="description-edit"
                        value={editedDescription}
                        onChange={(e) => setEditedDescription(e.target.value)}
                        onBlur={handleBlur}
                        autoFocus
                        rows={3}
                    />
                ) : (
                    <div
                        className="description-view"
                        onClick={() => setEditMode(true)}
                        style={{ cursor: "pointer" }}
                        title="Click to edit"
                    >
                        <p>{description || <em>(Click to add a description)</em>}</p>
                    </div>
                )}
            </div>
        </div>
    );
}

export default ImageModal;
