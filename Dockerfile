FROM golang:1-alpine

WORKDIR /app
# COPY go.mod .
# COPY go.sum .
COPY . .

RUN go mod tidy && go build .