import React from "react";

function UploadSection({ setFile, handleUpload }) {
    return (
        <div className="upload-section">
            <input type="file" onChange={(e) => setFile(e.target.files[0])} />
            <button onClick={handleUpload}>Upload</button>
        </div>
    );
}

export default UploadSection;
