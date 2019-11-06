package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"os"
)

var routeCache map[string]*IFile

const (
	DBLocation = "./db.json"
	DefaultUsername = "admin"
	DefaultPassword = "password"
)

type IDatabase struct {
	Users []*IUser
	Files []*IFile
}

type IUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type IFile struct {
	FileLocation string `json:"file_location"`
	DownloadOnLoad bool `json:"download_on_load"`
	WebRoute string `json:"web_route"`
	DateCreated int64 `json:"date_created"`
}

func StoreFile(file *IFile) {
	routeCache[file.WebRoute] = file
	db.Files = append(db.Files, file)
	StoreDB()
}

func LoadDB() {

	if _, err := os.Stat(DBLocation); os.IsNotExist(err) {
		StoreDB()
	}

	bV, err := ioutil.ReadFile(DBLocation)
	if err != nil {
		log.Fatalf("Cannot load database: %s\n", err)
	}
	err = json.Unmarshal(bV, &db)
	if err != nil {
		log.Fatalf("Error unmarshalling db from json: %s\n", err)
	}

	// build cache
	for _, file := range db.Files {
		routeCache[file.WebRoute] = file
	}

	// add admin user if there isn't one
	if len(db.Users) == 0 {
		pass, err := bcrypt.GenerateFromPassword([]byte(DefaultPassword), 10)
		if err != nil {
			log.Fatal(err)
		}
		db.Users = append(db.Users, &IUser{DefaultUsername, string(pass)})

		StoreDB()
	}
}

func StoreDB() {
	if _, err := os.Stat(DBLocation); !os.IsNotExist(err) {
		err := os.Rename(DBLocation, DBLocation+".backup")
		if err != nil {
			log.Fatalf("Cannot create backup: %s\n", err)
		}
	}
	b, err := json.Marshal(db)
	if err != nil {
		log.Printf("Cannot marshal db to JSON: %s\n", err)
		return
	}
	err = ioutil.WriteFile(DBLocation, b, 0644)
	if err != nil {
		log.Fatalf("Cannot write DB to file %s\n", err)
	}
	err = os.Remove(DBLocation + ".backup")
	if err != nil {
		log.Printf("Cannot delete backup: %s\n", err)
	}
}