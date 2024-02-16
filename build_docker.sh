#!/bin/bash/

IMAGE_NAME="ffmpeg_study"

TAG="0.1"

docker build --no-cache -t ${IMAGE_NAME}:${TAG} .
