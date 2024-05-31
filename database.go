package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// type Message struct {
// 	SenderId string
// 	RoomId   int
// 	Msg      string
// }
//
// type Room struct {
// 	Id   int
// 	Name string
// }
//
// type User struct {
// 	Nickname string `json:"nickname"`
// }

type Message struct {
	// id       int
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

	// createMessageTable := `CREATE TABLE IF NOT EXISTS Messages (
	// 	SenderId TEXT,
	// 	RoomId INTEGER,
	// 	Msg TEXT
	// );`
	//
	// createRoomTable := `CREATE TABLE IF NOT EXISTS Rooms (
	// 	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	Name TEXT
	// );`

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

func (sqliteDB SqliteDB) queryMessageByRoom(roomId string) (Message, error) {
	message := Message{}
	err := sqliteDB.DB.QueryRow("SELECT RoomId FROM Messages WHERE RoomId = ?", roomId).Scan(&message.RoomId)
	return message, err
}
