// model.go

package app

import (
	"fmt"
	"database/sql"
)

type dataBase struct {
	Database string `json:database`
}

func (u *dataBase) deleteDatabase(db *sql.DB) error {
	statement := fmt.Sprintf("DROP DATABASE %s", u.Database)
	fmt.Println(statement)
	_, err := db.Exec(statement)
	return err
}

func (u *dataBase) createDatabase(db *sql.DB) error {
	statement := fmt.Sprintf("CREATE DATABASE %s", u.Database)
	fmt.Println(statement)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	return nil
}

func getDatabases(db *sql.DB) ([]dataBase, error) {
	statement := fmt.Sprintf("SHOW DATABASES")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	databases := []dataBase{}

	for rows.Next() {
		var u dataBase

		if err := rows.Scan(&u.Database); err != nil {
			return nil, err
		}
		databases = append(databases, u)
	}

	return databases, nil
}