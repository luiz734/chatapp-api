package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getMessages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, []Message{
		{RoomId: "room1", SenderId: "sender1", Content: "message1"},
		{RoomId: "room2", SenderId: "sender2", Content: "message2"},
	})
}

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
	err = db.insertMessage(&Message{RoomId: "room1", SenderId: "sender1", Content: "message1"})

	if err != nil {
		log.Println("Error inserting test user:", err)
	} else {
		log.Println("Inserted test user successfully")
	}

	// Query the test user
	user, err := db.queryMessageByRoom("room1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Retrieved user:", user.Content)

	router := gin.Default()
	router.GET("/messages", getMessages)
	router.Run("localhost:55667")
}
