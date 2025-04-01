import React, { useState } from "react"; // state to handle file input
import './App.css';

function App() {
    const [file, setFile] = useState(null); // track selected file

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
            })
            .catch((error) => {
                console.error("Error uploading file:", error);
                alert("Failed to upload file");
            });
    };

    return (
    <div className="App">
      <header className="App-header">
          <img className="profile" src="/data/img/profile.jpg" alt="Image 10"/>
          <div>
              <h2>Lots of Love, Cindy</h2>
              <h3>bio: hello this is a test haha </h3>
          </div>
          <input type="file"
                 id="fileInput"
                 onChange={handleFileChange} // Update state when file is selected
          />
          <button onClick={uploadFile}>Upload</button> {/* Call uploadFile on button click */}
      </header>
    </div>
  );
}

export default App;



// 1. Set up state to handle file input in React using useState.
//
// 2. Handle the file input change event to store the file selected by the user.
//
// 3. Modify the uploadFile function to be called on button click.
//
// 4. Update the button's onClick event to call the function directly.