package main

import (
	"gorm.io/gorm"
)

type Persons struct {
	gorm.Model
	PersonID  int `gorm:"primaryKey"`
	LastName  string
	FirstName string
	Address   string
	City      string
}

func addPerson(db *gorm.DB, person *Persons) error {

	//fmt.Printf("%+v\n", person)
	p := db.Create(&person)
	if p.Error != nil {
		return p.Error
	}

	return nil
}

func getPersons(db *gorm.DB) ([]Persons, error) {

	var persons []Persons
	result := db.Find(&persons)

	if result.Error != nil {
		return nil, result.Error
	}

	return persons, nil
}
