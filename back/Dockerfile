
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


# Utilisation d'une image alpine compatible ARM
FROM arm64v8/alpine:3.14

# Créer un répertoire de travail dans le conteneur
WORKDIR /app

# Copier l'exécutable Go (qui a été compilé localement pour arm64) dans le conteneur
COPY ./main .

# Rendre le binaire exécutable
RUN chmod +x ./main

# Exposer le port sur lequel l'application tourne (par exemple 3006)
EXPOSE 3006

# Commande pour exécuter l'application
CMD ["./main"]

