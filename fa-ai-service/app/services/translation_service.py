"""
Translation Service

Provides local language translation using LibreTranslate (self-hosted).
Supports Swahili, and falls back to English for unsupported languages.

API: https://github.com/LibreTranslate/LibreTranslate
"""

import httpx
from typing import Optional, List, Dict
from dataclasses import dataclass

from app.core.config import get_settings

settings = get_settings()


@dataclass
class TranslationResult:
    """Translation result."""
    original_text: str
    translated_text: str
    source_language: str
    target_language: str
    detected_language: Optional[str] = None


class TranslationService:
    """Translation service using LibreTranslate API."""
    
    # Supported languages for East African farming context
    SUPPORTED_LANGUAGES = {
        "en": "English",
        "sw": "Swahili (Kiswahili)",
        "fr": "French",  # For Rwanda/DRC farmers
        "ar": "Arabic",
        "pt": "Portuguese",
    }
    
    # Common farming terms translations (fallback cache)
    FARMING_TERMS_SW = {
        "disease": "ugonjwa",
        "treatment": "matibabu",
        "pesticide": "dawa ya kuulia wadudu",
        "fertilizer": "mbolea",
        "watering": "kumwagilia",
        "harvest": "mavuno",
        "planting": "kupanda",
        "healthy": "afya njema",
        "infected": "imeambukizwa",
        "spray": "kunyunyizia",
        "maize": "mahindi",
        "cassava": "muhogo",
        "tomato": "nyanya",
        "beans": "maharagwe",
        "banana": "ndizi",
    }
    
    def __init__(self):
        self.base_url = settings.libretranslate_url
        self.api_key = settings.libretranslate_api_key
        self.client = httpx.AsyncClient(timeout=10.0)
        self.enabled = bool(self.base_url)
    
    async def translate(
        self,
        text: str,
        target_language: str = "sw",
        source_language: str = "en",
    ) -> TranslationResult:
        """
        Translate text to target language.
        
        Args:
            text: Text to translate
            target_language: Target language code (e.g., 'sw' for Swahili)
            source_language: Source language code (default: 'en')
        
        Returns:
            TranslationResult with original and translated text
        """
        if not self.enabled:
            # LibreTranslate not configured, return original
            return TranslationResult(
                original_text=text,
                translated_text=text,
                source_language=source_language,
                target_language=target_language,
            )
        
        try:
            payload = {
                "q": text,
                "source": source_language,
                "target": target_language,
                "format": "text",
            }
            
            if self.api_key:
                payload["api_key"] = self.api_key
            
            response = await self.client.post(
                f"{self.base_url}/translate",
                json=payload,
            )
            response.raise_for_status()
            data = response.json()
            
            return TranslationResult(
                original_text=text,
                translated_text=data.get("translatedText", text),
                source_language=source_language,
                target_language=target_language,
                detected_language=data.get("detectedLanguage", {}).get("language"),
            )
            
        except Exception as e:
            print(f"Translation error: {e}")
            # Fallback to original text
            return TranslationResult(
                original_text=text,
                translated_text=text,
                source_language=source_language,
                target_language=target_language,
            )
    
    async def detect_language(self, text: str) -> Optional[str]:
        """Detect the language of text."""
        if not self.enabled:
            return "en"
        
        try:
            response = await self.client.post(
                f"{self.base_url}/detect",
                json={"q": text},
            )
            response.raise_for_status()
            data = response.json()
            
            if data and len(data) > 0:
                return data[0].get("language")
            return None
            
        except Exception:
            return None
    
    async def get_supported_languages(self) -> List[Dict]:
        """Get list of supported languages from LibreTranslate."""
        if not self.enabled:
            return [
                {"code": k, "name": v}
                for k, v in self.SUPPORTED_LANGUAGES.items()
            ]
        
        try:
            response = await self.client.get(f"{self.base_url}/languages")
            response.raise_for_status()
            return response.json()
        except Exception:
            return []
    
    def translate_farming_term(self, term: str, target_language: str = "sw") -> str:
        """Quick translation for common farming terms (cached)."""
        if target_language == "sw":
            term_lower = term.lower()
            if term_lower in self.FARMING_TERMS_SW:
                return self.FARMING_TERMS_SW[term_lower]
        return term
    
    async def translate_recommendation(
        self,
        recommendation: str,
        target_language: str = "sw",
    ) -> str:
        """
        Translate a farming recommendation with context preservation.
        
        Keeps technical terms and translates the rest.
        """
        if target_language == "en":
            return recommendation
        
        result = await self.translate(
            recommendation,
            target_language=target_language,
            source_language="en",
        )
        return result.translated_text
    
    async def close(self):
        """Close the HTTP client."""
        await self.client.aclose()


# Singleton
_translation_service: Optional[TranslationService] = None


def get_translation_service() -> TranslationService:
    """Get translation service singleton."""
    global _translation_service
    if _translation_service is None:
        _translation_service = TranslationService()
    return _translation_service
