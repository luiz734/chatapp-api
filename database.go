package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	Id          int    `json:"id,omitempty"`
	SenderId    string `json:"senderid"`
	RoomId      string `json:"roomid"`
	Content     string `json:"content"`
	Attachment  []byte `json:"-"`
	ImageBase64 string `json:"imageBase64,omitempty"`
	Signature   []byte `json:"signature,omitempty"`
}

type SqliteDB struct {
	filepath string
	DB       *sql.DB
}

func NewSqliteDB(filepath string) (SqliteDB, error) {
	db, err := sql.Open("sqlite3", filepath)
	return SqliteDB{filepath: filepath, DB: db}, err
}

func (sqliteDB SqliteDB) Close() {
	sqliteDB.DB.Close()
}

func (sqliteDB SqliteDB) createTables() error {
	createMessageTable := `CREATE TABLE IF NOT EXISTS Messages (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
        SenderId TEXT,
        RoomId TEXT,
        Content TEXT,
        Attachment BLOB
	);`

	_, err := sqliteDB.DB.Exec(createMessageTable)
	return err
}

func (sqliteDB SqliteDB) insertMessage(message *Message) error {
	var err error
	if message.Attachment == nil {
		_, err = sqliteDB.DB.Exec(
			"INSERT INTO Messages (SenderId, RoomId, Content) VALUES (?, ?, ?)",
			message.SenderId, message.RoomId, message.Content)
	} else {
		_, err = sqliteDB.DB.Exec(
			"INSERT INTO Messages (SenderId, RoomId, Content, Attachment) VALUES (?, ?, ?, ?)",
			message.SenderId, message.RoomId, message.Content, message.Attachment)
	}
	return err
}

func (sqliteDB SqliteDB) queryMessagesByRoom(roomId string) ([]Message, error) {
	rows, err := sqliteDB.DB.Query("SELECT * FROM Messages WHERE RoomId = ?", roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		// var msg Message
		msg := Message{}
		var attachment sql.RawBytes

		if err := rows.Scan(&msg.Id, &msg.SenderId, &msg.RoomId, &msg.Content, &attachment); err != nil {
			return messages, err
		}
		if attachment != nil {
			msg.Attachment = []byte(attachment)
		}

		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return messages, err
	}
	return messages, nil
}

func (sqliteDB SqliteDB) deleteMessage(messageId string) error {
	_, err := sqliteDB.DB.Exec("DELETE FROM messages WHERE Id = ?", messageId)
	if err != nil {
		panic(err)
	}
	return err
}

func (sqliteDB SqliteDB) updateMessage(messageId string, newContent string) (int64, error) {
	stmt, err := sqliteDB.DB.Prepare("UPDATE Messages SET Content = ? WHERE Id = ?")
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Database preparation error"})
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(newContent, messageId)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Database execution error"})
		return 0, err
	}

	return res.RowsAffected()
}
