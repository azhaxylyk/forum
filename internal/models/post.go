package models

import (
	"database/sql"
	"html/template"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	ID                 string
	Content            template.HTML
	CreatedAt          time.Time
	CreatedAtFormatted string
	Likes              int
	Dislikes           int
	Author             string
	LoggedIn           bool
	UserHasLiked       bool
	UserHasDisliked    bool
	Categories         []string
}

type Category struct {
	ID   string
	Name string
}

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func CreatePost(userID, content string) (string, error) {
	postID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	_, err = db.Exec("INSERT INTO posts (id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		postID.String(), userID, content, time.Now())
	return postID.String(), err
}

func AddCategoryToPost(postID, categoryID string) error {
	_, err := db.Exec(`
        INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)
    `, postID, categoryID)
	return err
}

func GetCategoriesForPost(postID string) ([]string, error) {
	rows, err := db.Query(`
        SELECT categories.name 
        FROM categories
        JOIN post_categories ON categories.id = post_categories.category_id
        WHERE post_categories.post_id = ?
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func LikePost(userID, postID string) error {
	var interactionID string
	var isLike bool

	err := db.QueryRow("SELECT id, is_like FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&interactionID, &isLike)
	if err == sql.ErrNoRows {
		likeID, _ := uuid.NewV4()
		_, err = db.Exec("INSERT INTO post_likes (id, user_id, post_id, is_like) VALUES (?, ?, ?, TRUE)", likeID.String(), userID, postID)
		return err
	} else if err != nil {
		return err
	}

	if isLike {
		_, err = db.Exec("DELETE FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID)
		return err
	}

	_, err = db.Exec("UPDATE post_likes SET is_like = TRUE WHERE id = ?", interactionID)
	return err
}

func DislikePost(userID, postID string) error {
	var interactionID string
	var isLike bool

	// Получаем текущую запись для пользователя и поста
	err := db.QueryRow("SELECT id, is_like FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&interactionID, &isLike)
	if err == sql.ErrNoRows {
		// Если записи нет, создаем новую с is_like = FALSE (дизлайк)
		dislikeID, _ := uuid.NewV4()
		_, err = db.Exec("INSERT INTO post_likes (id, user_id, post_id, is_like) VALUES (?, ?, ?, FALSE)", dislikeID.String(), userID, postID)
		return err
	} else if err != nil {
		return err
	}

	// Если это дизлайк, удаляем запись
	if !isLike {
		_, err = db.Exec("DELETE FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID)
		return err
	}

	// Если это лайк, обновляем на дизлайк
	_, err = db.Exec("UPDATE post_likes SET is_like = FALSE WHERE id = ?", interactionID)
	return err
}

func UpdatePostLikesDislikes(postID string) error {
	var likeCount, dislikeCount int

	err := db.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND is_like = TRUE", postID).Scan(&likeCount)
	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND is_like = FALSE", postID).Scan(&dislikeCount)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE posts SET likes = ?, dislikes = ? WHERE id = ?", likeCount, dislikeCount, postID)
	return err
}

func GetFilteredPosts(loggedIn bool, userID, categoryID string) ([]Post, error) {
	var rows *sql.Rows
	var err error

	if categoryID != "" {
		rows, err = db.Query(`
            SELECT posts.id, posts.content, posts.created_at, posts.likes, posts.dislikes, users.username
            FROM posts
            JOIN users ON posts.user_id = users.id
            JOIN post_categories ON posts.id = post_categories.post_id
            WHERE post_categories.category_id = ?
            ORDER BY posts.created_at DESC
        `, categoryID)
	} else {
		rows, err = db.Query(`
            SELECT posts.id, posts.content, posts.created_at, posts.likes, posts.dislikes, users.username
            FROM posts
            JOIN users ON posts.user_id = users.id
            ORDER BY posts.created_at DESC
        `)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt time.Time

		err = rows.Scan(&post.ID, &post.Content, &createdAt, &post.Likes, &post.Dislikes, &post.Author)
		if err != nil {
			return nil, err
		}

		post.CreatedAtFormatted = createdAt.Format("02.01.2006 15:04")
		post.Content = template.HTML(strings.ReplaceAll(string(post.Content), "\n", "<br>"))
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}

	return posts, nil
}

func GetAllCategories() ([]Category, error) {
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetPostByID(postID string) (Post, error) {
	var post Post
	var createdAt time.Time

	err := db.QueryRow(`
        SELECT posts.id, posts.content, posts.created_at, posts.likes, posts.dislikes, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.id = ?`, postID).Scan(
		&post.ID, &post.Content, &createdAt, &post.Likes, &post.Dislikes, &post.Author,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return post, err
		}
		return post, err
	}

	categories, err := GetCategoriesForPost(post.ID)
	if err != nil {
		return post, err
	}
	post.Categories = categories

	post.CreatedAtFormatted = createdAt.Format("02.01.2006 15:04")

	return post, nil
}

func GetPostOwner(postID string) (string, error) {
	var ownerID string
	err := db.QueryRow(`
        SELECT users.id 
        FROM users
        JOIN posts ON users.id = posts.user_id
        WHERE posts.id = ?
    `, postID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Пост не найден
		}
		return "", err
	}
	return ownerID, nil
}
