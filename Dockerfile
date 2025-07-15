# Stage 1: Build
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o full-project-mock ./cmd/

EXPOSE 8080

CMD ["./full-project-mock"]