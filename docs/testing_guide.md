# FarmAgent - Developer Testing Guide

## Quick Reference

| Service | Port | Health Check |
|---------|------|--------------|
| Gateway | 8000 | `curl http://localhost:8000/health` |
| Auth | 8001 | `curl http://localhost:8001/health` |
| Crop | 8002 | `curl http://localhost:8002/health` |
| Notification | 8003 | `curl http://localhost:8003/health` |
| RabbitMQ UI | 15672 | <http://localhost:15672> (guest:guest) |
| Mailpit UI | 8025 | <http://localhost:8025> |

---

## 1. Start All Services

```bash
cd /Volumes/Personal/farmagent
docker-compose up -d
docker-compose logs -f  # Watch logs
```

---

## 2. Auth Service Endpoints

### Register User

```bash
curl -X POST http://localhost:8001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "farmer@example.com",
    "phone": "+256700123456",
    "password": "password123",
    "first_name": "John",
    "last_name": "Mukasa",
    "role": "farmer"
  }'
```

### Login

```bash
curl -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "farmer@example.com",
    "password": "password123"
  }'
```

### Get Profile (Protected)

```bash
# Replace <TOKEN> with access_token from login
curl http://localhost:8001/auth/me \
  -H "Authorization: Bearer <TOKEN>"
```

### Update Profile

```bash
curl -X PUT http://localhost:8001/auth/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "district": "Wakiso",
    "farm_size": 5.5
  }'
```

### Refresh Token

```bash
curl -X POST http://localhost:8001/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<REFRESH_TOKEN>"}'
```

### Forgot Password

```bash
curl -X POST http://localhost:8001/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"identifier": "farmer@example.com"}'
```

### Assign Role (Admin Only)

```bash
curl -X PUT http://localhost:8001/auth/users/<USER_ID>/role \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -d '{"role": "extension_officer"}'
```

---

## 3. Crop Service Endpoints

### Create Field

```bash
curl -X POST http://localhost:8002/fields \
  -H "Content-Type: application/json" \
  -H "X-User-ID: <USER_ID>" \
  -d '{
    "name": "Main Farm",
    "latitude": 0.3476,
    "longitude": 32.5825,
    "size_acres": 2.5,
    "soil_type": "loam"
  }'
```

### List Fields

```bash
curl http://localhost:8002/fields \
  -H "X-User-ID: <USER_ID>"
```

### Create Crop

```bash
curl -X POST http://localhost:8002/crops \
  -H "Content-Type: application/json" \
  -H "X-User-ID: <USER_ID>" \
  -d '{
    "field_id": "<FIELD_ID>",
    "crop_type": "maize",
    "variety": "Longe 5",
    "planting_date": "2026-01-15"
  }'
```

### List Crops

```bash
curl http://localhost:8002/crops \
  -H "X-User-ID: <USER_ID>"
```

### Create Health Record

```bash
curl -X POST http://localhost:8002/health-records \
  -H "Content-Type: application/json" \
  -H "X-User-ID: <USER_ID>" \
  -d '{
    "crop_id": "<CROP_ID>",
    "health_score": 85,
    "image_url": "https://storage.example.com/image.jpg",
    "disease_detected": "None"
  }'
```

### Create Treatment

```bash
curl -X POST http://localhost:8002/treatments \
  -H "Content-Type: application/json" \
  -H "X-User-ID: <USER_ID>" \
  -d '{
    "crop_id": "<CROP_ID>",
    "disease_name": "Fall Armyworm",
    "treatment_name": "Neem Oil Spray",
    "treatment_type": "organic",
    "application_date": "2026-01-20"
  }'
```

---

## 4. Notification Service Endpoints

### Send Email (via Queue)

```bash
curl -X POST http://localhost:8003/notifications/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "subject": "Test Email",
    "template": "welcome",
    "data": {"name": "John"}
  }'
```

### Send SMS (via Queue)

```bash
curl -X POST http://localhost:8003/notifications/sms \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+256700123456",
    "message": "Hello from FarmAgent!"
  }'
```

**Check sent emails**: <http://localhost:8025>

---

## 4. Run Unit Tests

### Auth Service

```bash
cd fa-auth-service
go test ./tests/... -v
```

### Notification Service

```bash
cd fa-notification-service
go test ./... -v
```

---

## 5. Generate Swagger Docs

### Install Swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Docs

```bash
cd fa-auth-service
swag init -g cmd/server/main.go -o docs
```

### View Swagger UI

Restart service, then visit: <http://localhost:8001/swagger/>

---

## 6. Export Postman Collection

1. Visit Swagger UI: <http://localhost:8001/swagger/>
2. Open: <http://localhost:8001/swagger/doc.json>
3. In Postman: **Import → Link → paste URL**

---

## 7. Verify Email Flow

1. Register a user with email
2. Check RabbitMQ: <http://localhost:15672> → Queues → `email_queue`
3. Check Mailpit: <http://localhost:8025> → See welcome email

---

## 8. Common Issues

| Issue | Solution |
|-------|----------|
| Services not starting | `docker-compose logs <service>` |
| Can't connect to DB | Wait for healthchecks, check `docker ps` |
| Emails not sending | Ensure notification-service is running |
| Token invalid | Token expired, refresh or re-login |

---

## 9. User Roles

| Role | Value | Permissions |
|------|-------|-------------|
| Farmer | `farmer` | Default, basic access |
| Extension Officer | `extension_officer` | Monitor farmers |
| Buyer | `buyer` | Market features |
| Admin | `admin` | Full access, assign roles |
