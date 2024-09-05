package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var mainDB *sql.DB

// connectMainDB établit une connexion à la base de données principale
func connectMainDB() error {
	connStr := "user=admin password=securepassword dbname=safebase sslmode=disable"
	var err error
	mainDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to main database: %v", err)
	}

	// Vérifiez que la connexion est bien établie
	if err := mainDB.Ping(); err != nil {
		return fmt.Errorf("unable to ping main database: %v", err)
	}

	fmt.Println("Connected to the main database successfully!")
	return nil
}

// connectDynamicDB établit une connexion dynamique à une autre base de données
func connectDynamicDB(userName, password, dbName, dbPort string) (*sql.DB, error) {
	// Construction de la chaîne de connexion
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", userName, password, dbName, dbPort)

	// Ouverture de la connexion
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %s: %v", dbName, err)
	}

	// Vérifiez que la connexion est bien établie
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database %s: %v", dbName, err)
	}

	fmt.Printf("Connected to the database %s successfully!\n", dbName)
	return db, nil
}

func main() {
	// Connexion à la base de données principale
	if err := connectMainDB(); err != nil {
		log.Fatalf("Error connecting to main database: %v", err)
	}
	defer mainDB.Close()

	app := fiber.New()

	// Route pour exécuter une requête sur la base de données principale
	app.Get("/main", func(c *fiber.Ctx) error {
		var result string
		err := mainDB.QueryRow("SELECT 'Main DB connection active'").Scan(&result)
		if err != nil {
			return c.Status(500).SendString("Query failed on main database: " + err.Error())
		}
		return c.SendString(result)
	})

	// Route pour se connecter à une autre base de données dynamiquement et insérer un utilisateur
	app.Post("/dynamic/adduser", func(c *fiber.Ctx) error {
		// Récupérer les informations de la base de données cible et de l'utilisateur
		var config struct {
			DBName   string `json:"dbName"`
			DBPort   string `json:"dbPort"`
			UserName string `json:"userName"`
			Password string `json:"password"`
			ID       int    `json:"id"`
			Name     string `json:"name"`
		}

		if err := c.BodyParser(&config); err != nil {
			return c.Status(400).SendString("Failed to parse JSON: " + err.Error())
		}

		// Se connecter dynamiquement à la base de données cible
		dynamicDB, err := connectDynamicDB(config.UserName, config.Password, config.DBName, config.DBPort)
		if err != nil {
			return c.Status(500).SendString("Failed to connect to dynamic database: " + err.Error())
		}
		defer dynamicDB.Close()

		// Exécuter la requête INSERT dans la base de données cible
		query := "INSERT INTO users (id, name) VALUES ($1, $2)"
		_, err = dynamicDB.Exec(query, config.ID, config.Name)
		if err != nil {
			return c.Status(500).SendString("Failed to insert user in dynamic database: " + err.Error())
		}

		return c.SendString("User added to dynamic database successfully")
	})

	// Démarrer le serveur
	log.Fatal(app.Listen(":8080"))
}
