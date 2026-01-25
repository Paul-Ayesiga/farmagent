# FA-AI-Service

AI-powered disease detection and agricultural assistant service.

## Overview

This service provides:

- **Disease Detection** - Image classification using HuggingFace pre-trained models
- **Treatment Recommendations** - AI-generated treatment plans via Ollama
- **Chat Assistant** - Conversational AI for farming questions

## Tech Stack

- **Framework**: FastAPI
- **ML**: HuggingFace Transformers, PyTorch
- **LLM**: Ollama (llama3.2:3b)
- **Language**: Python 3.11+

## Quick Start

### Local Development

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Run service
uvicorn app.main:app --reload --port 8005
```

### With Docker

```bash
docker build -t fa-ai-service .
docker run -p 8005:8005 fa-ai-service
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/ai/analyze` | POST | Analyze crop image |
| `/ai/analyze/base64` | POST | Analyze base64 image |
| `/ai/analyze/status` | GET | Model status |
| `/ai/analyze/diseases` | GET | List detectable diseases |
| `/ai/recommend` | POST | Treatment recommendations |
| `/ai/chat` | POST | Chat with AI |
| `/ai/chat/suggestions` | GET | Suggested questions |

## Environment Variables

```env
APP_ENV=development
APP_PORT=8005

# Ollama (for chat/recommendations)
OLLAMA_HOST=http://localhost:11434
OLLAMA_MODEL=llama3.2:3b
```

## Ollama Setup

For full AI chat functionality:

```bash
# Install Ollama
brew install ollama  # macOS

# Pull model
ollama pull llama3.2:3b

# Start server
ollama serve
```

## API Documentation

Interactive docs available at: <http://localhost:8005/docs>

## Testing

```bash
# Check model status
curl http://localhost:8005/ai/analyze/status

# Analyze image
curl -X POST http://localhost:8005/ai/analyze \
  -F "file=@plant.jpg"
```
