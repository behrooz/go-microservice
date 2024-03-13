package model

import (
	"fmt"
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

func Register(db *gorm.DB, username, password, email string) error {
	var u Persons
	u.Username = username
	u.Password = password
	u.email = email

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Printf("hash: %v\n", hash)
	if err != nil {
		log.Fatal(err)
	}

	var single_user Persons

	db.Find(&single_user, "username =?", username)

	fmt.Printf("user: %v\n", single_user)

	return nil

}
