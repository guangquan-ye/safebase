package main

import (
	"database/sql"
	"fmt"
	"log"
	"safebase/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var mainDB *sql.DB

// Structure représentant une base de données
type Database struct {
	DBName   string `json:"dbName"`
	DBType   string `json:"dbType"`
	DBPort   string `json:"dbPort"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

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

// insertDatabaseInfo insère les informations de la nouvelle base de données dans mainDB
func insertDatabaseInfo(dbName, dbType, dbPort, userName, password string) error {
	query := `INSERT INTO databases (db_name, db_type, db_port, user_name, password, ) VALUES ($1, $2, $3, $4, $5)`
	_, err := mainDB.Exec(query, dbName, dbType, dbPort, userName, password)
	if err != nil {
		return fmt.Errorf("failed to insert database info: %v", err)
	}
	return nil
}

// countDatabases retourne le nombre total d'entrées dans la table databases
func countDatabases() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM databases`
	err := mainDB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count databases: %v", err)
	}
	return count, nil
}

// getAllDatabases retourne toutes les bases de données sous forme de slice de Database
func getAllDatabases() ([]Database, error) {
	var databases []Database
	query := `SELECT db_name, db_type, db_port, user_name, password FROM databases`
	rows, err := mainDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var db Database
		if err := rows.Scan(&db.DBName, &db.DBType, &db.DBPort, &db.UserName, &db.Password); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		databases = append(databases, db)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

func main() {
	// Tenter la connexion à la base de données principale avant de démarrer l'application
	err := connectMainDB()
	if err != nil {
		log.Fatalf("Failed to connect to main database: %v", err)
	}

	// Créer une nouvelle instance de l'application Fiber
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// Rend la page index.html
		return c.Render("templates/index.html", fiber.Map{})
	})

	app.Get("/backups", func(c *fiber.Ctx) error {
		// Rend la page backups.html
		return c.Render("templates/backups.html", fiber.Map{})
	})

	app.Get("/databases", func(c *fiber.Ctx) error {
		// Rend la page databases.html
		return c.Render("templates/databases.html", fiber.Map{})
	})

	app.Get("/restores", func(c *fiber.Ctx) error {
		// Rend la page restores.html
		return c.Render("templates/restores.html", fiber.Map{})
	})

	// Route pour vérifier la connexion et enregistrer la base de données
	app.Post("/addDatabase", func(c *fiber.Ctx) error {
		// Structure pour les informations de la base de données cible
		var config struct {
			DBType   string `json:"dbType"`
			DBName   string `json:"dbName"`
			DBPort   string `json:"dbPort"`
			UserName string `json:"userName"`
			Password string `json:"password"`
		}

		// Parsing du corps de la requête JSON
		if err := c.BodyParser(&config); err != nil {
			return c.Status(400).SendString("Failed to parse JSON: " + err.Error())
		}

		// Vérification de la connexion à la base de données cible
		dynamicDB, err := database.ConnectDynamicDB(config.DBType, config.DBName, config.DBPort, config.UserName, config.Password)
		if err != nil {
			return c.Status(500).SendString("Failed to connect to dynamic database: " + err.Error())
		}
		defer dynamicDB.Close()

		// Si la connexion est réussie, insérer les informations dans mainDB
		if err := insertDatabaseInfo(config.DBName, config.DBType, config.DBPort, config.UserName, config.Password); err != nil {
			return c.Status(500).SendString("Failed to insert database info: " + err.Error())
		}

		return c.SendString("Database connected and info saved successfully!")
	})

	// Route pour récupérer toutes les bases de données en JSON
	app.Get("/getDatabases", func(c *fiber.Ctx) error {
		databases, err := getAllDatabases()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve databases"})
		}

		return c.JSON(databases)
	})

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

	// Route pour récupérer le nombre total d'entrées dans la table databases
	app.Get("/getDbCount", func(c *fiber.Ctx) error {
		count, err := countDatabases()
		if err != nil {
			return c.Status(500).SendString("Failed to count databases: " + err.Error())
		}

		return c.SendString(fmt.Sprintf("%d", count))
	})

	// Route pour faire un dump de la base de données
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

	// Démarrer le serveur après avoir défini toutes les routes
	log.Fatal(app.Listen(":8080"))
}
