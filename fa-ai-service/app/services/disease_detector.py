"""
Disease Detection Service

Handles crop disease classification using pre-trained models from HuggingFace.
"""

import os
from typing import Optional, Tuple, List, Dict, Any
from dataclasses import dataclass

import numpy as np

from app.core.config import get_settings
from app.core.diseases import (
    PLANTVILLAGE_CLASSES,
    DISEASE_INFO,
    get_crop_from_class,
    get_disease_from_class,
)
from app.services.image_service import get_image_service


settings = get_settings()


@dataclass
class DiseaseDetectionResult:
    """Result of disease detection."""
    disease_class: str
    disease_name: str
    crop_type: str
    confidence: float
    is_healthy: bool
    severity: str
    health_score: int
    top_predictions: List[Dict[str, Any]]


class DiseaseDetectorService:
    """Service for detecting crop diseases from images using pre-trained models."""
    
    def __init__(self):
        self.model = None
        self.processor = None
        self.classes = PLANTVILLAGE_CLASSES
        self.image_service = get_image_service()
        self.confidence_threshold = settings.confidence_threshold
        self.use_huggingface = True
        self._load_model()
    
    def _load_model(self):
        """Load pre-trained disease classification model from HuggingFace."""
        self.load_error = None
        
        try:
            from transformers import AutoImageProcessor, AutoModelForImageClassification
            import torch
            
            # Use a pre-trained plant disease model from HuggingFace
            model_name = "linkanjarad/mobilenet_v2_1.0_224-plant-disease-identification"
            
            print(f"📥 Loading pre-trained model: {model_name}")
            
            self.processor = AutoImageProcessor.from_pretrained(model_name)
            self.model = AutoModelForImageClassification.from_pretrained(model_name)
            self.model.eval()  # Set to evaluation mode
            
            # Get class labels from model config
            if hasattr(self.model.config, 'id2label'):
                self.classes = list(self.model.config.id2label.values())
            
            print(f"✅ Loaded model with {len(self.classes)} disease classes")
            self.use_huggingface = True
            
        except ImportError as e:
            self.load_error = "Required packages not installed. Run: pip install transformers torch"
            print(f"❌ {self.load_error}")
            self.use_huggingface = False
        except Exception as e:
            self.load_error = f"Failed to load model: {str(e)}"
            print(f"❌ {self.load_error}")
            self.use_huggingface = False
    
    def predict(self, image_data: bytes) -> DiseaseDetectionResult:
        """
        Predict disease from crop image.
        
        Args:
            image_data: Raw image bytes
            
        Returns:
            DiseaseDetectionResult with prediction details
            
        Raises:
            RuntimeError: If model is not loaded
        """
        if not self.use_huggingface or self.model is None:
            raise RuntimeError(
                "Disease detection model not available. "
                "Please ensure 'transformers' and 'torch' are installed: "
                "pip install transformers torch"
            )
        
        return self._predict_huggingface(image_data)
    
    def _predict_huggingface(self, image_data: bytes) -> DiseaseDetectionResult:
        """Predict using HuggingFace model."""
        import torch
        from PIL import Image
        import io
        
        # Open image
        img = Image.open(io.BytesIO(image_data))
        if img.mode != "RGB":
            img = img.convert("RGB")
        
        # Process image
        inputs = self.processor(images=img, return_tensors="pt")
        
        # Predict
        with torch.no_grad():
            outputs = self.model(**inputs)
            logits = outputs.logits
            probs = torch.nn.functional.softmax(logits, dim=-1)[0]
        
        # Get top predictions
        top_k = min(5, len(self.classes))
        top_probs, top_indices = torch.topk(probs, top_k)
        
        top_class_idx = top_indices[0].item()
        confidence = top_probs[0].item()
        disease_class = self.classes[top_class_idx]
        
        top_predictions = []
        for i in range(top_k):
            idx = top_indices[i].item()
            cls = self.classes[idx]
            top_predictions.append({
                "class": cls,
                "disease": get_disease_from_class(cls),
                "crop": get_crop_from_class(cls),
                "confidence": top_probs[i].item(),
            })
        
        # Extract info
        crop_type = get_crop_from_class(disease_class)
        disease_name = get_disease_from_class(disease_class)
        is_healthy = "healthy" in disease_class.lower()
        
        # Calculate severity and health score
        severity = self._calculate_severity(confidence, is_healthy)
        health_score = self._calculate_health_score(confidence, is_healthy)
        
        return DiseaseDetectionResult(
            disease_class=disease_class,
            disease_name=disease_name,
            crop_type=crop_type,
            confidence=confidence,
            is_healthy=is_healthy,
            severity=severity,
            health_score=health_score,
            top_predictions=top_predictions,
        )
    
    def is_model_loaded(self) -> bool:
        """Check if the model is loaded and ready."""
        return self.use_huggingface and self.model is not None
    
    def get_model_status(self) -> Dict[str, Any]:
        """Get status information about the model."""
        return {
            "loaded": self.is_model_loaded(),
            "model_type": "HuggingFace" if self.use_huggingface else "None",
            "num_classes": len(self.classes) if self.classes else 0,
            "error": self.load_error if hasattr(self, 'load_error') else None,
        }
    
    def _calculate_severity(self, confidence: float, is_healthy: bool) -> str:
        """Calculate disease severity based on confidence."""
        if is_healthy:
            return "none"
        
        if confidence >= 0.9:
            return "severe"
        elif confidence >= 0.7:
            return "moderate"
        elif confidence >= 0.5:
            return "mild"
        else:
            return "uncertain"
    
    def _calculate_health_score(self, confidence: float, is_healthy: bool) -> int:
        """Calculate health score (0-100)."""
        if is_healthy:
            return int(85 + confidence * 15)  # 85-100 for healthy
        else:
            # Lower score for diseased based on confidence
            return int(max(10, 70 - confidence * 60))  # 10-70 for diseased
    
    def get_disease_info(self, disease_class: str) -> Optional[Dict]:
        """Get detailed information about a disease."""
        # Try exact match
        if disease_class in DISEASE_INFO:
            return DISEASE_INFO[disease_class]
        
        # Try partial match
        disease_name = get_disease_from_class(disease_class)
        for key, info in DISEASE_INFO.items():
            if disease_name.lower() in key.lower():
                return info
        
        return None


# Singleton instance
_detector_service: Optional[DiseaseDetectorService] = None


def get_disease_detector() -> DiseaseDetectorService:
    """Get disease detector singleton."""
    global _detector_service
    if _detector_service is None:
        _detector_service = DiseaseDetectorService()
    return _detector_service

