# Sound Horee Backend Service

Backend service for **Sound Horee** (PT Algoritma Kita Digital), an Android application that converts payment notifications (QRIS/Bank) into voice alerts. This service handles user synchronization, transaction history, reporting, and subscription management using **Golang** and **MySQL**.

## 🛠 Tech Stack

- **Language:** Golang (1.20+)
- **Framework:** Gin Web Framework
- **Database:** MySQL (GORM)
- **Authentication:** JWT (JSON Web Token)
- **Configuration:** Godotenv
- **Testing:** Bash Script Automation

---

## Getting Started

### 1. Prerequisites
- [Go](https://go.dev/dl/)
- MySQL Server
- Git

### 2. Installation & Setup
```bash
# Clone repository
git clone [https://github.com/algoritma-kita/sound-horee-backend.git](https://github.com/algoritma-kita/sound-horee-backend.git)
cd sound-horee-backend

# Install dependencies
go mod tidy

# Create .env file
cp .env.example .env


Automated Testing & Simulation

We have included a powerful simulation script (test.sh) to test the API End-to-End without using the Android App.

1. Setup Permission
Before running the simulation for the first time:

chmod +x test.sh

2. Usage Scenarios
You can run the script with different flags to simulate specific user behaviors:

Scenario,Command,Description
Mode Default | ./test.sh | Creates a random user & normal transactions.
Custom User | ./test.sh -u yudha_01 | Simulates a specific User ID (Persistent Data).
Mode Sultan | ./test.sh -s rich | Simulates high-value transactions (Millions/Billions).
Trial Habis | ./test.sh -s locked| Simulates locked transactions (is_trial_limited=true).

3. Command Flags Reference
Flag,Description,Example
-u,Set User ID,-u yudha_akd
-n,Set Store Name,"-n ""Toko Yudha"""
-s,"Set Scenario (normal, rich, locked)",-s locked

API Endpoints Overview
Method,Endpoint,Description,Auth Required
POST,/api/v1/auth/login,Login or Register (Google Auth),❌
POST,/api/v1/profile/sync,Sync/Update Store Profile,✅
GET,/api/v1/profile/:uid,Get Store Profile,✅
POST,/api/v1/transactions/sync,Sync Offline Transactions (Batch),✅
GET,/api/v1/transactions,Get History (Filter by Day/Month),✅

Project Structure
sound-horee-backend/
├── config/         # Database connection & env loader
├── controllers/    # Business logic & request handlers
├── middlewares/    # Auth & JWT validation
├── models/         # Database structs & schemas
├── routes/         # API URL definitions
├── utils/          # Helper functions (JWT, Formatting)
├── main.go         # App entry point
└── test.sh         # Simulation script


© 2026 PT Algoritma Kita Digital