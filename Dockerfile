FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY . .

RUN go mod download && go mod tidy && go build .