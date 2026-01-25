"""
AI Analysis API Endpoints

Handles crop disease detection and image analysis.
"""

from typing import Optional
from fastapi import APIRouter, UploadFile, File, Form, HTTPException
from pydantic import BaseModel

from app.services.image_service import get_image_service
from app.services.disease_detector import get_disease_detector


router = APIRouter(prefix="/analyze", tags=["Analysis"])


class AnalysisResponse(BaseModel):
    """Response from image analysis."""
    disease_class: str
    disease_name: str
    crop_type: str
    confidence: float
    is_healthy: bool
    severity: str
    health_score: int
    top_predictions: list
    message: str


class Base64ImageRequest(BaseModel):
    """Request with base64 encoded image."""
    image: str  # Base64 encoded image
    crop_type: Optional[str] = None


@router.get("/status")
async def model_status():
    """
    Get the status of the disease detection model.
    
    Returns whether the model is loaded and ready for predictions.
    """
    detector = get_disease_detector()
    status = detector.get_model_status()
    
    if not status["loaded"]:
        return {
            "status": "unavailable",
            "message": status.get("error") or "Model not loaded. Install dependencies: pip install transformers torch",
            "ready": False,
        }
    
    return {
        "status": "ready",
        "message": f"Model loaded with {status['num_classes']} disease classes",
        "ready": True,
        "model_type": status["model_type"],
        "num_classes": status["num_classes"],
    }


@router.post("", response_model=AnalysisResponse)
async def analyze_image(
    file: UploadFile = File(...),
    crop_type: Optional[str] = Form(None),
):
    """
    Analyze a crop image for disease detection.
    
    Upload an image of a plant leaf to detect diseases.
    
    - **file**: Image file (JPEG, PNG)
    - **crop_type**: Optional hint about crop type (maize, cassava, tomato, etc.)
    
    Returns disease classification with confidence score and health assessment.
    """
    image_service = get_image_service()
    detector = get_disease_detector()
    
    # Check if model is ready
    if not detector.is_model_loaded():
        raise HTTPException(
            status_code=503,
            detail={
                "error": "Disease detection model not available",
                "message": detector.load_error or "Please install dependencies: pip install transformers torch",
                "hint": "Check /ai/analyze/status for more details",
            }
        )
    
    # Read image
    image_data = await file.read()
    
    # Validate
    is_valid, error = image_service.validate_image(image_data)
    if not is_valid:
        raise HTTPException(status_code=400, detail=error)
    
    # Detect disease
    try:
        result = detector.predict(image_data)
    except RuntimeError as e:
        raise HTTPException(status_code=503, detail=str(e))
    
    # Generate message
    if result.is_healthy:
        message = f"Your {result.crop_type} appears healthy! Health score: {result.health_score}/100"
    else:
        message = f"Detected {result.disease_name} on {result.crop_type} with {result.confidence:.0%} confidence. Severity: {result.severity}"
    
    return AnalysisResponse(
        disease_class=result.disease_class,
        disease_name=result.disease_name,
        crop_type=result.crop_type,
        confidence=result.confidence,
        is_healthy=result.is_healthy,
        severity=result.severity,
        health_score=result.health_score,
        top_predictions=result.top_predictions,
        message=message,
    )


@router.post("/base64", response_model=AnalysisResponse)
async def analyze_base64_image(request: Base64ImageRequest):
    """
    Analyze a base64-encoded crop image for disease detection.
    
    Alternative to file upload for mobile apps.
    
    - **image**: Base64 encoded image string
    - **crop_type**: Optional hint about crop type
    """
    image_service = get_image_service()
    detector = get_disease_detector()
    
    # Check if model is ready
    if not detector.is_model_loaded():
        raise HTTPException(
            status_code=503,
            detail={
                "error": "Disease detection model not available",
                "message": detector.load_error or "Please install dependencies: pip install transformers torch",
            }
        )
    
    # Decode base64
    try:
        image_data = image_service.decode_base64(request.image)
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Invalid base64 image: {e}")
    
    # Validate
    is_valid, error = image_service.validate_image(image_data)
    if not is_valid:
        raise HTTPException(status_code=400, detail=error)
    
    # Detect disease
    try:
        result = detector.predict(image_data)
    except RuntimeError as e:
        raise HTTPException(status_code=503, detail=str(e))
    
    # Generate message
    if result.is_healthy:
        message = f"Your {result.crop_type} appears healthy! Health score: {result.health_score}/100"
    else:
        message = f"Detected {result.disease_name} on {result.crop_type} with {result.confidence:.0%} confidence. Severity: {result.severity}"
    
    return AnalysisResponse(
        disease_class=result.disease_class,
        disease_name=result.disease_name,
        crop_type=result.crop_type,
        confidence=result.confidence,
        is_healthy=result.is_healthy,
        severity=result.severity,
        health_score=result.health_score,
        top_predictions=result.top_predictions,
        message=message,
    )


@router.get("/diseases")
async def list_diseases():
    """
    Get list of all detectable diseases.
    
    Returns the full list of crop diseases that can be detected.
    """
    from app.core.diseases import UGANDA_CROP_CLASSES, PLANTVILLAGE_CLASSES
    
    detector = get_disease_detector()
    
    return {
        "model_loaded": detector.is_model_loaded(),
        "model_classes": detector.classes if detector.is_model_loaded() else [],
        "uganda_crops": UGANDA_CROP_CLASSES,
        "plantvillage_classes": PLANTVILLAGE_CLASSES,
        "total_classes": len(detector.classes) if detector.is_model_loaded() else len(PLANTVILLAGE_CLASSES),
    }

