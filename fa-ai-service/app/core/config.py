"""
FarmAgent AI Service Configuration
"""

from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""
    
    # App
    app_env: str = "development"
    app_port: int = 8005
    
    # Ollama
    ollama_host: str = "http://localhost:11434"
    ollama_model: str = "llama3.2:3b"
    ollama_vision_model: str = "llava:7b"
    
    # Model
    model_path: str = "./models/disease_classifier.h5"
    model_input_size: int = 224
    confidence_threshold: float = 0.7
    
    # Image Processing
    max_image_size_mb: int = 10
    supported_formats: str = "jpg,jpeg,png"
    
    # Weather API (Open-Meteo - free, no key needed)
    weather_api_url: str = "https://api.open-meteo.com/v1/forecast"
    
    # Translation (LibreTranslate - self-hosted)
    libretranslate_url: str = ""  # e.g., "http://localhost:5000"
    libretranslate_api_key: str = ""  # Optional
    
    # Redis
    redis_url: str = "redis://localhost:6379"
    
    class Config:
        env_file = ".env"
        case_sensitive = False


@lru_cache()
def get_settings() -> Settings:
    """Get cached settings instance."""
    return Settings()
