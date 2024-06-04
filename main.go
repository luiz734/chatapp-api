package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func getMessages(c *gin.Context, db *SqliteDB, roomId string) {
	
	messages, err := db.queryMessagesByRoom(roomId)
	if err == nil {
		c.IndentedJSON(http.StatusOK, messages)
	}
}
func deleteMessage(c *gin.Context, db *SqliteDB, messageId string) {

	err := db.deleteMessage(messageId)
	if err == nil {
		c.String(http.StatusOK, "Deleted")
	}
}
func addNewMessage(c *gin.Context, db *SqliteDB, message Message) {

	err := db.insertMessage(&message)
	if err == nil {
		c.String(http.StatusCreated, "Added")
	}
}

func updateMessage(c *gin.Context, db *SqliteDB, messageId string, newContent string) {

	rowsAffected, err := db.updateMessage(messageId, newContent)
	if err == nil {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		} else {
			c.String(http.StatusOK, "Message updated successfully")
		}
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

    // Use the CORS middleware
    router.Use(CORSMiddleware())

	router.GET("/messages/:roomid", func(c *gin.Context) {
		roomId := c.Param("roomid")
		getMessages(c, &db, roomId)

	})

	router.DELETE("/delete/:messageid", func(c *gin.Context) {
		messageId := c.Param("messageid")
		deleteMessage(c, &db, messageId)

	})

	router.POST("/newMessage", func(c *gin.Context) {
		var newMessage Message
		if err := c.ShouldBindJSON(&newMessage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		addNewMessage(c, &db, newMessage)
	})

	router.PUT("/updateMessage/:id", func(c *gin.Context) {
		var newMessage Message
		id := c.Param("id")

		if err := c.ShouldBindJSON(&newMessage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		updateMessage(c, &db, id, newMessage.Content)
	})
	router.Run("0.0.0.0:55667")
}
