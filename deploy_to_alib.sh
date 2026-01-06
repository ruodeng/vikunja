#!/bin/bash

# Configuration
REMOTE_HOST="alib"
IMAGE_NAME="vikunja:latest"
REMOTE_DIR="~/docker/vikunja"

echo ">>> Starting Deployment to $REMOTE_HOST"

# 1. Build Docker Image
echo ">>> Building Docker image: $IMAGE_NAME..."
docker build -t $IMAGE_NAME .
if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi
echo "✅ Build successful."

# 2. Upload Image
echo ">>> Uploading image to $REMOTE_HOST..."
docker save $IMAGE_NAME | gzip | ssh $REMOTE_HOST "gunzip | docker load"
if [ $? -ne 0 ]; then
    echo "❌ Upload failed!"
    exit 1
fi
echo "✅ Upload successful."

# 3. Restart Remote Service
echo ">>> Restarting remote service..."
ssh $REMOTE_HOST "cd $REMOTE_DIR && docker compose up -d --force-recreate vikunja"
if [ $? -ne 0 ]; then
    echo "❌ Restart failed!"
    exit 1
fi

echo "✅ Deployment completed successfully!"
