"""
Chat API Endpoints

Handles conversational AI for farming assistance.
"""

import json
from typing import Optional, List
from fastapi import APIRouter
from pydantic import BaseModel

from app.services.llm_service import get_llm_service


router = APIRouter(prefix="/chat", tags=["Chat"])


class ChatMessage(BaseModel):
    """A single chat message."""
    role: str  # 'user' or 'assistant'
    content: str


class ChatRequest(BaseModel):
    """Chat request payload."""
    message: str
    context: Optional[str] = None
    history: Optional[List[ChatMessage]] = None


class ChatResponse(BaseModel):
    """Chat response payload."""
    message: str
    suggestions: List[str]


@router.post("", response_model=ChatResponse)
async def chat(request: ChatRequest):
    """
    Chat with FarmAgent AI assistant.
    
    Ask questions about:
    - Crop diseases and treatments
    - Farming best practices
    - Weather-based advice
    - Market information
    
    - **message**: Your question or message
    - **context**: Optional context (e.g., recent disease detection result)
    - **history**: Optional previous messages for context continuity
    
    Returns AI response with follow-up suggestions.
    """
    llm = get_llm_service()
    
    result = await llm.chat(
        message=request.message,
        context=request.context,
        history=[m.model_dump() for m in request.history] if request.history else None,
    )
    
    return ChatResponse(
        message=result.message,
        suggestions=result.suggestions,
    )


@router.post("/stream")
async def chat_stream(request: ChatRequest):
    """
    Stream chat response using Server-Sent Events (SSE).
    
    Returns responses token-by-token in real-time.
    Each chunk is sent as a JSON object with a 'content' field.
    The final chunk includes 'done': true and 'suggestions' array.
    
    - **message**: Your question or message
    - **context**: Optional context
    - **history**: Optional previous messages
    """
    from fastapi.responses import StreamingResponse
    
    async def generate():
        llm = get_llm_service()
        full_response = ""
        
        async for chunk in llm.chat_stream(
            message=request.message,
            context=request.context,
            history=[m.model_dump() for m in request.history] if request.history else None,
        ):
            full_response += chunk
            yield f"data: {json.dumps({'content': chunk, 'done': False})}\n\n"
        
        # Send final message with suggestions
        suggestions = [
            "What organic treatments are available?",
            "How can I prevent this in the future?",
            "What's the best time to apply treatment?",
        ]
        yield f"data: {json.dumps({'content': '', 'done': True, 'suggestions': suggestions})}\n\n"
    
    return StreamingResponse(
        generate(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no",
        }
    )


@router.get("/suggestions")
async def get_suggestions():
    """
    Get suggested questions to ask the AI.
    
    Returns common questions farmers might want to ask.
    """
    return {
        "suggestions": [
            "How do I identify Fall Armyworm in maize?",
            "What are the best organic pesticides for tomatoes?",
            "When is the best time to plant cassava in Uganda?",
            "How do I prevent disease spread in my crops?",
            "What are the current market prices for maize?",
            "How can I improve soil fertility naturally?",
            "Which crop varieties are drought-resistant?",
            "How do I manage pests during the rainy season?",
        ],
        "categories": {
            "disease": [
                "What disease is affecting my crop?",
                "How do I treat plant diseases naturally?",
            ],
            "prevention": [
                "How can I prevent crop diseases?",
                "What are integrated pest management techniques?",
            ],
            "cultivation": [
                "When should I plant?",
                "How much water do my crops need?",
            ],
            "market": [
                "What are current crop prices?",
                "Where can I sell my produce?",
            ],
        },
    }


@router.post("/quick")
async def quick_chat(question: str):
    """
    Quick endpoint for simple questions.
    
    Example: POST /chat/quick?question=How do I treat late blight?
    """
    llm = get_llm_service()
    
    result = await llm.chat(message=question)
    
    return {
        "answer": result.message,
        "suggestions": result.suggestions[:3],
    }
