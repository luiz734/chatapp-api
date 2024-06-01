package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getMessages(c *gin.Context, db *SqliteDB, roomId string) {
	messages, err := db.queryMessagesByRoom(roomId)
	if err == nil {
		c.IndentedJSON(http.StatusOK, messages)
	}
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
	// err = db.insertMessage(&Message{RoomId: "room1", SenderId: "sender1", Content: "message1"})

	// if err != nil {
	// 	log.Println("Error inserting test user:", err)
	// } else {
	// 	log.Println("Inserted test user successfully")
	// }

	// Query the test user
	// user, err := db.queryMessageByRoom("room1")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Retrieved user:", user.Content)

	router := gin.Default()
    router.GET("/messages/:roomid", func(c *gin.Context) {
        roomId := c.Param("roomid")
		getMessages(c, &db, roomId)

	})
	router.Run("0.0.0.0:55667")
}
