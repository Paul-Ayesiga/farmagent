"""
FarmAgent AI Service

FastAPI application for crop disease detection and agricultural AI assistance.
"""

from contextlib import asynccontextmanager
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.core.config import get_settings
from app.api import analyze, recommend, chat, weather


settings = get_settings()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan events."""
    # Startup
    print("🌱 FarmAgent AI Service starting...")
    
    # Preload models
    from app.services.disease_detector import get_disease_detector
    from app.services.llm_service import get_llm_service
    
    detector = get_disease_detector()
    llm = get_llm_service()
    
    # Check Ollama connection
    if await llm.check_health():
        print(f"✅ Connected to Ollama at {settings.ollama_host}")
        print(f"   Using model: {settings.ollama_model}")
    else:
        print(f"⚠️ Ollama not available at {settings.ollama_host}")
        print("   Recommendations will use fallback responses")
    
    print("🚀 AI Service ready!")
    
    yield
    
    # Shutdown
    print("👋 AI Service shutting down...")


# Create FastAPI app
app = FastAPI(
    title="FarmAgent AI Service",
    description="""
AI-powered agricultural assistant for East African farmers.

## Features

* **Disease Detection** - Upload crop images for disease identification
* **Treatment Recommendations** - Get organic and chemical treatment options
* **Chat Assistant** - Ask questions about farming practices
* **Weather Integration** - Location-based farming advisories

## Endpoints

* `/ai/analyze` - Image analysis for disease detection
* `/ai/recommend` - Treatment recommendations
* `/ai/chat` - Conversational AI assistant
* `/ai/weather` - Weather and spray advisories
    """,
    version="1.0.0",
    lifespan=lifespan,
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# Health check
@app.get("/health")
async def health_check():
    """Check service health."""
    from app.services.llm_service import get_llm_service
    
    llm = get_llm_service()
    ollama_ok = await llm.check_health()
    
    return {
        "status": "healthy",
        "service": "fa-ai-service",
        "ollama": "connected" if ollama_ok else "disconnected",
        "model": settings.ollama_model,
    }


# Mount API routers
app.include_router(analyze.router, prefix="/api/v1/ai")
app.include_router(recommend.router, prefix="/api/v1/ai")
app.include_router(chat.router, prefix="/api/v1/ai")
app.include_router(weather.router, prefix="/api/v1/ai")


# Direct mount for gateway compatibility
app.include_router(analyze.router, prefix="/ai")
app.include_router(recommend.router, prefix="/ai")
app.include_router(chat.router, prefix="/ai")
app.include_router(weather.router, prefix="/ai")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=settings.app_port,
        reload=settings.app_env == "development",
    )
