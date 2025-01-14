package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBConnection struct {
	DBPath string
	DB     *sql.DB
}

type DBcredentials struct {
	DbPath string
}

func (dbCred *DBcredentials) GetDbPath() string {
	return dbCred.DbPath
}

type DBconfig interface {
	GetDbPath() string
}

func NewDBConnection(DBconfig DBconfig) (*DBConnection, error) {
	db, err := initDB(DBconfig.GetDbPath())
	if err != nil {
		return nil, err
	}

	return &DBConnection{
		DBPath: DBconfig.GetDbPath(),
		DB:     db,
	}, nil
}

func (dbConn *DBConnection) Close() error {
	return dbConn.DB.Close()
}

func initDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("SELECT 1 FROM posts LIMIT 1")
	if err != nil {
		if _, err := db.Exec(DbSchemaSQL()); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func DbSchemaSQL() string {
	schema := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL UNIQUE,
		vk BOOLEAN NOT NULL DEFAULT FALSE,
		youtube BOOLEAN NOT NULL DEFAULT FALSE,
		updated BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	return schema
}
