package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

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
	now := time.Now()
	date := now.Format("02-01-2006_15:04:05")
	fileName := dbName + "_" + date + ".sql"

	switch dbType {
	case "postgres":
		// Commande pour faire un dump PostgreSQL
		dumpCmd = exec.Command(
			"docker", "exec", dbName,
			"sh", "-c", fmt.Sprintf(
				"PGPASSWORD='%s' pg_dump -U %s -h localhost -p %s %s > /backups/%s",
				password, userName, dbPort, dbName, fileName,
			),
		)

	case "mysql":
		// Commande pour faire un dump MySQL
		dumpCmd = exec.Command(
			"docker", "exec", "huawey",
			"sh", "-c", fmt.Sprintf(
				"mysqldump -u %s -p%s %s > /backups/%s",
				userName, password, dbName, fileName,
			),
		)

	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Créer le fichier de dump
	outputFile, err := os.Create("dumpFiles/" + fileName)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Rediriger la sortie de la commande vers le fichier
	dumpCmd.Stdout = outputFile
	// Capturer la sortie d'erreur (stderr)
	var stderr bytes.Buffer
	dumpCmd.Stderr = &stderr

	// Exécuter la commande
	if err := dumpCmd.Run(); err != nil {
		// Afficher l'erreur avec le contenu de stderr
		return fmt.Errorf("failed to dump database: %v - %s", err, stderr.String())
	}

	fmt.Printf("Database %s dumped successfully!\n", dbName)
	return nil
}
