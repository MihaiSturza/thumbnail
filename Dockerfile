FROM golang:1.17-alpine as build
WORKDIR /app
# add libvips for bimg pdf support
RUN apk add vips
# cache dependencies
ADD go.mod go.sum ./
RUN go mod download 
# push files
ADD . .
# build
ENV GOOS linux
ENV CGO_ENABLED 0
ENV GOARCH amd64
RUN go build -o /function ./cmd/main.go

# copy artifacts to a clean image
FROM amazon/aws-lambda-go:latest
COPY --from=build /function ${LAMBDA_TASK_ROOT}
CMD [ "function" ]