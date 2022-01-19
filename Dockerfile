FROM golang:1.17-alpine as build
WORKDIR /app
# cache dependencies
ENV GOOS linux
ENV CGO_ENABLED=0
ENV GOARCH amd64
ADD go.mod go.sum ./
RUN go mod download 
# push files
ADD . .
# build
RUN go build -o /function ./cmd/main.go
# copy artifacts to a clean image
FROM amazon/aws-lambda-go:latest
COPY --from=build /function ${LAMBDA_TASK_ROOT}
CMD [ "function" ]




# FROM golang:1.17-alpine as build
# WORKDIR /app
# RUN apk add --no-cache build-base imagemagick-dev imagemagick
# RUN apk add --no-cache tzdata ca-certificates imagemagick-dev imagemagick

# # set environment variables
# ENV GOOS linux
# ENV CGO_ENABLED=1
# ENV GOARCH amd64
# ENV CGO_CFLAGS_ALLOW=-Xpreprocessor

# ADD go.mod go.sum ./
# RUN go mod download 
# # push files
# ADD . .

# # build
# RUN go build -a -installsuffix cgo -o /function ./cmd/main.go

# # copy artifacts to a clean image
# FROM amazon/aws-lambda-go:latest
# COPY --from=build /function ${LAMBDA_TASK_ROOT}
# CMD [ "function" ]