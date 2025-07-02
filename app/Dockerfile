#syntax=docker/dockerfile:1.7-labs
FROM golang:1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY app/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o application

# Copy static assets after
# Lets us change them without affecting build cache
COPY app/static/ ./static

EXPOSE 8080

CMD ["./application"]