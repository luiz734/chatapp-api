package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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


func getMessages(c *gin.Context, db *SqliteDB, roomId string, crypt Crypt) {
	messages, err := db.queryMessagesByRoom(roomId)
    for i, m := range messages {
        messages[i].Content = string(crypt.encrypt([]byte(m.Content)))
    }
    fmt.Println(messages)
	if err == nil {
		c.IndentedJSON(http.StatusOK, messages)
	} else {
        panic(err)
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	var crypt Crypt
	if !KeysExists() {
		CreateKeyPair()
	}
	crypt.loadKeys()
	//    a := crypt.encrypt([]byte("hello"))
	//    b := crypt.decrypt(a)
	//    fmt.Println(string(a))
	//    _=b

	// image, err := os.ReadFile("plane.jpeg")
	// if err != nil {
	// 	panic("Can't read file")
	// }
	// _ = image
	// _ = fmt.Println
	// for i := 0; i < 1; i++ {
	// 	enqueueImage(image, "plane.jpeg")
	// }
	// return
    _=os.Args

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

	router := gin.Default()

	// Use the CORS middleware
	router.Use(CORSMiddleware())

	router.GET("/key", func(c *gin.Context) {
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename=public.pem")
		c.Header("Content-Type", "application/octet-stream")
		c.Writer.Write(crypt.PublicKeyPem)
	})

	router.GET("/messages/:roomid", func(c *gin.Context) {
		roomId := c.Param("roomid")
		getMessages(c, &db, roomId, crypt)

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
