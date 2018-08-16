// app.go

package app

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/globalsign/mgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// App - web base structure
type App struct {
	Router   *mux.Router
	DB       *sql.DB
	msession *mgo.Session
}

// Initialize - DB connection and web routes
func (a *App) Initialize(user, password, dbname string, mongoHost string, mongoTimeout time.Duration) {
	// connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err, mErr error
	//Initializing sql connection
	// a.DB, err = sql.Open("mysql", connectionString)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//Initializing mongo connection
	log.Println("mongo init")
	a.msession, mErr = newSession(mongoHost, mongoTimeout)
	if mErr != nil {
		log.Fatal(err)
	}
	log.Println("after mongo initialize")
	//a.DB, err = sql.Open("mysql", connectionString)
	//if err != nil {
	//	log.Fatal(err)
	//}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func newSession(url string, timeout time.Duration) (*mgo.Session, error) {
	// Set the default timeout for the session.
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	ses, err := mgo.DialWithTimeout(url, timeout)
	if err != nil {
		return nil, err
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	ses.SetMode(mgo.Monotonic, true)

	return ses, nil
}

// Run - Starts web services
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/db", a.getDatabases).Methods("GET")
	a.Router.HandleFunc("/db", a.createDatabase).Methods("POST")
	a.Router.HandleFunc("/db/{database}", a.deleteDatabase).Methods("DELETE")
	a.Router.HandleFunc("/mdb", a.MgetDatabases).Methods("GET")
	a.Router.HandleFunc("/mdb", a.McreateDatabase).Methods("POST")
	a.Router.HandleFunc("/mdb/{database}", a.MdeleteDatabase).Methods("DELETE")
}

func (a *App) getDatabases(w http.ResponseWriter, r *http.Request) {
	products, err := getDatabases(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) MgetDatabases(w http.ResponseWriter, r *http.Request) {
	databaseList, err := mgetDatabases(a.msession)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, databaseList)
}

func (a *App) McreateDatabase(w http.ResponseWriter, r *http.Request) {
	var u []dataBase
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err := json.Unmarshal(body, &u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	for index := range u {
		if err := u[index].mcreateDatabase(a.msession); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) MdeleteDatabase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	database, err := vars["database"]
	if !err {
		respondWithError(w, http.StatusBadRequest, "Invalid Database ID")
		return
	}

	u := dataBase{Database: database}
	if err := u.mDeleteDatabase(a.msession); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) createDatabase(w http.ResponseWriter, r *http.Request) {
	var u []dataBase
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err := json.Unmarshal(body, &u); err != nil {
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
