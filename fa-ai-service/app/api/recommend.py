"""
Recommendation API Endpoints

Handles treatment recommendations from LLM.
"""

from typing import Optional, List
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.services.llm_service import get_llm_service


router = APIRouter(prefix="/recommend", tags=["Recommendations"])


class RecommendationRequest(BaseModel):
    """Request for treatment recommendations."""
    disease: str
    crop: str
    severity: str = "moderate"
    location: str = "Uganda"
    latitude: Optional[float] = None
    longitude: Optional[float] = None


class TreatmentResponse(BaseModel):
    """Treatment recommendation response."""
    disease: str
    crop: str
    treatments: dict
    prevention: List[str]
    local_tips: List[str]
    weather_advice: Optional[str]


@router.post("", response_model=TreatmentResponse)
async def get_recommendations(request: RecommendationRequest):
    """
    Get treatment recommendations for a detected disease.
    
    - **disease**: Name of the disease (e.g., "Late Blight", "Fall Armyworm")
    - **crop**: Type of crop (e.g., "Tomato", "Maize")
    - **severity**: Severity level (mild, moderate, severe)
    - **location**: Location for localized advice (default: Uganda)
    - **latitude/longitude**: For weather-specific advice
    
    Returns organic and chemical treatments, prevention tips, and local advice.
    """
    llm = get_llm_service()
    
    # Get weather if coordinates provided
    weather = None
    if request.latitude and request.longitude:
        weather = await _get_weather(request.latitude, request.longitude)
    
    # Get recommendations
    result = await llm.get_treatment_recommendation(
        disease=request.disease,
        crop=request.crop,
        severity=request.severity,
        location=request.location,
        weather=weather,
    )
    
    return TreatmentResponse(
        disease=result.disease,
        crop=result.crop,
        treatments=result.treatments,
        prevention=result.prevention,
        local_tips=result.local_tips,
        weather_advice=result.weather_advice,
    )


@router.get("/quick/{disease}/{crop}")
async def quick_recommend(disease: str, crop: str):
    """
    Quick recommendation endpoint for simple cases.
    
    Example: /recommend/quick/late_blight/tomato
    """
    llm = get_llm_service()
    
    result = await llm.get_treatment_recommendation(
        disease=disease.replace("_", " ").title(),
        crop=crop.title(),
        severity="moderate",
    )
    
    return {
        "disease": result.disease,
        "crop": result.crop,
        "organic_treatments": result.treatments.get("organic", [])[:3],
        "chemical_treatments": result.treatments.get("chemical", [])[:2],
        "prevention": result.prevention[:3],
    }


async def _get_weather(lat: float, lon: float) -> Optional[dict]:
    """Fetch weather from Open-Meteo (free, no API key)."""
    import httpx
    
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(
                "https://api.open-meteo.com/v1/forecast",
                params={
                    "latitude": lat,
                    "longitude": lon,
                    "current_weather": True,
                    "hourly": "relativehumidity_2m",
                },
                timeout=5.0,
            )
            
            if response.status_code == 200:
                data = response.json()
                current = data.get("current_weather", {})
                hourly = data.get("hourly", {})
                
                humidity = hourly.get("relativehumidity_2m", [None])[0]
                
                return {
                    "temperature": current.get("temperature"),
                    "condition": _weather_code_to_text(current.get("weathercode", 0)),
                    "humidity": humidity,
                    "windspeed": current.get("windspeed"),
                }
    except Exception as e:
        print(f"Weather fetch failed: {e}")
    
    return None


def _weather_code_to_text(code: int) -> str:
    """Convert WMO weather code to text."""
    codes = {
        0: "Clear sky",
        1: "Mainly clear",
        2: "Partly cloudy",
        3: "Overcast",
        45: "Foggy",
        48: "Depositing rime fog",
        51: "Light drizzle",
        53: "Moderate drizzle",
        61: "Slight rain",
        63: "Moderate rain",
        65: "Heavy rain",
        80: "Slight rain showers",
        81: "Moderate rain showers",
        95: "Thunderstorm",
    }
    return codes.get(code, "Unknown")
