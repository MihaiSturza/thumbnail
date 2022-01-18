#!/usr/bin/env bash

# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go -ldflags="-w -s -extldflags=-static"

# aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin 934008704165.dkr.ecr.eu-central-1.amazonaws.com

# aws ecr create-repository \
#     --repository-name lambdas \
#     --image-scanning-configuration scanOnPush=true \
#     --region eu-central-1

# aws ecr set-repository-policy \
#     --repository-name lambdas \
#     --policy-text file://$(pwd)/policy.json

# docker images

# docker build -t lambdas .

# docker tag lambdas:latest 934008704165.dkr.ecr.eu-central-1.amazonaws.com/lambdas:latest

# docker push 934008704165.dkr.ecr.eu-central-1.amazonaws.com/lambdas:latest


aws lambda update-function-code \
--function-name DockerTest \
--image-uri 934008704165.dkr.ecr.eu-central-1.amazonaws.com/lambdas:latest   \
--profile Admin
# --dry-run

