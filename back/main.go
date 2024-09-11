package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"safebase/database"

	_ "github.com/go-sql-driver/mysql"
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

func main() {
	// Connexion à la base de données principale
	if err := connectMainDB(); err != nil {
		log.Fatalf("Error connecting to main database: %v", err)
	}
	defer mainDB.Close()

	app := fiber.New()

	cmd := exec.Command("bash", "-c", "echo Hello, World!")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la commande:", err)
		return
	}
	fmt.Println(string(output))

	// ================================================================================
	// ================================Route en GET ===================================

	// Route pour exécuter une requête sur la base de données principale
	app.Get("/", func(c *fiber.Ctx) error {
		var result string
		err := mainDB.QueryRow("SELECT 'Main DB connection active'").Scan(&result)
		if err != nil {
			return c.Status(500).SendString("Query failed on main database: " + err.Error())
		}
		return c.SendString(result)
	})

	// ================================Route en GET ===================================
	// ================================================================================

	// ================================================================================
	// ================================Route en POST ===================================

	// Route pour se connecter à une autre base de données dynamiquement
	app.Post("/connexion", func(c *fiber.Ctx) error {
		// Récupérer les informations de la base de données cible et de l'utilisateur
		var config struct {
			DBType   string `json:"dbType"`
			DBName   string `json:"dbName"`
			DBPort   string `json:"dbPort"`
			UserName string `json:"userName"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&config); err != nil {
			return c.Status(400).SendString("Failed to parse JSON: " + err.Error())
		}

		// Se connecter dynamiquement à la base de données cible
		dynamicDB, err := database.ConnectDynamicDB(config.DBType, config.DBName, config.DBPort, config.UserName, config.Password)
		if err != nil {
			return c.Status(500).SendString("Failed to connect to dynamic database: " + err.Error())
		}
		defer dynamicDB.Close()

		return c.SendString("Connected to the database successfully!")

	})

	app.Post("/getAll", func(c *fiber.Ctx) error {
		// Récupérer les informations de la base de données cible et de l'utilisateur
		var config struct {
			DBType   string `json:"dbType"`
			DBName   string `json:"dbName"`
			DBPort   string `json:"dbPort"`
			UserName string `json:"userName"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&config); err != nil {
			return c.Status(400).SendString("Failed to parse JSON: " + err.Error())
		}

		// Se connecter dynamiquement à la base de données cible
		dynamicDB, err := database.ConnectDynamicDB(config.DBType, config.DBName, config.DBPort, config.UserName, config.Password)
		if err != nil {
			return c.Status(500).SendString("Failed to connect to dynamic database: " + err.Error())
		}
		defer dynamicDB.Close()

		return c.SendString("Connected to the database successfully!")
	})

	app.Post("/dump", func(c *fiber.Ctx) error {

		var dataDump struct {
			DBType   string `json:"dbType"`
			DBName   string `json:"dbName"`
			DBPort   string `json:"dbPort"`
			UserName string `json:"userName"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&dataDump); err != nil {
			return c.Status(400).SendString("Failed to parse JSON: " + err.Error())
		}
		// Récupérer les informations de la base de données cible et de l'utilisateur
		err := database.DumpBdd(dataDump.DBType, dataDump.DBName, dataDump.DBPort, dataDump.UserName, dataDump.Password)
		if err != nil {
			return c.Status(500).SendString("Failed to dump Database: " + err.Error())
		}

		return c.SendString("Dumped the database successfully!")
	})
	// ================================Route en POST ===================================
	// ================================================================================

	// Démarrer le serveur
	log.Fatal(app.Listen(":8080"))
}
