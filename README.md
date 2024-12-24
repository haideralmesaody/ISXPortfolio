# ISX Portfolio

A full-stack web application for portfolio management with Golang backend and React frontend.

## Tech Stack

- Backend: Golang with Gin framework
- Database: SQLite
- Authentication: Google OAuth2
- Docker for containerization

## Project Structure

```
project/
├── backend/
│   ├── config/      # Configuration files
│   ├── handlers/    # HTTP handlers
│   ├── models/      # Data models
│   └── main.go      # Entry point
└── docker-compose.yml
```

## Setup

1. Clone the repository:
```bash
git clone https://github.com/haideralmesaody/isxportfolio.git
```

2. Set up environment variables:
```bash
cp backend/.env.example backend/.env
# Edit backend/.env with your Google OAuth credentials
```

3. Start the application:
```bash
docker-compose up --build
```

4. Access the application:
- Backend: http://localhost:8000
- Health Check: http://localhost:8000/health
- Google Login: http://localhost:8000/auth/google/login

## Features

- [x] Google OAuth2 Authentication
- [x] Health Check Endpoint
- [ ] User Management
- [ ] Portfolio Management
- [ ] Data Visualization 