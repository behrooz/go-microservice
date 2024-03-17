package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"main/model"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

var secretKey = []byte("eyforShipApplication!@#@!")

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

	db.AutoMigrate(&model.Persons{})
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
	a.Router.HandleFunc("/person/{id}", a.deletePerson).Methods("DELETE")

	a.Router.HandleFunc("/register", a.register).Methods("POST")
	a.Router.HandleFunc("/login", a.login).Methods("POST")
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

	var p model.Persons
	err := json.NewDecoder(r.Body).Decode(&p)
	model.AddPerson(a.DB, &p)
	if err != nil {
		ResponseWithJson(w, 500, map[string]string{"error": err.Error()})
	}
}

func (a *App) getPersons(w http.ResponseWriter, r *http.Request) {
	result, err := model.GetPersons(a.DB)
	if err != nil {
		ResponseWithError(w, 500, err.Error())
	}

	ResponseWithJson(w, 200, result)
}

func (a *App) updatePerson(w http.ResponseWriter, r *http.Request) {

	var p model.Persons
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		ResponseWithError(w, 500, err.Error())
	}

	result, err := model.UpdatePerson(a.DB, &p)
	if err != nil {
		ResponseWithError(w, 500, err.Error())
	}

	ResponseWithJson(w, 200, result)

}

func (a *App) deletePerson(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ResponseWithError(w, 400, "Invalid ID")
		return
	}
	err = model.DeletePerson(a.DB, idInt)
	if err != nil {
		ResponseWithError(w, 500, err.Error())
		return
	}
	ResponseWithJson(w, 200, "ok")
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {

	var u model.Register

	_ = json.NewDecoder(r.Body).Decode(&u)

	_, err := model.Registeration(a.DB, &u)

	if err != nil {
		ResponseWithError(w, 500, err.Error())
		return
	}

	ResponseWithJson(w, 200, "ok")
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var u model.Register

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		ResponseWithError(w, 500, "Missing Credentials")
		return
	}

	if u.Username == "" || u.Password == "" {
		ResponseWithError(w, 500, "Invalid Credentials")
		return
	}

	result, err := model.Login(a.DB, &u)

	if err != nil {
		ResponseWithError(w, 500, err.Error())
		return
	}

	if result.Username != "" {
		result, _ := createToken(u.Username)

		item := model.LoginToken{
			Username: u.Username,
			Token:    result,
		}

		ResponseWithJson(w, 200, &item)
	}

	if err != nil {
		ResponseWithError(w, 500, err.Error())
	}
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
