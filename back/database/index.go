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
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, password, dbPort, dbPort, dbName)

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

func InsertDumpRoute(dbType, dbName, dbPort, userName, password string) error {
	// // Vérifier le type de base de données
	// if dbType != "postgres" {
	// 	return fmt.Errorf("unsupported database type: %s", dbType)
	// }

	// // Connexion à PostgreSQL
	// connStr := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", dbPort, userName, password, dbName)
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to the database: %v", err)
	// }
	// defer db.Close()

	// // Vérifier la connexion
	// if err := db.Ping(); err != nil {
	// 	return fmt.Errorf("failed to ping the database: %v", err)
	// }

	// // Exemple d'insertion dans une table fictive 'dumps'
	// _, err = db.Exec("INSERT INTO backup (dump_data) VALUES ($1)", "Sample dump data")
	// if err != nil {
	// 	return fmt.Errorf("failed to insert dump: %v", err)
	// }

	return nil
}

func DumpBdd(dbType, dbName, dbPort, userName, password string) error {

	now := time.Now()
	date := now.Format("02-01-2006_15:04:05")
	fileName := fmt.Sprintf("%s_%s.sql", dbName, date)

	// Définir le chemin de sauvegarde local
	localPath := "/app/docker_dumpFiles/" + fileName

	var dumpCmd *exec.Cmd

	switch dbType {
	case "postgres":

		// var host string
		// switch dbName {
		// case "safebase":
		// 	host = "safebase"
		// case "apple":
		// 	host = "apple"
		// case "blizzard":
		// 	host = "blizzard"
		// }

		dumpCmdStr := fmt.Sprintf(
			"PGPASSWORD=%s pg_dump -U %s -h %s -p 5432 %s",
			password,
			userName,
			dbName,
			dbName,
		)
		dumpCmd = exec.Command("/bin/sh", "-c", dumpCmdStr)

	case "mysql":
		dumpCmdStr := fmt.Sprintf(
			"mysqldump -h huawey -u %s -p%s %s",
			userName,
			password,
			dbName,
		)

		dumpCmd = exec.Command("/bin/sh", "-c", dumpCmdStr)

	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Créer le fichier localement
	outputFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local output file: %v", err)
	}
	defer outputFile.Close()

	// Rediriger la sortie de la commande vers le fichier local
	dumpCmd.Stdout = outputFile

	// Capturer la sortie d'erreur (stderr)
	var stderr bytes.Buffer
	dumpCmd.Stderr = &stderr

	// Exécuter la commande
	if err := dumpCmd.Run(); err != nil {
		return fmt.Errorf("failed to dump database: %v - %s", err, stderr.String())
	}

	fmt.Printf("Database %s dumped successfully to %s!\n", dbName, localPath)

	// // Insertion de la route de dump (peut être mise à jour selon le besoin)
	// if err := InsertDumpRoute(dbType, dbName, dbPort, userName, password); err != nil {
	// 	return fmt.Errorf("failed to insert dump route: %v", err)
	// }

	return nil
}

// func addBdd(dbType, dbName, dbPort, userName, password string) error {

// 	dumpCmdStr := fmt.Sprintf(
// 		"PGPASSWORD=%s pg_dump -U %s -h %s -p 5432 %s",
// 		password,
// 		userName,
// 		host,
// 		dbName,
// 	)
// 	dumpCmd = exec.Command("/bin/sh", "-c", dumpCmdStr)

// 	return nil
// }
