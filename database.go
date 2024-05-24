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

type User struct {
	Nickname string
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

	createUserTable := `CREATE TABLE IF NOT EXISTS Users (
		Nickname TEXT PRIMARY KEY
	);`

	// _, err := db.Exec(createMessageTable)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = db.Exec(createRoomTable)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_, err := sqliteDB.DB.Exec(createUserTable)
	return err
}

func (sqliteDB SqliteDB) insertUser(user *User) error {
	_, err := sqliteDB.DB.Exec("INSERT INTO Users (Nickname) VALUES (?)", user.Nickname)
	return err
}
func (sqliteDB SqliteDB) queryUser(userId string) (User, error) {
	user := User{}
	err := sqliteDB.DB.QueryRow("SELECT Nickname FROM Users WHERE Nickname = ?", userId).Scan(&user.Nickname)
	return user, err
}
