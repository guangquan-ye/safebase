# Étape de build avec Go
FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY . .

# Compiler l'application Go
RUN GOARCH=arm64 GOOS=linux go build -o main .

# Étape finale avec une version plus récente d'Alpine
FROM arm64v8/ubuntu:20.04

# Installer les clients PostgreSQL et MySQL
RUN apt-get update && apt-get install -y postgresql-client mysql-client

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3006
CMD ["./main"]
