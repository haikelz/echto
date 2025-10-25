#!/bin/bash

# Swagger Documentation Test Script
# This script tests the Swagger documentation endpoints

set -e

BASE_URL="http://localhost:8080"

echo "Testing Swagger Documentation..."

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/health" | jq '.' || echo "Health check failed"

echo -e "\n2. Testing Swagger JSON endpoint..."
curl -s "$BASE_URL/swagger/doc.json" | jq '.info' || echo "Swagger JSON not found"

echo -e "\n3. Testing Swagger UI endpoint..."
SWAGGER_UI_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/swagger/index.html")
if [ "$SWAGGER_UI_RESPONSE" = "200" ]; then
    echo "Swagger UI is accessible at: $BASE_URL/swagger/index.html"
else
    echo "Swagger UI not accessible (HTTP $SWAGGER_UI_RESPONSE)"
fi

echo -e "\n4. Testing Swagger YAML endpoint..."
curl -s "$BASE_URL/swagger/doc.yaml" | head -20 || echo "Swagger YAML not found"

echo -e "\nSwagger documentation testing completed!"
echo "You can view the interactive documentation at: $BASE_URL/swagger/index.html"
