import React, { useState } from "react";
import { ReactComponent as DeleteIcon } from '../icons/delete.svg';

function ImageModal({ image, description, onClose, updatedDescription, dateCreated, onDelete}) {
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
            updatedDescription(editedDescription); // Update the parent component's description
        }
        setEditMode(false); // Exit edit mode
    };


    const handleDelete = async () => {
        try {
            const formData = new FormData();
            formData.append("url", image);

            const response = await fetch("http://localhost:8080/delete-post", {
                method: "POST",
                body: formData,
            });

            if (!response.ok) throw new Error("Delete failed");

            console.log("Delete successful, notifying parent for refresh");
            onDelete?.(image);

        } catch (error) {
            console.error("Error deleting post:", error);
            alert("Failed to delete post.");
        }
    };


    return (
        <div className="modal" onClick={handleBackgroundClick}>
            <div className="modal-content">
                <span className="close" onClick={onClose}>&times;</span>
                <img src={image} alt="Enlarged view" className="enlarged-img" />

                <div>
                    <div
                        className="delete-post"
                        onClick={handleDelete}
                        style={{ cursor: "pointer" }}
                    >
                        <DeleteIcon width={24} height={24} />
                    </div>
                    <p>
                        <em>Date Uploaded: </em>
                        {new Date(dateCreated).toLocaleDateString()}
                    </p>
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
        </div>
    );
}

export default ImageModal;
