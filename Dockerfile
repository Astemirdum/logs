# Builder

FROM golang:1.19-alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/log ./cmd/main.go

# App

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

EXPOSE 8081

COPY --from=builder /app/log /app/log
COPY --from=builder /build/.env /app/.env

RUN mkdir configs
COPY --from=builder /build/configs /app/configs
CMD ["./log", "-c", "configs/config.yml"]
