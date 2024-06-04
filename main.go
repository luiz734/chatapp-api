package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
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

func getMessages(c *gin.Context, db *SqliteDB, roomId string, crypt Crypt) {
	messages, err := db.queryMessagesByRoom(roomId)
	if err != nil {
		panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching messages"})
		return
	}

	for i, m := range messages {
		// Encrypt the content
		encryptedContent := crypt.encrypt([]byte(m.Content))
		messages[i].Content = string(encryptedContent)

		// Encode the attachment to base64 if it exists
		if m.Attachment != nil {
			messages[i].ImageBase64 = base64.StdEncoding.EncodeToString(m.Attachment)
		}
	}

	c.IndentedJSON(http.StatusOK, messages)
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
	} else {
		panic(err)
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
	// _ = os.Args

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
		// c.Header("Content-Description", "File Transfer")
		// c.Header("Content-Transfer-Encoding", "binary")
		// c.Header("Content-Disposition", "attachment; filename=public.pem")
		// c.Header("Content-Type", "application/octet-stream")
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
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form data"})
			return
		}

		// Extract form values
		senderId := c.Request.FormValue("senderid")
		roomId := c.Request.FormValue("roomid")
		content := c.Request.FormValue("content")

		// Create the new message
		newMessage := Message{
			SenderId: senderId,
			RoomId:   roomId,
			Content:  content,
		}

		// Get the file part
		file, _, err := c.Request.FormFile("attachment")
		if err != nil {
			if err != http.ErrMissingFile { // If the file is not provided, it's okay
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving the file"})
				return
			}
		} else {
			defer file.Close()
			attachment, err := io.ReadAll(file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading the file"})
				return
			}
			compressedImg := enqueueImage(attachment,
				"helloyou.jpeg")
			// fmt.Sprint("%s%s%s", newMessage.SenderId, newMessage.RoomId))
			newMessage.Attachment = compressedImg
		}
		// Insert the new message
		addNewMessage(c, &db, newMessage)
		c.JSON(http.StatusCreated, gin.H{"message": "Message created successfully"})
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

	_ = fmt.Append
}
