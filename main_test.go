// main_test.go

package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("root", "", "dbaas")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableDBCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(tableDBTypeCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM dbs")
	a.DB.Exec("DELETE FROM db_type")
	a.DB.Exec("ALTER TABLE dbs AUTO_INCREMENT = 1")
	a.DB.Exec("ALTER TABLE db_type AUTO_INCREMENT = 1")
	a.DB.Exec("INSERT INTO db_type(name)VALUES('MYSQL')")

}

const tableDBCreationQuery = `
CREATE TABLE IF NOT EXISTS db_type (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(50) NOT NULL,
  PRIMARY KEY (id)
)`

const tableDBTypeCreationQuery =  `
CREATE TABLE IF NOT EXISTS dbs (
  id int NOT NULL AUTO_INCREMENT,
  typeID int NOT NULL,
  name varchar(50) NOT NULL,
  ip varchar(15) NOT NULL,
  dbPort int NOT NULL,
  createdData timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY typeID (typeID),
  CONSTRAINT dbs_fk_dbType FOREIGN KEY (typeID) REFERENCES db_type (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/db", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateDatabase(t *testing.T) {
	clearTable()

	payload := []byte(`{"TypeID":1,"Name":"TEST_DB","IP":"12345","Port":12345}`)

	req, _ := http.NewRequest("POST", "/db", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "TEST_DB" {
		t.Errorf("Expected user name to be 'TEST_DB'. Got '%v'", m["name"])
	}

	if m["ip"] != "12345" {
		t.Errorf("Expected ip to be '12345'. Got '%v'", m["ip"])
	}

	if m["Port"] != 12345.0 {
		t.Errorf("Expected user port to be '12345'. Got '%v'", m["Port"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func addDatabase(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO dbs(typeID,name,ip,dbPort)VALUES(%d,'%s','%s', %d)", 1, ("DB" + strconv.Itoa(i)), ("IP" + strconv.Itoa(i)), (8000 + i))
		a.DB.Exec(statement)
	}
}

func TestGetDatabase(t *testing.T) {
	clearTable()
	addDatabase(1)

	req, _ := http.NewRequest("GET", "/db", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteDatabase(t *testing.T) {
	clearTable()
	addDatabase(1)

	req, _ := http.NewRequest("DELETE", "/db/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}