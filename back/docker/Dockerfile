
# # Utilise une image Alpine légère
# FROM alpine:3.14

# # Installe curl avec le gestionnaire de paquets apk
# RUN apk update && apk add --no-cache mysql-client postgresql postgresql-client

# # Crée un répertoire pour l'application dans le conteneur
# WORKDIR /app

# # Copie le binaire Go que tu as déjà compilé localement dans le conteneur
# COPY ./main .

# # Rendre le binaire exécutable (au cas où)
# RUN chmod +x ./main

# # Expose le port sur lequel l'application écoute
# EXPOSE 3006

# # Commande pour exécuter le binaire
# CMD ["./main"]

# --------------
# # Utilisation d'une image alpine compatible ARM
# FROM arm64v8/alpine:3.14

# # Créer un répertoire de travail dans le conteneur
# WORKDIR /app


# # Rendre le binaire exécutable
# RUN chmod +x ./main
# RUN GOARCH=arm64 GOOS=linux go build -o main .

# # Copier l'exécutable Go (qui a été compilé localement pour arm64) dans le conteneur
# COPY ../main .

# # Exposer le port sur lequel l'application tourne (par exemple 3006)
# EXPOSE 3006

# # Commande pour exécuter l'application
# CMD ["./main"]

# ------------------------

# Étape 1 : Construction de l'application Go
FROM golang:1.23.0-alpine AS builder

# Définir le répertoire de travail
WORKDIR /app

# Copier le code source dans le conteneur
COPY . .
COPY --from=builder /app/main /app/main

# RUN chmod +x ./main
# RUN chmod +x ./app/main

# Compiler l'application Go
RUN GOARCH=arm64 GOOS=linux go build -o main .

# Étape 2 : Création de l'image finale
FROM arm64v8/alpine:3.14

# Installer les clients pour PostgreSQL et MySQL
RUN apk update && apk add --no-cache mysql-client postgresql-client

# Créer un répertoire de travail dans le conteneur
WORKDIR /app

# Copier le binaire construit depuis l'étape de build
COPY --from=builder /app/main .

# Exposer le port de l'application
EXPOSE 3006

# Commande pour exécuter l'application
CMD ["./main"]
