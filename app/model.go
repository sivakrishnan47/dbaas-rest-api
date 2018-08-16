// model.go

package app

import (
	"database/sql"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

func mgetDatabases(session *mgo.Session) ([]dataBase, error) {
	var databases []dataBase
	mongoDatabaseNames, err := session.DatabaseNames()
	if err != nil {
		return nil, err
	}

	for _, value := range mongoDatabaseNames {
		db := dataBase{}
		db.Database = value
		databases = append(databases, db)
	}
	return databases, nil
}

func (u *dataBase) mcreateDatabase(session *mgo.Session) error {
	mcollection := session.DB(u.Database).C("test")
	var mockInterface interface{}
	mcollection.Upsert(bson.M{"test": "test"}, mockInterface)
	return nil
}

func (u *dataBase) mDeleteDatabase(session *mgo.Session) error {
	var mgoDelete mgo.Database
	mgoDelete.Session = session
	mgoDelete.Name = u.Database
	err := mgoDelete.DropDatabase()
	return err
}
