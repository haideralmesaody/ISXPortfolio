# Full-Stack Development Document for a Golang Web Application

## Overview
This document provides a step-by-step guide for building a full-stack web application with a Golang backend, a React frontend, and SQLite as the lightweight database. The app will use Google Auth for user registration and login and will be containerized with Docker.

---

## Tech Stack

| Component     | Technology          | Reason                                                                 |
|---------------|---------------------|------------------------------------------------------------------------|
| Backend       | Golang              | High performance, scalability, and robust support for APIs.           |
| Frontend      | React               | Flexible, ecosystem-rich, and compatible with React Native for mobile.|
| Database      | SQLite              | Lightweight, file-based, and ideal for small to medium-sized apps.    |
| Authentication| Google OAuth2       | Industry-standard authentication for easy user registration/login.     |
| Charting      | D3.js or Chart.js   | Advanced, interactive visualizations.                                 |
| Containerization | Docker           | Simplifies setup, deployment, and environment consistency.            |

---

## Project Structure

```plaintext
project/
├── backend/
│   ├── main.go
│   ├── handlers/
│   ├── models/
│   └── config/
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   └── App.js
├── docker-compose.yml

### 3. Create `docker-compose.yml`
```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
    ports:
      - "8080:8080"
    volumes:
      - sqlite_data:/app/data

  frontend:
    build:
      context: ./frontend
    ports:
      - "3000:3000"

volumes:
  sqlite_data:
    driver: local
```
└── Dockerfile
```

---

## Backend Setup (Golang)

### 1. Initialize Project
```bash
mkdir backend
cd backend
go mod init isxportfolio-backend
go get github.com/gin-gonic/gin@v1.8.1
```

### 2. Directory Structure
```plaintext
backend/
├── main.go
├── handlers/
│   ├── auth.go
├── models/
│   ├── user.go
├── config/
│   ├── database.go
```

### 3. Install Dependencies
```bash
go get github.com/gin-gonic/gin
go get github.com/mattn/go-sqlite3
```

### 4. Example `main.go`
```go
package main

import (
	"github.com/gin-gonic/gin"
	"isxportfolio-backend/handlers"
	"isxportfolio-backend/config"
)

func main() {
	config.InitDB()
	r := gin.Default()
	r.POST("/login", handlers.LoginHandler)
	r.Run(":8080")
}
```

### 5. Google OAuth Handler (`handlers/auth.go`)
```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	// Google OAuth logic here
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
```

### 6. Database Configuration (`config/database.go`)
```go
package config

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./isxportfolio.db")
	if err != nil {
		panic(err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL
	)`

	_, err = DB.Exec(createTable)
	if err != nil {
		panic(err)
	}
}
```

---

## Frontend Setup (React)

### 1. Initialize Project
```bash
npx create-react-app frontend
cd frontend
npm install axios react-google-login chart.js react-chartjs-2
```

### 2. Directory Structure
```plaintext
frontend/
├── src/
│   ├── components/
│   │   ├── Auth.js
│   │   ├── Portfolio.js
│   │   └── NewsFeed.js
│   ├── hooks/
│   ├── App.js
```

### 3. Example `Auth.js`
```javascript
import React from "react";
import { GoogleLogin } from "react-google-login";

const Auth = () => {
  const handleSuccess = (response) => {
    console.log(response);
    // Send token to backend
  };

  const handleFailure = (response) => {
    console.error(response);
  };

  return (
    <GoogleLogin
      clientId="YOUR_GOOGLE_CLIENT_ID"
      buttonText="Login with Google"
      onSuccess={handleSuccess}
      onFailure={handleFailure}
      cookiePolicy={'single_host_origin'}
    />
  );
};

export default Auth;
```

---

## Docker Setup

### 1. Create `Dockerfile` for Backend
```dockerfile
# Backend Dockerfile
FROM golang:1.19
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .
CMD ["./main"]
EXPOSE 8080
```

### 2. Create `Dockerfile` for Frontend
```dockerfile
# Frontend Dockerfile
FROM node:20
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

### 3. Create `docker-compose.yml`
```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
    ports:
      - "8080:8080"

  frontend:
    build:
      context: ./frontend
    ports:
      - "3000:3000"

volumes:
  sqlite_data:
```

---

## Deployment

1. **Start Services:**
```bash
docker-compose up --build
```

2. **Access Application:**
   - Backend: `http://localhost:8080`
   - Frontend: `http://localhost:3000`

3. **Verify Database:**
Check `isxportfolio.db` file for user data.

---

## Future Steps
1. **Extend Backend:** Add endpoints for portfolio management and news fetching.
2. **Enhance Frontend:** Build dashboards and charts using `Chart.js`.
3. **Mobile App:** Use React Native for a seamless transition.
4. **Deployment:** Use Kubernetes or Docker Swarm for scalable deployment.

