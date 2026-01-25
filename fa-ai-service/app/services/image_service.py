"""
Image Processing Service

Handles image preprocessing for disease classification.
"""

import io
import base64
from typing import Tuple, Optional

import numpy as np
from PIL import Image

from app.core.config import get_settings


settings = get_settings()


class ImageService:
    """Service for image preprocessing and validation."""
    
    def __init__(self, input_size: int = 224):
        self.input_size = input_size
        self.supported_formats = settings.supported_formats.split(",")
        self.max_size_bytes = settings.max_image_size_mb * 1024 * 1024
    
    def validate_image(self, image_data: bytes) -> Tuple[bool, Optional[str]]:
        """
        Validate image data.
        Returns (is_valid, error_message)
        """
        # Check size
        if len(image_data) > self.max_size_bytes:
            return False, f"Image exceeds maximum size of {settings.max_image_size_mb}MB"
        
        # Check format
        try:
            img = Image.open(io.BytesIO(image_data))
            fmt = img.format.lower() if img.format else ""
            if fmt not in self.supported_formats:
                return False, f"Unsupported format: {fmt}. Supported: {self.supported_formats}"
        except Exception as e:
            return False, f"Invalid image: {str(e)}"
        
        return True, None
    
    def preprocess_for_model(self, image_data: bytes) -> np.ndarray:
        """
        Preprocess image for the disease classification model.
        
        Returns:
            Numpy array of shape (1, input_size, input_size, 3)
        """
        # Open image
        img = Image.open(io.BytesIO(image_data))
        
        # Convert to RGB if needed
        if img.mode != "RGB":
            img = img.convert("RGB")
        
        # Resize to model input size
        img = img.resize((self.input_size, self.input_size), Image.Resampling.LANCZOS)
        
        # Convert to numpy array
        img_array = np.array(img, dtype=np.float32)
        
        # Normalize to [0, 1]
        img_array = img_array / 255.0
        
        # Add batch dimension
        img_array = np.expand_dims(img_array, axis=0)
        
        return img_array
    
    def decode_base64(self, base64_string: str) -> bytes:
        """Decode base64 image string to bytes."""
        # Remove data URL prefix if present
        if "," in base64_string:
            base64_string = base64_string.split(",")[1]
        
        return base64.b64decode(base64_string)
    
    def encode_base64(self, image_data: bytes) -> str:
        """Encode image bytes to base64 string."""
        return base64.b64encode(image_data).decode("utf-8")
    
    def get_image_info(self, image_data: bytes) -> dict:
        """Get image metadata."""
        img = Image.open(io.BytesIO(image_data))
        return {
            "format": img.format,
            "mode": img.mode,
            "width": img.width,
            "height": img.height,
            "size_bytes": len(image_data),
        }
    
    def compress_image(
        self, 
        image_data: bytes, 
        max_dimension: int = 800,
        quality: int = 85
    ) -> bytes:
        """
        Compress image for storage.
        """
        img = Image.open(io.BytesIO(image_data))
        
        # Convert to RGB if needed
        if img.mode != "RGB":
            img = img.convert("RGB")
        
        # Resize if larger than max dimension
        if max(img.width, img.height) > max_dimension:
            img.thumbnail((max_dimension, max_dimension), Image.Resampling.LANCZOS)
        
        # Save to bytes
        buffer = io.BytesIO()
        img.save(buffer, format="JPEG", quality=quality, optimize=True)
        
        return buffer.getvalue()


# Singleton instance
_image_service: Optional[ImageService] = None


def get_image_service() -> ImageService:
    """Get image service singleton."""
    global _image_service
    if _image_service is None:
        _image_service = ImageService(settings.model_input_size)
    return _image_service
