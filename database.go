package main

import (
	"database/sql"
	// "fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	Id       int
	SenderId string
	RoomId   string
	Content  string
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
        Content TEXT
	);`

	_, err := sqliteDB.DB.Exec(createMessageTable)
	return err
}

func (sqliteDB SqliteDB) insertMessage(message *Message) error {
	_, err := sqliteDB.DB.Exec("INSERT INTO Messages (SenderId, RoomId, Content) VALUES (?, ?, ?)",
		message.SenderId, message.RoomId, message.Content)
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
		if err := rows.Scan(&msg.Id, &msg.SenderId, &msg.RoomId, &msg.Content); err != nil {
			return messages, err
		} 
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return messages, err
	}
	return messages, nil
}

// func (sqliteDB SqliteDB) queryMessageByRoom(roomId string) (Message, error) {
// 	message := Message{}
// 	err := sqliteDB.DB.QueryRow("SELECT * FROM Messages WHERE RoomId = ?", roomId).
// 		Scan(&message.Id, &message.SenderId, &message.RoomId, &message.Content)
// 	return message, err
// }
