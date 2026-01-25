# FA-Payment-Service

MTN Mobile Money payment integration service for subscriptions and transactions.

## Overview

This service handles:

- **MTN MoMo Collections** - Request payments from users
- **Transaction Management** - Track payment status and history
- **Subscriptions** - Manage user subscription plans

## Tech Stack

- **Framework**: Go 1.22+, Chi Router
- **Database**: PostgreSQL
- **Payment**: MTN MoMo API

## Quick Start

### Prerequisites

1. Register at [momodeveloper.mtn.com](https://momodeveloper.mtn.com)
2. Subscribe to Collections API
3. Create API User and get API Key

### Local Development

```bash
# Set environment
export MTN_SUBSCRIPTION_KEY=your-key
export MTN_API_USER=your-uuid
export MTN_API_KEY=your-api-key

# Run service
go run cmd/server/main.go
```

### With Docker

```bash
docker build -t fa-payment-service .
docker run -p 8004:8004 \
  -e MTN_SUBSCRIPTION_KEY=xxx \
  -e MTN_API_USER=xxx \
  -e MTN_API_KEY=xxx \
  fa-payment-service
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/api/v1/payments/initiate` | POST | Start MoMo payment |
| `/api/v1/payments/{id}/status` | GET | Check transaction status |
| `/api/v1/payments/callback` | POST | MTN webhook |
| `/api/v1/payments/history` | GET | User transactions |
| `/api/v1/subscriptions` | GET | Current subscription |
| `/api/v1/subscriptions/plans` | GET | Available plans |
| `/api/v1/subscriptions/subscribe` | POST | Subscribe to plan |

## Subscription Plans

| Plan | Price (UGX) | Features |
|------|-------------|----------|
| Free | 0 | 2 scans/month |
| Basic | 5,000 | 10 scans, recommendations |
| Premium | 20,000 | Unlimited, AI chat, alerts |

## Environment Variables

```env
APP_PORT=8004

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=farmagent
DB_PASSWORD=secret
DB_NAME=farmagent_payments

# MTN MoMo
MTN_ENVIRONMENT=sandbox
MTN_BASE_URL=https://sandbox.momodeveloper.mtn.com
MTN_SUBSCRIPTION_KEY=your-key
MTN_API_USER=your-uuid
MTN_API_KEY=your-api-key
MTN_CALLBACK_URL=http://your-domain.com/api/v1/payments/callback
```

## Testing (Sandbox)

Use sandbox test numbers:

- `256772000000` - Successful payment
- `256772000001` - Failed payment

```bash
curl -X POST http://localhost:8004/api/v1/payments/initiate \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -d '{
    "amount": 5000,
    "phone": "256772000000",
    "reason": "Test payment"
  }'
```

## Payment Flow

1. Client calls `/payments/initiate`
2. Service calls MTN `requesttopay` API
3. User receives USSD prompt on phone
4. User enters PIN to confirm
5. MTN sends callback to `/payments/callback`
6. Service updates transaction status
7. If subscription payment, subscription is activated
