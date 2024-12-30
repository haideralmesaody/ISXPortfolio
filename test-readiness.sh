#!/bin/bash

echo "Testing backend health..."
curl http://localhost:8000/health

echo "\nTesting frontend..."
curl http://localhost:3000

echo "\nChecking if services are ready..."
docker-compose ps 