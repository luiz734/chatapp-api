package main

import (
	"fmt"
	"log"
)

// import database

func main() {
    // Create a database or open if it not exists 
	db, err := NewSqliteDB("database.db")
    if err != nil {
        panic("Can't create/open database")
    }
	defer db.Close()

	err = db.createTables()
    if err != nil {
        panic("Can't create tables")
    }

	// Insert a test user
    err = db.insertUser(&User{Nickname: "testuser"})
	if err != nil {
		log.Println("Error inserting test user:", err)
	} else {
		log.Println("Inserted test user successfully")
	}

	// Query the test user
    user, err := db.queryUser("testuser")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Retrieved user:", user.Nickname)
}
