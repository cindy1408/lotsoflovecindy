import React from "react";

function Gallery({ images, onImageClick }) {
    if (!Array.isArray(images)) {
        console.error("Expected 'images' to be an array, but got:", images);
        return <p>Image data is invalid</p>;
    }

    return (
        <div className="image-gallery">
            {images.length > 0 ? (
                images.map((post, index) => {
                    console.log("Type of post:", typeof post, "Value:", post);
                    return (
                        <img
                            key={index}
                            src={post.ContentURL}
                            alt={`Image ${index}`}
                            className="img"
                            onClick={() => onImageClick(post)}
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
