package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func (dbConn *DBConnection) GetPostDirectories(scope string, limit int) ([]string, error) {
	query := `SELECT path FROM posts WHERE ` + scope + ` = FALSE LIMIT ?`
	log.Printf("SQL: %s", query)
	rows, err := dbConn.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func (dbConn *DBConnection) SetPostedPostDir(postDir string, scope string) {
	insertSQL := `UPDATE posts SET ` + scope + ` = TRUE WHERE path = ?`
	log.Printf("SQL: %s", insertSQL)

	_, err := dbConn.DB.Exec(insertSQL, postDir)
	if err != nil {
		log.Printf("FAILED: to save post directory %s: %v", postDir, err)
	} else {
		log.Printf("SAVED: successfully for directory %s", postDir)
	}
}
