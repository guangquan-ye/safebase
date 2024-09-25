package main

import (
	"database/sql"
	"fmt"
	"log"
	"safebase/database"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var mainDB *sql.DB

// Database structure representing a database
type Database struct {
	DBId     int    `json:"id"`
	DBName   string `json:"dbName"`
	DBType   string `json:"dbType"`
	DBPort   string `json:"dbPort"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// Common structure for DB connection
type DBConnection struct {
	DBType   string `json:"dbType"`
	DBName   string `json:"dbName"`
	DBPort   string `json:"dbPort"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// Function to connect to the main database
func connectMainDB() error {
	connStr := "user=admin password=securepassword dbname=safebase host=safebase port=5432 sslmode=disable"
	var err error
	mainDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to main database: %v", err)
	}

	if err := mainDB.Ping(); err != nil {
		return fmt.Errorf("unable to ping main database: %v", err)
	}

	fmt.Println("Connected to the main database successfully!")
	return nil
}

// Insert new database info into mainDB
func insertDatabaseInfo(db Database) error {
	query := `INSERT INTO databases (dbName, dbType, dbPort, userName, password) VALUES ($1, $2, $3, $4, $5)`
	_, err := mainDB.Exec(query, db.DBName, db.DBType, db.DBPort, db.UserName, db.Password)
	if err != nil {
		return fmt.Errorf("failed to insert database info: %v", err)
	}
	return nil
}

// Count total databases
func countDatabases() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM databases`
	err := mainDB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count databases: %v", err)
	}
	return count, nil
}

// Retrieve all databases from the main DB
func getAllDatabases() ([]Database, error) {
	var databases []Database
	query := `SELECT id, dbname, dbtype, dbport, username, password FROM databases`
	rows, err := mainDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get databases: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var db Database
		if err := rows.Scan(&db.DBId, &db.DBName, &db.DBType, &db.DBPort, &db.UserName, &db.Password); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		databases = append(databases, db)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

// Delete database by ID
func deleteDatabase(id int) error {
	query := `DELETE FROM databases WHERE id = $1`
	_, err := mainDB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete database: %v", err)
	}
	return nil
}

// Unified error response handler
func errorResponse(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{"status": "error", "message": message})
}

func main() {
	// Connect to the main database
	if err := connectMainDB(); err != nil {
		log.Fatalf("Failed to connect to main database: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New()

	// Basic routes for rendering pages
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("templates/index.html", fiber.Map{})
	})

	app.Get("/backups", func(c *fiber.Ctx) error {
		return c.Render("templates/backups.html", fiber.Map{})
	})

	app.Get("/databases", func(c *fiber.Ctx) error {
		return c.Render("templates/databases.html", fiber.Map{})
	})

	app.Get("/restores", func(c *fiber.Ctx) error {
		return c.Render("templates/restores.html", fiber.Map{})
	})

	// Add database route
	app.Post("/addDatabase", func(c *fiber.Ctx) error {
		var db Database
		if err := c.BodyParser(&db); err != nil {
			return errorResponse(c, 400, "Invalid input")
		}

		log.Printf("Received data: DBType=%s, DBName=%s, DBPort=%s, UserName=%s\n", db.DBType, db.DBName, db.DBPort, db.UserName)

		dynamicDB, err := database.ConnectDynamicDB(db.DBType, db.DBName, db.DBPort, db.UserName, db.Password)
		if err != nil {
			return errorResponse(c, 500, "Failed to connect to dynamic database: "+err.Error())
		}
		defer dynamicDB.Close()

		if err := insertDatabaseInfo(db); err != nil {
			return errorResponse(c, 500, "Failed to insert database info: "+err.Error())
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Database connected and info saved successfully!"})
	})

	// Route to get all databases
	app.Get("/getDatabases", func(c *fiber.Ctx) error {
		databases, err := getAllDatabases()
		if err != nil {
			return errorResponse(c, 500, "Failed to retrieve databases")
		}
		return c.JSON(databases)
	})

	// Dynamic connection route
	app.Post("/connexion", func(c *fiber.Ctx) error {
		var config DBConnection
		if err := c.BodyParser(&config); err != nil {
			return errorResponse(c, 400, "Failed to parse JSON: "+err.Error())
		}

		dynamicDB, err := database.ConnectDynamicDB(config.DBType, config.DBName, config.DBPort, config.UserName, config.Password)
		if err != nil {
			return errorResponse(c, 500, "Failed to connect to dynamic database: "+err.Error())
		}
		defer dynamicDB.Close()

		return c.SendString("Connected to the database successfully!")
	})

	// Get total count of databases
	app.Get("/getDbCount", func(c *fiber.Ctx) error {
		count, err := countDatabases()
		if err != nil {
			return errorResponse(c, 500, "Failed to count databases: "+err.Error())
		}

		return c.SendString(fmt.Sprintf("%d", count))
	})

	// Dump database route
	app.Post("/dump", func(c *fiber.Ctx) error {
		var dataDump DBConnection
		if err := c.BodyParser(&dataDump); err != nil {
			return errorResponse(c, 400, "Failed to parse JSON: "+err.Error())
		}

		err := database.DumpBdd(dataDump.DBType, dataDump.DBName, dataDump.DBPort, dataDump.UserName, dataDump.Password)
		if err != nil {
			return c.Status(500).SendString("Failed to dump Database: " + err.Error())
		}

		return c.SendString("Dumped the database successfully!")
	})

	// Delete database by ID route
	app.Delete("/deleteDatabase/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return errorResponse(c, 400, "Invalid database ID")
		}

		if err := deleteDatabase(id); err != nil {
			return errorResponse(c, 500, "Failed to delete database: "+err.Error())
		}

		return c.SendString("Database deleted successfully!")
	})

	// Start the Fiber app
	log.Fatal(app.Listen(":8080"))
}
