package database

import (
	"database/sql"
	"fmt"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// connectDynamicDB établit une connexion dynamique à une autre base de données
func ConnectDynamicDB(dbType, dbName, dbPort, userName, password string) (*sql.DB, error) {
	var connStr string

	switch dbType {
	case "postgres":
		connStr = fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", userName, password, dbName, dbPort)

	case "mysql":
		connStr = fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s", userName, password, dbPort, dbName)

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	db, err := sql.Open(dbType, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %s: %v", dbName, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database %s: %v", dbName, err)
	}

	fmt.Printf("Connected to the %s database %s successfully!\n", dbType, dbName)
	return db, nil
}

func DumpBdd(dbType, dbName, dbPort, userName, password string) error {
	var dumpCmd *exec.Cmd

	fmt.Print("dbtype: ", dbType)
	fmt.Print("dbName: ", dbName)
	fmt.Print("dbPort: ", dbPort)
	fmt.Print("userName: ", userName)
	fmt.Print("password: ", password)

	outputFile := fmt.Sprintf("/dumpFiles/%s.sql", dbName) // Crée le chemin du fichier

	switch dbType {
	case "postgres":

		// Commande pour faire un dump PostgreSQL
		dumpCmd = exec.Command("docker", "exec", "safebase", "pg_dump", "-U", userName, "-h", "localhost", "-p", dbPort, dbName, "-f", outputFile)

	case "mysql":
		// Commande pour faire un dump MySQL
		dumpCmd = exec.Command("mysqldump", "-u", userName, "-p"+password, dbName, "-r", outputFile)

	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Exécution de la commande Bash
	dumpCmd.Env = append(dumpCmd.Env, fmt.Sprintf("PGPASSWORD=%s", password)) // Pour PostgreSQL

	output, err := dumpCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to dump database: %v - %s", err, string(output))
	}

	fmt.Printf("Database %s dumped successfully!\n", dbName)
	return nil
}
