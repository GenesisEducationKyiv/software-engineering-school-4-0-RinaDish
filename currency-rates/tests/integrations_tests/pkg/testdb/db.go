package testdb

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	_ "github.com/lib/pq"
)

const (
	DBName         = "testdb"
	TemplateDBName = "api_template"
	dsnTemplate    = "postgres://test_user:secret@127.0.0.1:8632/%s?sslmode=disable"
)

func GetDBDSN() string {
	return fmt.Sprintf(dsnTemplate, DBName)
}

func GetTemplateDBDSN() string {
	return fmt.Sprintf(dsnTemplate, TemplateDBName)
}

func GetEmail(t *testing.T, dsn, dbName, email string) string {
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	var res string
	
	row := db.QueryRow(fmt.Sprintf(`SELECT emails.email FROM emails where emails.email = '%s';`, email))

	_ = row.Scan(&res)

	return res
}

func ResetDBTemplate(t *testing.T, dsn, dbName, tplDBName string) {
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %q`, dbName))
	require.NoError(t, err)

	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE %q TEMPLATE %q`, dbName, tplDBName))
	require.NoError(t, err)
}
