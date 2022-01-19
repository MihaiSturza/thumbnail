#!/usr/bin/env bash

# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go -ldflags="-w -s -extldflags=-static"

aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin 934008704165.dkr.ecr.eu-central-1.amazonaws.com

# aws ecr create-repository \
#     --repository-name thumbnail \
#     --image-scanning-configuration scanOnPush=true \
#     --region eu-central-1

# aws ecr set-repository-policy \
#     --repository-name thumbnail \
#     --policy-text file://$(pwd)/policy.json

# docker images


# aws ecr batch-delete-image \
#       --repository-name thumbnail \
#       --image-ids imageTag=latest \
#       --region region

# we build using a new tag
docker build -t pdf1 .

# we add the newtag:latest to repo:newtag
docker tag pdf1:latest 934008704165.dkr.ecr.eu-central-1.amazonaws.com/thumbnail:pdf1

# push the /repo:newtag
docker push 934008704165.dkr.ecr.eu-central-1.amazonaws.com/thumbnail:pdf1

# update function with /repo:newtag
aws lambda update-function-code \
--function-name DockerTest \
--image-uri 934008704165.dkr.ecr.eu-central-1.amazonaws.com/thumbnail:pdf1   \
--profile Admin
# --dry-run

