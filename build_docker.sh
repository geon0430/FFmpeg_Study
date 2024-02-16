#!/bin/bash/

IMAGE_NAME="hub.inbic.duckdns.org/dev-1-team/go_vms"

TAG="goCuda_FFmpeg-0.2"

docker build --no-cache -t ${IMAGE_NAME}:${TAG} .
