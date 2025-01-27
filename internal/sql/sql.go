package sql

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	err = CreateTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	sqlFiles := []string{
		"./migrations/users_table.sql",
		"./migrations/posts_table.sql",
		"./migrations/categories_table.sql",
		"./migrations/comments_table.sql",
		"./migrations/post_likes_table.sql",
		"./migrations/comment_likes_table.sql",
		"./migrations/post_categories_table.sql",
		"./migrations/notifications_table.sql",
		"./migrations/moderation_requests.sql",
	}

	for _, file := range sqlFiles {
		query, err := LoadSQLFile(file)
		if err != nil {
			return err
		}

		_, err = db.Exec(query)
		if err != nil {
			log.Fatalf("Ошибка выполнения SQL из файла %s: %v", file, err)
		}
	}

	err := SeedCategories(db)
	if err != nil {
		return err
	}

	log.Println("Все таблицы успешно созданы.")
	return nil
}

func SeedCategories(db *sql.DB) error {
	query, err := LoadSQLFile("./migrations/seed_categories.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Ошибка выполнения seed_categories.sql: %v", err)
		return err
	}

	log.Println("Категории успешно добавлены.")
	return nil
}

func LoadSQLFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
