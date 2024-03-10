package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(username, password, dbname string) {

	dns := fmt.Sprintf("%s:%s@tcp(localhost:3306)/?charset=utf8mb4&parseTime=True&loc=Local", username, password)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	createDatabaseCommand := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname)
	db.Exec(createDatabaseCommand)

	dns = fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbname)
	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{})

	db.AutoMigrate(&Persons{})
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db
	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) {

	fmt.Print("application is running on :8010")
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/person", a.getPersons).Methods("GET")
	a.Router.HandleFunc("/person", a.addPerson).Methods("POST")
	a.Router.HandleFunc("/person", a.updatePerson).Methods("PATCH")
}

func ResponseWithError(w http.ResponseWriter, code int, message string) {
	ResponseWithJson(w, code, map[string]string{"error": message})
}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	respose, _ := json.Marshal(payload)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respose)
}

func (a *App) addPerson(w http.ResponseWriter, r *http.Request) {

	var p Persons
	err := json.NewDecoder(r.Body).Decode(&p)
	addPerson(a.DB, &p)
	if err != nil {
		ResponseWithJson(w, 500, map[string]string{"error": err.Error()})
	}
}

func (a *App) getPersons(w http.ResponseWriter, r *http.Request) {
	result, err := getPersons(a.DB)
	if err != nil {
		ResponseWithError(w, 500, err.Error())
	}

	ResponseWithJson(w, 200, result)
}

func (a *App) updatePerson(w http.ResponseWriter, r *http.Request) {

	var p Persons
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		ResponseWithError(w, 500, err.Error())
	}

	result, err := updatePerson(a.DB, &p)
	if err != nil {
		ResponseWithError(w, 500, err.Error())
	}

	ResponseWithJson(w, 200, result)

}
