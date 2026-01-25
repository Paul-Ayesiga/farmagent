# FarmAgent - Software Requirements Specification (SRS)

**Version:** 2.0  
**Date:** January 24, 2026  
**Project:** FarmAgent - Autonomous AI Agent for Small-Scale Farmers in Uganda

---

## 1. Introduction

### 1.1 Purpose

This document defines the functional and non-functional requirements for **FarmAgent**, an autonomous AI-powered agricultural assistant for Uganda's small-scale farmers.

### 1.2 Scope

FarmAgent provides:

- **Crop Disease Detection** — AI-powered image analysis
- **Treatment Recommendations** — Context-aware advice via LLM
- **Market Intelligence** — Real-time pricing and buyer matching
- **Financial Services** — Mobile money integration
- **Multi-language Support** — Luganda, Runyankole, English

### 1.3 Technology Stack

| Component | Technology |
|-----------|------------|
| Authentication | Keycloak |
| API Gateway | Go (Gin) |
| Crop Service | Go (Gin + GORM) |
| Payment Service | Go (Gin + GORM) |
| Notification Service | Go (Gin) |
| AI Agent | Python (Agno + FastAPI) |
| Mobile App | React Native |
| Databases | PostgreSQL, MongoDB, Redis |

---

## 2. User Classes

| Role | Description |
|------|-------------|
| **Farmer** | Primary user, small-scale agricultural producer |
| **Extension Officer** | Agricultural advisor, monitors multiple farmers |
| **Buyer** | Wholesale buyers, processors |
| **Admin** | System administrator |

---

## 3. Functional Requirements

### 3.1 Authentication (Keycloak)

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-AUTH-01 | User registration via phone number | High |
| FR-AUTH-02 | Login with OAuth2/OIDC | High |
| FR-AUTH-03 | Role-based access control | High |
| FR-AUTH-04 | Token refresh mechanism | High |
| FR-AUTH-05 | Password reset flow | Medium |

### 3.2 Crop Management

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-CROP-01 | Register agricultural fields (name, GPS, size, soil type) | High |
| FR-CROP-02 | Register crops (type, variety, planting date) | High |
| FR-CROP-03 | Track crop health over time (0-100 score) | High |
| FR-CROP-04 | Log treatments applied | High |
| FR-CROP-05 | View crop lifecycle status | Medium |
| FR-CROP-06 | Group farmers into cooperatives | Low |

### 3.3 AI Disease Detection

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-AI-01 | Analyze crop photos for diseases | Critical |
| FR-AI-02 | Detect 50+ diseases across major crops | High |
| FR-AI-03 | Provide confidence score (0-100%) | High |
| FR-AI-04 | Assess severity (mild/moderate/severe/critical) | High |
| FR-AI-05 | Support offline image queuing | Medium |

### 3.4 Treatment Recommendations

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-REC-01 | Recommend treatments based on diagnosis | Critical |
| FR-REC-02 | Factor in weather conditions | High |
| FR-REC-03 | Provide cost estimates in UGX | Medium |
| FR-REC-04 | Support organic/chemical preferences | Medium |
| FR-REC-05 | Schedule follow-up reminders | High |

### 3.5 Market Intelligence

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-MKT-01 | Display real-time commodity prices | High |
| FR-MKT-02 | Create sell requests | High |
| FR-MKT-03 | Match farmers with buyers | High |
| FR-MKT-04 | Manage orders (confirm, deliver, pay) | High |
| FR-MKT-05 | Show price trends | Medium |

### 3.6 Payment Services

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-PAY-01 | MTN Mobile Money integration | Critical |
| FR-PAY-02 | Airtel Money integration | High |
| FR-PAY-03 | Transaction history | High |
| FR-PAY-04 | Invoice generation | Medium |
| FR-PAY-05 | Refund processing | Medium |

### 3.7 Notifications

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-NOT-01 | SMS notifications (Africa's Talking) | Critical |
| FR-NOT-02 | Push notifications (FCM) | High |
| FR-NOT-03 | Scheduled reminders | High |
| FR-NOT-04 | User notification preferences | Medium |

---

## 4. Non-Functional Requirements

### 4.1 Performance

| Requirement | Specification |
|-------------|---------------|
| API response time | < 500ms (95th percentile) |
| Image analysis | < 5 seconds |
| Concurrent users | 1,000+ |
| Mobile app cold start | < 3 seconds |

### 4.2 Security

| Requirement | Specification |
|-------------|---------------|
| Authentication | OAuth2/OIDC via Keycloak |
| Transport | TLS 1.3 |
| Token validation | RS256 JWT |
| Rate limiting | 100 req/min per user |

### 4.3 Availability

| Requirement | Specification |
|-------------|---------------|
| Uptime | 99.5% |
| Offline mode | Core features work offline |
| Data sync | Auto-sync on connectivity |

### 4.4 Scalability

| Requirement | Specification |
|-------------|---------------|
| Horizontal scaling | Kubernetes auto-scaling |
| Target users | 50,000 within 2 years |

---

## 5. Use Cases

### UC-001: Diagnose Crop Disease

**Actor:** Farmer  
**Flow:**

1. Farmer opens app, takes photo of crop
2. App sends image to AI Agent
3. Agent analyzes, returns diagnosis + treatment
4. App displays results, schedules follow-up
5. Notification Service sends SMS reminder

### UC-002: Sell Produce

**Actor:** Farmer, Buyer  
**Flow:**

1. Farmer creates sell request (crop, quantity, price)
2. System matches with suitable buyers
3. Buyer accepts, order created
4. Farmer delivers, buyer confirms
5. Payment processed via Mobile Money

### UC-003: Make Payment

**Actor:** Farmer  
**Flow:**

1. User initiates payment
2. Payment Service calls MTN MoMo API
3. User receives USSD prompt, enters PIN
4. Callback confirms payment
5. Transaction recorded, notification sent

---

## 6. External Integrations

| System | Purpose | API Type |
|--------|---------|----------|
| Keycloak | Authentication | OAuth2/OIDC |
| MTN MoMo | Payments | REST |
| Airtel Money | Payments | REST |
| Africa's Talking | SMS, Voice | REST |
| OpenWeather | Weather data | REST |
| Claude API | AI recommendations | REST |
| Firebase | Push notifications | SDK |

---

## 7. Constraints

| Constraint | Impact |
|------------|--------|
| 2-person team | Simplified architecture |
| Low bandwidth (2G/3G) | Compressed payloads |
| Local languages | Translation required |
| Offline areas | Offline-first design |

---

## Document History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-01-23 | Initial SRS |
| 2.0 | 2026-01-24 | Revised for Go stack, Keycloak, Agno |
