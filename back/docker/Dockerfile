# Étape de build avec Go
FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY . .

# Compiler l'application Go
RUN GOARCH=arm64 GOOS=linux go build -o main .

# Étape finale avec une version plus récente d'Alpine
FROM arm64v8/ubuntu:20.04

# Installer les clients PostgreSQL et MySQL et ajouter les dépôts PostgreSQL pour obtenir la version 15 du client

# RUN apt-get update && apt-get install -y postgresql-client mysql-client

RUN apt-get update && apt-get install -y wget gnupg2 && \
    echo "deb http://apt.postgresql.org/pub/repos/apt focal-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    apt-get update && apt-get install -y postgresql-client-15 mysql-client

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 3006
CMD ["./main"]
