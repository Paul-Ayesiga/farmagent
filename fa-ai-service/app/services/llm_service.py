"""
LLM Service

Handles integration with Ollama for treatment recommendations and chat.
"""

import json
from typing import Optional, Dict, List, Any
from dataclasses import dataclass

import httpx

from app.core.config import get_settings
from app.core.diseases import DISEASE_INFO


settings = get_settings()


@dataclass
class RecommendationResult:
    """Result of treatment recommendation."""
    disease: str
    crop: str
    treatments: Dict[str, List[str]]
    prevention: List[str]
    local_tips: List[str]
    weather_advice: Optional[str]
    raw_response: str


@dataclass
class ChatResponse:
    """Response from chat endpoint."""
    message: str
    suggestions: List[str]
    context_used: Optional[str]


class LLMService:
    """Service for LLM-based recommendations using Ollama."""
    
    def __init__(self):
        self.base_url = settings.ollama_host
        self.model = settings.ollama_model
        self.vision_model = settings.ollama_vision_model
        self.timeout = 60.0
    
    async def check_health(self) -> bool:
        """Check if Ollama is running."""
        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(f"{self.base_url}/api/tags", timeout=5.0)
                return response.status_code == 200
        except Exception:
            return False
    
    async def generate(
        self, 
        prompt: str, 
        system_prompt: Optional[str] = None,
        temperature: float = 0.7
    ) -> str:
        """
        Generate text response from LLM.
        
        Args:
            prompt: User prompt
            system_prompt: System instructions
            temperature: Creativity level (0-1)
            
        Returns:
            Generated text response
        """
        messages = []
        
        if system_prompt:
            messages.append({"role": "system", "content": system_prompt})
        
        messages.append({"role": "user", "content": prompt})
        
        try:
            async with httpx.AsyncClient() as client:
                response = await client.post(
                    f"{self.base_url}/api/chat",
                    json={
                        "model": self.model,
                        "messages": messages,
                        "stream": False,
                        "options": {"temperature": temperature},
                    },
                    timeout=self.timeout,
                )
                
                if response.status_code != 200:
                    raise Exception(f"Ollama error: {response.text}")
                
                data = response.json()
                return data.get("message", {}).get("content", "")
                
        except httpx.ConnectError:
            return self._fallback_response(prompt)
        except Exception as e:
            print(f"LLM Error: {e}")
            return self._fallback_response(prompt)
    
    async def get_treatment_recommendation(
        self,
        disease: str,
        crop: str,
        severity: str,
        location: Optional[str] = "Uganda",
        weather: Optional[Dict] = None,
    ) -> RecommendationResult:
        """
        Get treatment recommendations for a detected disease.
        """
        # First check if we have local info
        local_info = self._get_local_disease_info(disease, crop)
        
        # Build prompt for LLM
        weather_context = ""
        if weather:
            weather_context = f"""
Current weather: {weather.get('temperature', 'N/A')}°C, {weather.get('condition', 'N/A')}
Humidity: {weather.get('humidity', 'N/A')}%
"""
        
        system_prompt = """You are an expert agricultural advisor for East African farmers.
Provide practical, actionable advice that considers:
- Local availability of treatments
- Cost-effectiveness for smallholder farmers
- Both organic and chemical options
- Traditional farming knowledge
Keep responses concise and practical."""

        prompt = f"""A farmer in {location} has detected {disease} on their {crop} crop.
Severity: {severity}
{weather_context}

Provide:
1. Immediate treatment options (2-3 organic, 2-3 chemical)
2. Prevention measures for future
3. Local tips specific to East Africa

Format as JSON with keys: organic_treatments, chemical_treatments, prevention, local_tips"""

        response = await self.generate(prompt, system_prompt, temperature=0.5)
        
        # Parse response and merge with local info
        result = self._parse_recommendation(response, disease, crop, local_info)
        
        # Add weather advice if available
        weather_advice = None
        if weather:
            weather_advice = await self._get_weather_advice(disease, crop, weather)
        
        return RecommendationResult(
            disease=disease,
            crop=crop,
            treatments=result.get("treatments", {"organic": [], "chemical": []}),
            prevention=result.get("prevention", []),
            local_tips=result.get("local_tips", []),
            weather_advice=weather_advice,
            raw_response=response,
        )
    
    async def chat(
        self,
        message: str,
        context: Optional[str] = None,
        history: Optional[List[Dict]] = None,
    ) -> ChatResponse:
        """
        Have a conversation about farming topics with memory of previous messages.
        """
        system_prompt = """You are FarmAgent AI, a helpful agricultural assistant for East African farmers.
You specialize in:
- Crop disease identification and treatment
- Best farming practices for Uganda
- Sustainable agriculture techniques
- Market and pricing advice

Be concise, practical, and friendly. Use simple language.
If asked about something outside farming, politely redirect to agricultural topics."""

        # Build messages array with history
        messages = [{"role": "system", "content": system_prompt}]
        
        # Add conversation history
        if history:
            for msg in history:
                role = msg.get("role", "user")
                content = msg.get("content", "")
                if role in ["user", "assistant"] and content:
                    messages.append({"role": role, "content": content})
        
        # Add context and current message
        full_prompt = message
        if context:
            full_prompt = f"Context: {context}\n\nQuestion: {message}"
        
        messages.append({"role": "user", "content": full_prompt})
        
        # Call Ollama with full conversation
        try:
            async with httpx.AsyncClient() as client:
                response = await client.post(
                    f"{self.base_url}/api/chat",
                    json={
                        "model": self.model,
                        "messages": messages,
                        "stream": False,
                        "options": {"temperature": 0.7},
                    },
                    timeout=self.timeout,
                )
                
                if response.status_code != 200:
                    raise Exception(f"Ollama error: {response.text}")
                
                data = response.json()
                response_text = data.get("message", {}).get("content", "")
                
        except httpx.ConnectError:
            response_text = self._fallback_response(message)
        except Exception as e:
            print(f"LLM Error: {e}")
            response_text = self._fallback_response(message)
        
        # Generate follow-up suggestions
        suggestions = self._generate_suggestions(message, response_text)
        
        return ChatResponse(
            message=response_text,
            suggestions=suggestions,
            context_used=context,
        )
    
    async def chat_stream(
        self,
        message: str,
        context: Optional[str] = None,
        history: Optional[List[Dict]] = None,
    ):
        """
        Stream a conversation response token by token.
        
        Yields chunks of text as they are generated by Ollama.
        """
        system_prompt = """You are FarmAgent AI, a helpful agricultural assistant for East African farmers.
You specialize in:
- Crop disease identification and treatment
- Best farming practices for Uganda
- Sustainable agriculture techniques
- Market and pricing advice

Be concise, practical, and friendly. Use simple language.
If asked about something outside farming, politely redirect to agricultural topics."""

        messages = [{"role": "system", "content": system_prompt}]
        
        # Add history if provided
        if history:
            for msg in history:
                messages.append({"role": msg.get("role", "user"), "content": msg.get("content", "")})
        
        # Add context and current message
        full_prompt = message
        if context:
            full_prompt = f"Context: {context}\n\nQuestion: {message}"
        
        messages.append({"role": "user", "content": full_prompt})
        
        try:
            async with httpx.AsyncClient() as client:
                async with client.stream(
                    "POST",
                    f"{self.base_url}/api/chat",
                    json={
                        "model": self.model,
                        "messages": messages,
                        "stream": True,
                        "options": {"temperature": 0.7},
                    },
                    timeout=self.timeout,
                ) as response:
                    async for line in response.aiter_lines():
                        if line:
                            try:
                                import json
                                data = json.loads(line)
                                content = data.get("message", {}).get("content", "")
                                if content:
                                    yield content
                                # Check if done
                                if data.get("done", False):
                                    break
                            except json.JSONDecodeError:
                                continue
        except httpx.ConnectError:
            yield self._fallback_response(message)
        except Exception as e:
            print(f"LLM Stream Error: {e}")
            yield f"Sorry, I encountered an error: {str(e)}"
    
    def _get_local_disease_info(self, disease: str, crop: str) -> Optional[Dict]:
        """Get locally stored disease information."""
        for key, info in DISEASE_INFO.items():
            if disease.lower() in key.lower() or disease.lower() in info.get("disease", "").lower():
                return info
        return None
    
    def _parse_recommendation(
        self, 
        response: str, 
        disease: str, 
        crop: str,
        local_info: Optional[Dict]
    ) -> Dict:
        """Parse LLM response and merge with local info."""
        result = {
            "treatments": {"organic": [], "chemical": []},
            "prevention": [],
            "local_tips": [],
        }
        
        # Try to parse as JSON
        try:
            # Find JSON in response
            start = response.find("{")
            end = response.rfind("}") + 1
            if start >= 0 and end > start:
                parsed = json.loads(response[start:end])
                
                if "organic_treatments" in parsed:
                    result["treatments"]["organic"] = parsed["organic_treatments"]
                if "chemical_treatments" in parsed:
                    result["treatments"]["chemical"] = parsed["chemical_treatments"]
                if "prevention" in parsed:
                    result["prevention"] = parsed["prevention"]
                if "local_tips" in parsed:
                    result["local_tips"] = parsed["local_tips"]
        except json.JSONDecodeError:
            pass
        
        # Merge with local info
        if local_info:
            local_treatments = local_info.get("treatments", {})
            if local_treatments.get("organic"):
                result["treatments"]["organic"] = list(set(
                    result["treatments"]["organic"] + local_treatments["organic"]
                ))[:5]
            if local_treatments.get("chemical"):
                result["treatments"]["chemical"] = list(set(
                    result["treatments"]["chemical"] + local_treatments["chemical"]
                ))[:5]
            if local_info.get("prevention"):
                result["prevention"] = list(set(
                    result["prevention"] + local_info["prevention"]
                ))[:5]
        
        return result
    
    async def _get_weather_advice(
        self, 
        disease: str, 
        crop: str, 
        weather: Dict
    ) -> str:
        """Get weather-specific advice."""
        prompt = f"""Given:
- Disease: {disease}
- Crop: {crop}
- Temperature: {weather.get('temperature')}°C
- Humidity: {weather.get('humidity')}%
- Condition: {weather.get('condition')}

Give ONE brief sentence of weather-specific advice for treating this disease."""

        return await self.generate(prompt, temperature=0.3)
    
    def _generate_suggestions(self, message: str, response: str) -> List[str]:
        """Generate follow-up question suggestions."""
        suggestions = [
            "What organic treatments are available?",
            "How can I prevent this in the future?",
            "What's the best time to apply treatment?",
        ]
        
        # Context-aware suggestions
        if "maize" in message.lower() or "corn" in message.lower():
            suggestions.append("How do I protect maize from Fall Armyworm?")
        if "cassava" in message.lower():
            suggestions.append("Which cassava varieties are disease-resistant?")
        if "tomato" in message.lower():
            suggestions.append("How do I manage late blight in tomatoes?")
        
        return suggestions[:4]
    
    def _fallback_response(self, prompt: str) -> str:
        """Fallback response when Ollama is not available."""
        return """I apologize, but I'm currently unable to connect to the AI service.

Here are some general recommendations:
1. For disease treatment, consult your local agricultural extension officer
2. Check the National Agricultural Research Organisation (NARO) website
3. Visit your nearest agro-input shop for treatment options

Please try again later when the service is available."""


# Singleton instance
_llm_service: Optional[LLMService] = None


def get_llm_service() -> LLMService:
    """Get LLM service singleton."""
    global _llm_service
    if _llm_service is None:
        _llm_service = LLMService()
    return _llm_service
