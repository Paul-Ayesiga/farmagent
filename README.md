<div align="center">

# 🌱 FarmAgent

**AI-Powered Agricultural Assistant for East African Farmers**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.11+-3776AB?style=flat&logo=python)](https://python.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://docker.com)

[Features](#-features) • [Architecture](#-architecture) • [Quick Start](#-quick-start) • [API Docs](#-api-documentation) • [Development](#-development)

</div>

---

## 📋 Overview

FarmAgent is a comprehensive agricultural technology platform designed to help East African farmers improve crop yields through AI-powered disease detection, personalized treatment recommendations, and mobile money payment integration.

### 🎯 Features

- **🔬 AI Disease Detection** - Upload crop images to identify diseases using deep learning
- **💬 AI Chat Assistant** - Get farming advice from an intelligent chatbot powered by Ollama
- **💊 Treatment Recommendations** - Receive organic and chemical treatment suggestions
- **🌦️ Weather Integration** - Location-based weather advice for farming decisions
- **💳 Mobile Money Payments** - MTN MoMo integration for subscription payments
- **🔐 Secure Authentication** - JWT-based auth with email verification
- **📱 Mobile-First API** - RESTful APIs designed for mobile applications

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Client Applications                              │
│                    (Flutter Mobile App / Web Dashboard)                       │
└─────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           API Gateway (Port 8000)                            │
│                    JWT Validation • Rate Limiting • Routing                   │
└─────────────────────────────────────────────────────────────────────────────┘
           │              │              │              │              │
           ▼              ▼              ▼              ▼              ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
    │   Auth   │   │   Crop   │   │ Notifi-  │   │    AI    │   │ Payment  │
    │ Service  │   │ Service  │   │ cation   │   │ Service  │   │ Service  │
    │  :8001   │   │  :8002   │   │  :8003   │   │  :8005   │   │  :8004   │
    └────┬─────┘   └────┬─────┘   └────┬─────┘   └────┬─────┘   └────┬─────┘
         │              │              │              │              │
         ▼              ▼              ▼              ▼              ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
    │ Postgres │   │ Postgres │   │ MongoDB  │   │  Ollama  │   │   MTN    │
    │          │   │          │   │          │   │   LLM    │   │   MoMo   │
    └──────────┘   └──────────┘   └──────────┘   └──────────┘   └──────────┘
```

### Services Overview

| Service | Port | Tech Stack | Description |
|---------|------|------------|-------------|
| **API Gateway** | 8000 | Go, Chi | Request routing, JWT validation, rate limiting |
| **Auth Service** | 8001 | Go, Chi, PostgreSQL | User authentication, registration, password reset |
| **Crop Service** | 8002 | Go, Chi, PostgreSQL | Crop management, health records, treatments |
| **Notification** | 8003 | Go, Chi, MongoDB | Email, SMS, push notifications |
| **Payment** | 8004 | Go, Chi, PostgreSQL | MTN MoMo integration, subscriptions |
| **AI Service** | 8005 | Python, FastAPI | Disease detection, recommendations, chat |

---

## 🚀 Quick Start

### Prerequisites

Before you begin, ensure you have the following installed:

| Tool | Version | Installation |
|------|---------|--------------|
| **Docker Desktop** | 24.0+ | [Download](https://docker.com/products/docker-desktop) |
| **Docker Compose** | 2.20+ | Included with Docker Desktop |
| **Git** | 2.40+ | [Download](https://git-scm.com) |
| **Ollama** (optional) | Latest | [Download](https://ollama.ai) |

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/farmagent/farmagent.git
cd farmagent
```

### 2️⃣ Environment Setup

Create environment files for each service:

```bash
# Copy environment templates
cp fa-auth-service/.env.example fa-auth-service/.env
cp fa-notification-service/.env.example fa-notification-service/.env
cp fa-payment-service/.env.example fa-payment-service/.env
cp fa-ai-service/.env.example fa-ai-service/.env
```

#### Configure MTN MoMo (Payment Service)

1. Register at [momodeveloper.mtn.com](https://momodeveloper.mtn.com)
2. Subscribe to the **Collections** API
3. Get your credentials and update `fa-payment-service/.env`:

```env
MTN_SUBSCRIPTION_KEY=your-subscription-key
MTN_API_USER=your-api-user-uuid
MTN_API_KEY=your-api-key
```

### 3️⃣ Start All Services

```bash
# Start infrastructure and all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### 4️⃣ Verify Installation

Check that all services are healthy:

```bash
# Check service health
curl http://localhost:8000/health          # API Gateway
curl http://localhost:8001/health          # Auth Service
curl http://localhost:8002/health          # Crop Service
curl http://localhost:8003/health          # Notification Service
curl http://localhost:8004/health          # Payment Service
curl http://localhost:8005/health          # AI Service
```

Expected response:

```json
{"status": "healthy", "service": "fa-api-gateway"}
```

### 5️⃣ Test the API

```bash
# Register a new user
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "farmer@example.com",
    "password": "SecurePass123!",
    "full_name": "John Mukasa",
    "phone": "+256771234567"
  }'

# Login
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "farmer@example.com",
    "password": "SecurePass123!"
  }'
```

---

## 🤖 AI Service Setup

The AI Service requires additional setup for full functionality:

### Option A: With Ollama (Recommended)

For AI chat and recommendations:

```bash
# Install Ollama
brew install ollama  # macOS
# or download from https://ollama.ai

# Pull the model
ollama pull llama3.2:3b

# Start Ollama
ollama serve
```

The AI service will automatically connect to Ollama at `localhost:11434`.

### Option B: Without Ollama

The AI service works without Ollama for disease detection. Chat/recommendations will return helpful fallback messages.

### Test AI Endpoints

```bash
# Check AI model status
curl http://localhost:8005/ai/analyze/status

# Get chat suggestions
curl http://localhost:8005/ai/chat/suggestions

# Test disease detection (requires image file)
curl -X POST http://localhost:8005/ai/analyze \
  -F "file=@/path/to/plant-image.jpg"
```

---

## 📚 API Documentation

### Interactive Documentation

Once services are running, access Swagger/OpenAPI docs:

| Service | Documentation URL |
|---------|-------------------|
| AI Service | <http://localhost:8005/docs> |

### Postman Collection

Import the complete API collection:

```bash
# Located at:
docs/FarmAgent.postman_collection.json
```

### API Endpoints Summary

#### Authentication (`/api/v1/auth`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register new user |
| POST | `/login` | User login |
| POST | `/refresh` | Refresh access token |
| POST | `/forgot-password` | Request password reset |
| POST | `/reset-password` | Reset password with token |
| GET | `/users/me` | Get current user profile |

#### Crops (`/api/v1/crops`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/fields` | List user's fields |
| POST | `/fields` | Create new field |
| GET | `/crops` | List crops |
| POST | `/health-records` | Add health record |

#### AI (`/api/v1/ai`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/analyze` | Analyze crop image |
| POST | `/recommend` | Get treatment recommendations |
| POST | `/chat` | Chat with AI assistant |
| GET | `/analyze/status` | Check model status |

#### Payments (`/api/v1/payments`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/initiate` | Initiate MoMo payment |
| GET | `/{id}/status` | Check payment status |
| GET | `/history` | Payment history |

#### Subscriptions (`/api/v1/subscriptions`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/plans` | List available plans |
| GET | `/` | Current subscription |
| POST | `/subscribe` | Subscribe to plan |

---

## 🛠️ Development

### Project Structure

```
farmagent/
├── fa-api-gateway/          # API Gateway (Go)
├── fa-auth-service/         # Authentication Service (Go)
├── fa-crop-service/         # Crop Management Service (Go)
├── fa-notification-service/ # Notification Service (Go)
├── fa-payment-service/      # Payment Service (Go)
├── fa-ai-service/           # AI Service (Python)
├── docs/                    # API documentation
│   └── FarmAgent.postman_collection.json
├── docker-compose.yml       # Docker orchestration
└── README.md
```

### Running Individual Services

#### Go Services

```bash
cd fa-auth-service
go mod download
go run cmd/server/main.go
```

#### AI Service (Python)

```bash
cd fa-ai-service
python -m venv venv
source venv/bin/activate  # or `venv\Scripts\activate` on Windows
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8005
```

### Running Tests

```bash
# Auth Service
cd fa-auth-service
go test ./... -v

# All Go services
for service in fa-auth-service fa-crop-service fa-notification-service fa-payment-service; do
  echo "Testing $service..."
  (cd $service && go test ./... -v)
done
```

### Database Migrations

Migrations run automatically on service startup. To reset:

```bash
# Drop and recreate databases
docker-compose down -v
docker-compose up -d postgres mongodb
# Wait for DBs to start, then restart services
docker-compose up -d
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service

# Last 100 lines
docker-compose logs --tail=100 ai-service
```

---

## 🔧 Configuration

### Environment Variables

#### Common Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (development/production) | development |
| `APP_PORT` | Service port | varies |
| `JWT_SECRET` | JWT signing secret | (required) |

#### Database (PostgreSQL)

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_USER` | Database user | farmagent |
| `DB_PASSWORD` | Database password | (required) |
| `DB_NAME` | Database name | varies |

#### AI Service

| Variable | Description | Default |
|----------|-------------|---------|
| `OLLAMA_HOST` | Ollama API URL | <http://localhost:11434> |
| `OLLAMA_MODEL` | Chat model | llama3.2:3b |

#### Payment Service (MTN MoMo)

| Variable | Description | Default |
|----------|-------------|---------|
| `MTN_ENVIRONMENT` | sandbox/production | sandbox |
| `MTN_SUBSCRIPTION_KEY` | Primary subscription key | (required) |
| `MTN_API_USER` | API User UUID | (required) |
| `MTN_API_KEY` | API Key | (required) |

---

## 📧 Email Testing

FarmAgent uses [Mailpit](https://github.com/axllent/mailpit) for email testing in development.

**Access the Mailpit UI:** <http://localhost:8025>

All emails sent by the notification service (verification, password reset, etc.) will appear here.

---

## 🐛 Troubleshooting

### Common Issues

#### Services won't start

```bash
# Check if ports are in use
lsof -i :8000-8005

# Reset everything
docker-compose down -v
docker-compose up -d --build
```

#### AI Service shows "Model not loaded"

```bash
# Check if transformers is installed
docker-compose exec ai-service pip list | grep transformers

# Rebuild AI service
docker-compose build --no-cache ai-service
docker-compose up -d ai-service
```

#### MTN Payment fails

1. Verify credentials in `.env`
2. Check you're using sandbox phone numbers (e.g., `256772000000`)
3. View logs: `docker-compose logs payment-service`

#### Database connection errors

```bash
# Check if postgres is healthy
docker-compose ps postgres

# Restart postgres
docker-compose restart postgres
```

---

## 📊 Monitoring

### Health Endpoints

All services expose `/health` endpoints. Use for monitoring:

```bash
#!/bin/bash
# health-check.sh
services=("8000" "8001" "8002" "8003" "8004" "8005")
for port in "${services[@]}"; do
  status=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$port/health)
  echo "Port $port: $status"
done
```

### Container Stats

```bash
docker stats --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Team

Built with ❤️ for East African farmers.

---

<div align="center">

**[⬆ Back to Top](#-farmagent)**

</div>
