import React from "react";

function Gallery({ images, onImageClick }) {
    return (
        <div className="image-gallery">
            {images.length > 0 ? (
                images.map((post, index) => {
                    console.log("Type of URL:", typeof post, "Value:", post);
                    return (
                        <img
                            key={index}
                            src={post.ContentURL}
                            alt={`Image ${index}`}
                            className="img"
                            onClick={() => onImageClick(post, `Description for image ${index}`)}
                        />
                    );
                })
            ) : (
                <p>No images found</p>
            )}
        </div>
    );
}
export default Gallery;
