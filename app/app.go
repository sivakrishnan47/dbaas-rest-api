// app.go

package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/globalsign/mgo"
	"github.com/gorilla/mux"
)

// App - web base structure
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize - DB connection and web routes
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	//a.DB, err = sql.Open("mysql", connectionString)
	//if err != nil {
	//	log.Fatal(err)
	//}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run - Starts web services
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/db", a.getDatabases).Methods("GET")
	a.Router.HandleFunc("/db", a.createDatabase).Methods("POST")
	a.Router.HandleFunc("/db/{database}", a.deleteDatabase).Methods("DELETE")
}

func (a *App) getDatabases(w http.ResponseWriter, r *http.Request) {
	products, err := getDatabases(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createDatabase(w http.ResponseWriter, r *http.Request) {
	var u []dataBase
	body,_ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err := json. Unmarshal(body,&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	for i, _ := range u {
		if err := u[i].createDatabase(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) deleteDatabase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	database, err := vars["database"]
	if !err {
		respondWithError(w, http.StatusBadRequest, "Invalid Database ID")
		return
	}

	u := dataBase{Database: database}
	if err := u.deleteDatabase(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}