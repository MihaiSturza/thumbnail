# thumbnail
Create thumbnail using golang

# docker file config

FROM golang:1.15-alpine3.13 as build

RUN apk add --no-cache build-base imagemagick-dev imagemagick

COPY . /code
WORKDIR /code
ENV CGO_ENABLED=1
ENV CGO_CFLAGS_ALLOW=-Xpreprocessor
ENV GOPROXY=https://proxy.golang.org
RUN go mod download && \
    go test -cover ./... && \
    go build -a -installsuffix cgo -o app

FROM alpine:3.13
RUN apk add --no-cache tzdata ca-certificates imagemagick-dev imagemagick
COPY --from=build /code/app /app/app
COPY --from=build /code/api/docs /app/api/docs

WORKDIR /app
ENTRYPOINT [ "/app/app" ]


/ You used to build this programm some CI function with golang-alpine base image, so try to use full golang image instead.
