#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

APP_NAME="shuttlers-app"

echo "Building Docker image..."
docker build -t $APP_NAME .

echo "Stopping existing container (if any)..."
if [ "$(docker ps -aq -f name=$APP_NAME)" ]; then
    docker stop $APP_NAME || true
    docker rm $APP_NAME || true
fi

echo "Starting new container..."
docker run -d --name $APP_NAME -p 9877:9877 $APP_NAME

echo "Deployment complete! App is running on port 9877"
