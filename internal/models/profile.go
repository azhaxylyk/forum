package models

import (
	"log"
	"time"
)

func GetPostsByUser(userID string) ([]Post, error) {
	rows, err := db.Query(`
        SELECT posts.id, posts.content, posts.created_at, posts.likes, posts.dislikes, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.user_id = ?
        ORDER BY posts.created_at DESC
    `, userID)
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

		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		post.CreatedAtFormatted = createdAt.Format("02.01.2006 15:04")
		posts = append(posts, post)
	}

	return posts, nil
}

func GetLikedPostsByUser(userID string) ([]Post, error) {
	rows, err := db.Query(`
        SELECT posts.id, posts.content, posts.created_at, posts.likes, posts.dislikes, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        JOIN post_likes ON posts.id = post_likes.post_id
        WHERE post_likes.user_id = ?
        ORDER BY posts.created_at DESC
    `, userID)
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

		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		post.CreatedAtFormatted = createdAt.Format("02.01.2006 15:04")
		posts = append(posts, post)
	}

	return posts, nil
}

func GetDislikedPostsByUser(userID string) ([]Post, error) {
	rows, err := db.Query(`
        SELECT posts.id, posts.content, posts.created_at, posts.likes, posts.dislikes, users.username
        FROM posts
        JOIN users ON posts.user_id = users.id
        JOIN post_likes ON posts.id = post_likes.post_id
        WHERE post_likes.user_id = ? AND post_likes.is_like = FALSE
        ORDER BY posts.created_at DESC
    `, userID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt time.Time

		err = rows.Scan(&post.ID, &post.Content, &createdAt, &post.Likes, &post.Dislikes, &post.Author)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			log.Printf("Error fetching categories: %v", err)
			return nil, err
		}
		post.Categories = categories

		post.CreatedAtFormatted = createdAt.Format("02.01.2006 15:04")
		posts = append(posts, post)
	}

	return posts, nil
}

func GetCommentsByUser(userID string) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
        SELECT comments.id, comments.post_id, comments.content, comments.created_at, comments.likes, comments.dislikes, posts.content
        FROM comments
        JOIN posts ON comments.post_id = posts.id
        WHERE comments.user_id = ?
        ORDER BY comments.created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var comment Comment
		var postContent string
		var createdAt time.Time

		err = rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &createdAt, &comment.Likes, &comment.Dislikes, &postContent)
		if err != nil {
			return nil, err
		}

		comment.CreatedAtFormatted = createdAt.Format("02.01.2006 15:04")

		results = append(results, map[string]interface{}{
			"Comment":     comment,
			"PostContent": postContent,
		})
	}
	return results, nil
}
