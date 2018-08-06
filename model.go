// model.go

package main

import (
	"fmt"
	"database/sql"
)

type dataBase struct {
	ID   int    `json:"id"`
	TypeID int  `json:"typeID"`
	Name string `json:"name"`
	IP  string  `json:"ip"`
	Port int	`json:port`
	CreatedDate string `json:createdData`
}

type dataBaseType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u *dataBase) deleteDatabase(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM dbs WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *dataBase) createDatabase(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO dbs (typeID, name, ip, dbPort) VALUES (%d, '%s', '%s', %d)", u.TypeID, u.Name, u.IP, u.Port)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getDatabases(db *sql.DB, start, count int) ([]dataBase, error) {
	statement := fmt.Sprintf("SELECT id, typeID, name, ip, dbPort, createdData FROM dbs LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	databases := []dataBase{}

	for rows.Next() {
		var u dataBase
		if err := rows.Scan(&u.ID, &u.TypeID, &u.Name, &u.IP, &u.Port, &u.CreatedDate); err != nil {
			return nil, err
		}
		databases = append(databases, u)
	}

	return databases, nil
}

func getDatabaseTypes(db *sql.DB, start, count int) ([]dataBaseType, error) {
	statement := fmt.Sprintf("SELECT ID, NAME FROM db_type LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	dataBaseTypes := []dataBaseType{}

	for rows.Next() {
		var u dataBaseType
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		dataBaseTypes = append(dataBaseTypes, u)
	}

	return dataBaseTypes, nil
}