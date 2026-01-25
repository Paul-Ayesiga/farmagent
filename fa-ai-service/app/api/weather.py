"""
Weather API Endpoints

Provides weather information and farming advisories.
"""

from fastapi import APIRouter, HTTPException, Query
from pydantic import BaseModel
from typing import Optional

from app.services.weather_service import (
    get_weather_service,
    get_weather_for_region,
    UGANDA_LOCATIONS,
)


router = APIRouter(prefix="/weather", tags=["Weather"])


class WeatherResponse(BaseModel):
    """Weather response model."""
    latitude: float
    longitude: float
    temperature: float
    humidity: float
    precipitation: float
    precipitation_probability: float
    wind_speed: float
    weather_code: int
    weather_description: str
    is_rainy: bool
    forecast_rain_24h: bool
    forecast_rain_48h: bool
    farming_advisory: str
    spray_advisory: str


@router.get("", response_model=WeatherResponse)
async def get_weather(
    latitude: float = Query(..., description="Latitude", ge=-90, le=90),
    longitude: float = Query(..., description="Longitude", ge=-180, le=180),
):
    """
    Get weather for a specific location.
    
    Returns current weather conditions and farming advisories.
    
    - **latitude**: Location latitude (e.g., 0.3476 for Kampala)
    - **longitude**: Location longitude (e.g., 32.5825 for Kampala)
    """
    service = get_weather_service()
    weather = await service.get_weather(latitude, longitude)
    
    if not weather:
        raise HTTPException(
            status_code=503,
            detail="Weather service temporarily unavailable"
        )
    
    return WeatherResponse(
        latitude=weather.latitude,
        longitude=weather.longitude,
        temperature=weather.temperature,
        humidity=weather.humidity,
        precipitation=weather.precipitation,
        precipitation_probability=weather.precipitation_probability,
        wind_speed=weather.wind_speed,
        weather_code=weather.weather_code,
        weather_description=weather.weather_description,
        is_rainy=weather.is_rainy,
        forecast_rain_24h=weather.forecast_rain_24h,
        forecast_rain_48h=weather.forecast_rain_48h,
        farming_advisory=weather.farming_advisory,
        spray_advisory=weather.spray_advisory,
    )


@router.get("/region/{region}", response_model=WeatherResponse)
async def get_weather_by_region(region: str):
    """
    Get weather for a named Uganda region.
    
    Available regions: kampala, mbarara, gulu, jinja, mbale, fort_portal, soroti, lira
    """
    weather = await get_weather_for_region(region)
    
    if not weather:
        raise HTTPException(
            status_code=503,
            detail="Weather service temporarily unavailable"
        )
    
    return WeatherResponse(
        latitude=weather.latitude,
        longitude=weather.longitude,
        temperature=weather.temperature,
        humidity=weather.humidity,
        precipitation=weather.precipitation,
        precipitation_probability=weather.precipitation_probability,
        wind_speed=weather.wind_speed,
        weather_code=weather.weather_code,
        weather_description=weather.weather_description,
        is_rainy=weather.is_rainy,
        forecast_rain_24h=weather.forecast_rain_24h,
        forecast_rain_48h=weather.forecast_rain_48h,
        farming_advisory=weather.farming_advisory,
        spray_advisory=weather.spray_advisory,
    )


@router.get("/regions")
async def list_regions():
    """
    List available Uganda regions with their coordinates.
    """
    return {
        "regions": [
            {"name": name, "latitude": lat, "longitude": lon}
            for name, (lat, lon) in UGANDA_LOCATIONS.items()
        ]
    }


@router.get("/spray-check")
async def check_spray_conditions(
    latitude: float = Query(..., ge=-90, le=90),
    longitude: float = Query(..., ge=-180, le=180),
):
    """
    Quick check: Is it safe to spray pesticides/fertilizers?
    
    Returns a simple yes/no with explanation.
    """
    service = get_weather_service()
    weather = await service.get_weather(latitude, longitude)
    
    if not weather:
        return {
            "safe_to_spray": None,
            "message": "Cannot check weather. Use your judgment based on local conditions.",
        }
    
    safe = not weather.is_rainy and not weather.forecast_rain_24h and weather.wind_speed < 15
    
    return {
        "safe_to_spray": safe,
        "message": weather.spray_advisory,
        "current_conditions": {
            "temperature": weather.temperature,
            "is_rainy": weather.is_rainy,
            "rain_expected_24h": weather.forecast_rain_24h,
            "wind_speed": weather.wind_speed,
        }
    }
