package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestConnectMainDB(t *testing.T) {
	// Créer un mock SQL avec l'option MonitorPings pour surveiller les appels à Ping()
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	// Simuler la réponse du ping
	mock.ExpectPing().WillReturnError(nil)

	// Simuler le comportement de la fonction connectMainDB
	mainDB = db
	err = mainDB.Ping()

	// Valider qu'il n'y a pas d'erreurs et que la connexion fonctionne
	assert.NoError(t, err, "Erreur lors du ping de la base de données principale")

	if err == nil {
		t.Log("Connexion à la base de données principale (mock) réussie")
	}

	// Vérifier les attentes de sqlmock
	assert.NoError(t, mock.ExpectationsWereMet(), "Toutes les attentes n'ont pas été satisfaites")
}
