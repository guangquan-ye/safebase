package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	// Connexion à la base de données PostgreSQL
	connStr := "user=admin password=securepassword dbname=safebase sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Création de l'application Fiber
	app := fiber.New()

	// Définition de la route "/"
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("The /!")
	})

	// Définition de la route "/hello"Ò
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	// // Démarrage du serveur sur le port 8080
	err = app.Listen(":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
