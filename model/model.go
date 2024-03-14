package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Persons struct {
	gorm.Model
	PersonID  int `gorm:"primaryKey"`
	LastName  string
	FirstName string
	Address   string
	City      string
	Username  string
	Password  string
	email     string
}

type Register struct {
	Username string
	Password string
	Email    string
}

func AddPerson(db *gorm.DB, person *Persons) error {

	//fmt.Printf("%+v\n", person)
	p := db.Create(&person)
	if p.Error != nil {
		return p.Error
	}

	return nil
}

func GetPersons(db *gorm.DB) ([]Persons, error) {

	var persons []Persons
	result := db.Find(&persons)

	if result.Error != nil {
		return nil, result.Error
	}

	return persons, nil
}

func UpdatePerson(db *gorm.DB, person *Persons) (*Persons, error) {

	var p Persons
	p.ID = person.ID

	result := db.Model(&p).Where("id = ?", person.ID).Updates(person)
	if result.Error != nil {
		return nil, result.Error
	}

	return person, nil

}

func DeletePerson(db *gorm.DB, id int) error {

	var p Persons
	p.ID = uint(id)

	result := db.Delete(&p)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func Registeration(db *gorm.DB, user *Register) (Persons, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	var single_user Persons

	db.Find(&single_user, "username =?", user.Username)
	if single_user.Username != "" {
		return single_user, nil
	}

	single_user.Password = string(hash)
	single_user.Username = user.Username
	single_user.email = user.Email

	r := db.Create(&single_user)
	single_user.Password = "***"
	if r.Error != nil {
		return single_user, r.Error
	}
	single_user.Password = "***"
	return single_user, nil

}

func Login(db *gorm.DB, user *Register) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	db.Where(map[string]interface{}{"username": user.Username, "password": hash}).Find(&user)
	if user.Username != "" {
		return nil
	}

	return nil

}
