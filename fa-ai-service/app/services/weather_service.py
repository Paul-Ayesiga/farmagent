"""
Weather Service

Provides weather-aware farming recommendations using Open-Meteo API.
Free, no API key required.

API: https://open-meteo.com
"""

import httpx
from typing import Optional, Dict, Any, List
from datetime import datetime, timedelta
from dataclasses import dataclass
from functools import lru_cache

from app.core.config import get_settings

settings = get_settings()


@dataclass
class WeatherData:
    """Current and forecast weather data."""
    latitude: float
    longitude: float
    temperature: float  # Celsius
    humidity: float  # Percentage
    precipitation: float  # mm
    precipitation_probability: float  # Percentage
    wind_speed: float  # km/h
    weather_code: int
    weather_description: str
    is_rainy: bool
    forecast_rain_24h: bool
    forecast_rain_48h: bool
    farming_advisory: str
    spray_advisory: str
    raw_data: Dict


# WMO Weather interpretation codes
WEATHER_CODES = {
    0: "Clear sky",
    1: "Mainly clear",
    2: "Partly cloudy",
    3: "Overcast",
    45: "Foggy",
    48: "Depositing rime fog",
    51: "Light drizzle",
    53: "Moderate drizzle",
    55: "Dense drizzle",
    61: "Slight rain",
    63: "Moderate rain",
    65: "Heavy rain",
    71: "Slight snow",
    73: "Moderate snow",
    75: "Heavy snow",
    80: "Slight rain showers",
    81: "Moderate rain showers",
    82: "Violent rain showers",
    95: "Thunderstorm",
    96: "Thunderstorm with slight hail",
    99: "Thunderstorm with heavy hail",
}


class WeatherService:
    """Weather service using Open-Meteo free API."""
    
    BASE_URL = "https://api.open-meteo.com/v1/forecast"
    
    def __init__(self):
        self.client = httpx.AsyncClient(timeout=10.0)
    
    async def get_weather(
        self,
        latitude: float,
        longitude: float,
    ) -> Optional[WeatherData]:
        """
        Get current weather and 48-hour forecast for a location.
        
        Args:
            latitude: Location latitude (e.g., 0.3476 for Kampala)
            longitude: Location longitude (e.g., 32.5825 for Kampala)
        
        Returns:
            WeatherData with current conditions and farming advisories
        """
        params = {
            "latitude": latitude,
            "longitude": longitude,
            "current": [
                "temperature_2m",
                "relative_humidity_2m",
                "precipitation",
                "weather_code",
                "wind_speed_10m",
            ],
            "hourly": [
                "precipitation_probability",
                "precipitation",
            ],
            "forecast_days": 3,
            "timezone": "Africa/Kampala",
        }
        
        try:
            response = await self.client.get(self.BASE_URL, params=params)
            response.raise_for_status()
            data = response.json()
            
            return self._parse_weather_data(data, latitude, longitude)
            
        except Exception as e:
            print(f"Weather API error: {e}")
            return None
    
    def _parse_weather_data(
        self,
        data: Dict,
        latitude: float,
        longitude: float,
    ) -> WeatherData:
        """Parse Open-Meteo response into WeatherData."""
        current = data.get("current", {})
        hourly = data.get("hourly", {})
        
        # Current conditions
        temperature = current.get("temperature_2m", 25)
        humidity = current.get("relative_humidity_2m", 70)
        precipitation = current.get("precipitation", 0)
        weather_code = current.get("weather_code", 0)
        wind_speed = current.get("wind_speed_10m", 5)
        
        # Check forecast for rain
        precip_probs = hourly.get("precipitation_probability", [])
        precip_amounts = hourly.get("precipitation", [])
        
        # Next 24 hours (indices 0-23)
        rain_24h = any(p > 50 for p in precip_probs[:24]) if precip_probs else False
        # Next 48 hours (indices 0-47)
        rain_48h = any(p > 50 for p in precip_probs[:48]) if precip_probs else False
        
        # Current precipitation probability (average of next 6 hours)
        precip_prob = sum(precip_probs[:6]) / 6 if precip_probs else 0
        
        # Is it currently rainy?
        is_rainy = weather_code >= 51 or precipitation > 0.1
        
        # Weather description
        weather_desc = WEATHER_CODES.get(weather_code, "Unknown")
        
        # Generate farming advisories
        farming_advisory = self._generate_farming_advisory(
            temperature, humidity, is_rainy, rain_24h, wind_speed
        )
        spray_advisory = self._generate_spray_advisory(
            is_rainy, rain_24h, wind_speed, precip_prob
        )
        
        return WeatherData(
            latitude=latitude,
            longitude=longitude,
            temperature=temperature,
            humidity=humidity,
            precipitation=precipitation,
            precipitation_probability=precip_prob,
            wind_speed=wind_speed,
            weather_code=weather_code,
            weather_description=weather_desc,
            is_rainy=is_rainy,
            forecast_rain_24h=rain_24h,
            forecast_rain_48h=rain_48h,
            farming_advisory=farming_advisory,
            spray_advisory=spray_advisory,
            raw_data=data,
        )
    
    def _generate_farming_advisory(
        self,
        temp: float,
        humidity: float,
        is_rainy: bool,
        rain_24h: bool,
        wind: float,
    ) -> str:
        """Generate farming advisory based on weather."""
        advisories = []
        
        if is_rainy:
            advisories.append("🌧️ Current rain - avoid fieldwork and spraying.")
        elif rain_24h:
            advisories.append("⛈️ Rain expected in 24 hours - complete spraying today.")
        else:
            advisories.append("☀️ Good conditions for fieldwork.")
        
        if temp > 32:
            advisories.append("🌡️ High temperature - water crops in early morning.")
        elif temp < 15:
            advisories.append("❄️ Cool conditions - protect sensitive seedlings.")
        
        if humidity > 85:
            advisories.append("💧 High humidity - watch for fungal diseases.")
        
        if wind > 20:
            advisories.append("💨 Windy - avoid pesticide spraying.")
        
        return " ".join(advisories)
    
    def _generate_spray_advisory(
        self,
        is_rainy: bool,
        rain_24h: bool,
        wind: float,
        precip_prob: float,
    ) -> str:
        """Generate spray/treatment advisory."""
        if is_rainy:
            return "❌ DO NOT SPRAY - Currently raining. Wait for dry conditions."
        
        if precip_prob > 70:
            return "❌ DO NOT SPRAY - High chance of rain. Pesticides will wash off."
        
        if wind > 15:
            return "⚠️ CAUTION - High winds may cause spray drift. Wait for calmer conditions."
        
        if rain_24h:
            return "⚠️ SPRAY EARLY - Rain expected within 24 hours. Apply in the morning."
        
        return "✅ GOOD TO SPRAY - Favorable conditions for pesticide/fertilizer application."
    
    async def close(self):
        """Close the HTTP client."""
        await self.client.aclose()


# Singleton
_weather_service: Optional[WeatherService] = None


def get_weather_service() -> WeatherService:
    """Get weather service singleton."""
    global _weather_service
    if _weather_service is None:
        _weather_service = WeatherService()
    return _weather_service


# Default locations for Uganda regions
UGANDA_LOCATIONS = {
    "kampala": (0.3476, 32.5825),
    "mbarara": (-0.6167, 30.6500),
    "gulu": (2.7809, 32.2995),
    "jinja": (0.4244, 33.2041),
    "mbale": (1.0750, 34.1756),
    "fort_portal": (0.6710, 30.2750),
    "soroti": (1.7150, 33.6111),
    "lira": (2.2499, 32.5339),
}


async def get_weather_for_region(region: str) -> Optional[WeatherData]:
    """Get weather for a named Uganda region."""
    region_lower = region.lower().replace(" ", "_")
    
    if region_lower not in UGANDA_LOCATIONS:
        # Default to Kampala
        lat, lon = UGANDA_LOCATIONS["kampala"]
    else:
        lat, lon = UGANDA_LOCATIONS[region_lower]
    
    service = get_weather_service()
    return await service.get_weather(lat, lon)
