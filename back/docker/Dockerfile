# Étape de build avec Go
FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY . .

# Compiler l'application Go
RUN GOARCH=arm64 GOOS=linux go build -o main .

# Étape finale avec une version plus récente d'Alpine
FROM arm64v8/alpine:3.18

# Installer les clients PostgreSQL et MySQL
RUN apk update && apk add --no-cache postgresql15-client mysql-client

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3006
CMD ["./main"]
