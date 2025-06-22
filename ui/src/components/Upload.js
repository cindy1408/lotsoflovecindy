import React from "react";

function UploadSection({ setFile, handleUpload }) {
    const handleFileChange = (e) => {
        const selectedFile = e.target.files[0];
        if (selectedFile) {
            setFile(selectedFile);
            handleUpload(selectedFile);
        }
    };

    return (
        <div className="upload-section">
            <input type="file" onChange={handleFileChange} />
        </div>
    );
}


export default UploadSection;
