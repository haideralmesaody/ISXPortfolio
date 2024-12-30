# Full-Stack Development Document for a Golang Web Application

## Overview
This document provides a step-by-step guide for building a full-stack web application with a Golang backend, a Flutter frontend, and SQLite as the lightweight database. The app will use Google Auth for user registration and login and will be containerized with Docker.

---

## Tech Stack

| Component     | Technology          | Reason                                                                 |
|---------------|---------------------|------------------------------------------------------------------------|
| Backend       | Golang              | High performance, scalability, and robust support for APIs.           |
| Frontend      | Flutter             | Cross-platform, high-performance, and native-like UI for web and mobile.|
| Database      | SQLite              | Lightweight, file-based, and ideal for small to medium-sized apps.    |
| Authentication| Google OAuth2       | Industry-standard authentication for easy user registration/login.     |
| Charting      | Flutter Packages    | Advanced, interactive visualizations using packages like `fl_chart`.  |
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
│   ├── lib/
│   │   ├── components/
│   │   ├── screens/
│   │   ├── services/
│   │   └── main.dart
│   ├── pubspec.yaml
│   ├── android/
│   ├── ios/
│   ├── web/
│   └── test/
├── docker-compose.yml
└── Dockerfile
```

---

## Backend Setup (Golang)

### 1. Initialize Project
```bash
mkdir backend
cd backend
go mod init isxportfolio-backend
go get -u github.com/gin-gonic/gin
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

## Frontend Setup (Flutter)

### 1. Initialize Project
```bash
flutter create frontend
cd frontend
```

### 2. Directory Structure
```plaintext
frontend/
├── lib/
│   ├── components/
│   ├── screens/
│   ├── services/
│   └── main.dart
├── pubspec.yaml
├── android/
├── ios/
├── web/
└── test/
```

### 3. Update `pubspec.yaml`
Add dependencies for HTTP requests and Google Sign-In.
```yaml
dependencies:
  flutter:
    sdk: flutter
  google_sign_in: ^5.4.2
  http: ^0.15.0
  fl_chart: ^0.40.0
```

### 4. Example `main.dart`
```dart
import 'package:flutter/material.dart';
import 'screens/login_screen.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'ISX Portfolio',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const LoginScreen(),
    );
  }
}
```

### 5. Example Login Screen (`screens/login_screen.dart`)
```dart
import 'package:flutter/material.dart';
import 'package:google_sign_in/google_sign_in.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({Key? key}) : super(key: key);

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final GoogleSignIn _googleSignIn = GoogleSignIn();

  Future<void> _handleSignIn() async {
    try {
      await _googleSignIn.signIn();
      // Handle login logic here
    } catch (error) {
      print(error);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Login')),
      body: Center(
        child: ElevatedButton(
          onPressed: _handleSignIn,
          child: const Text('Sign in with Google'),
        ),
      ),
    );
  }
}
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
FROM google/dart:stable
WORKDIR /app
COPY pubspec.* ./
RUN dart pub get
COPY . .
RUN dart pub global activate webdev
RUN dart pub global run webdev build
EXPOSE 8080
CMD ["webdev", "serve", "--hostname", "0.0.0.0"]
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
      - "3000:8080"

volumes:
  sqlite_data:
    driver: local
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
2. **Enhance Frontend:** Build dashboards and charts using `fl_chart`.
3. **Mobile App:** The Flutter frontend is already cross-platform and ready for mobile deployment.
4. **Deployment:** Use Kubernetes or Docker Swarm for scalable deployment.

