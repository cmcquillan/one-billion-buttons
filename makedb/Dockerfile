FROM golang:1-bookworm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY makedb/*.go ./

RUN GOOS=linux go build -o application

# Copy static assets after
# Lets us change them without affecting build cache
COPY makedb/migrations/*.sql ./migrations

EXPOSE 8080

CMD ["./application"]