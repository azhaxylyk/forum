package handlers

import (
	"database/sql"
	"forum/internal/models"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Неподдерживаемый метод запроса")
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		log.Println("Пользователь не аутентифицирован")
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		log.Printf("Ошибка получения ID пользователя: %v", err)
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	content := r.FormValue("content")
	categories := r.Form["categories"]

	// logs
	log.Printf("Получен контент: %q", content)
	log.Printf("Получены категории: %v", categories)

	content = models.SanitizeInput(content)
	if !models.IsValidContent(content) || len(categories) == 0 {
		log.Println("Некорректный контент или отсутствуют категории")
		ErrorHandler(w, r, http.StatusBadRequest, "Content and at least one category are required to create a post")
		return
	}

	// image-upload
	var imagePath string
	file, header, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("Ошибка при загрузке изображения: %v", err)
		ErrorHandler(w, r, http.StatusBadRequest, "Invalid image upload")
		return
	}
	if err == nil {
		defer file.Close()

		if header.Size > 20*1024*1024 {
			log.Println("Размер изображения превышает 20MB")
			ErrorHandler(w, r, http.StatusBadRequest, "Image size exceeds 20MB limit")
			return
		}

		fileType := header.Header.Get("Content-Type")
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
		}
		if !allowedTypes[fileType] {
			log.Printf("Неподдерживаемый тип изображения: %s", fileType)
			ErrorHandler(w, r, http.StatusBadRequest, "Unsupported image type. Allowed types: JPEG, PNG, GIF")
			return
		}

		// Генерация уникального имени файла
		fileExtension := filepath.Ext(header.Filename)
		uniqueFileName, err := uuid.NewV4() // Использование NewV4 из github.com/gofrs/uuid
		if err != nil {
			log.Printf("Ошибка при генерации UUID: %v", err)
			ErrorHandler(w, r, http.StatusInternalServerError, "Unable to generate unique file name")
			return
		}
		uploadPath := "./web/static/uploads/" + uniqueFileName.String() + fileExtension

		// create directory
		err = os.MkdirAll(filepath.Dir(uploadPath), os.ModePerm)
		if err != nil {
			log.Printf("Ошибка при создании директории для загрузок: %v", err)
			ErrorHandler(w, r, http.StatusInternalServerError, "Unable to create directory for uploads")
			return
		}

		out, err := os.Create(uploadPath)
		if err != nil {
			log.Printf("Ошибка при сохранении изображения: %v", err)
			ErrorHandler(w, r, http.StatusInternalServerError, "Unable to save the image")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			log.Printf("Ошибка при копировании изображения: %v", err)
			ErrorHandler(w, r, http.StatusInternalServerError, "Error while saving the image")
			return
		}

		// way
		imagePath = "/static/uploads/" + uniqueFileName.String() + fileExtension
	}

	// Создание поста
	postID, err := models.CreatePost(userID, content, imagePath)
	if err != nil {
		log.Printf("Ошибка при создании поста в базе данных: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError, "Error creating post")
		return
	}

	for _, categoryID := range categories {
		err = models.AddCategoryToPost(postID, categoryID)
		if err != nil {
			log.Printf("Ошибка при ассоциации категории %s с постом %s: %v", categoryID, postID, err)
			ErrorHandler(w, r, http.StatusInternalServerError, "Error associating category")
			return
		}
	}

	log.Printf("Пост %s успешно создан пользователем %s", postID, userID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postID := r.FormValue("post_id")

	err = models.LikePost(userID, postID)
	if err != nil {
		http.Error(w, "Error liking post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = models.UpdatePostLikesDislikes(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating like count")
		return
	}

	postOwnerID, err := models.GetPostOwner(postID)
	if err == nil && postOwnerID != userID {
		_ = models.CreateNotification(postOwnerID, userID, "like", postID, "post")
	}

	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postID := r.FormValue("post_id")

	err = models.DislikePost(userID, postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error disliking post")
		return
	}

	err = models.UpdatePostLikesDislikes(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating dislike count")
		return
	}

	postOwnerID, err := models.GetPostOwner(postID)
	if err == nil && postOwnerID != userID {
		_ = models.CreateNotification(postOwnerID, userID, "dislike", postID, "post")
	}

	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	postID := r.URL.Query().Get("id")
	if postID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing post ID")
		return
	}

	post, err := models.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorHandler(w, r, http.StatusNotFound, "Post not found")
			return
		}
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching post")
		return
	}

	comments, err := models.GetCommentsForPost(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching comments")
		return
	}

	notification := r.URL.Query().Get("notification")

	tmpl, err := template.ParseFiles("web/templates/comments.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	var loggedIn bool
	var isModerator bool
	var username string

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie == nil || cookie.Value == "" {
		loggedIn = false
	} else {
		userID, username, err := models.GetIDBySessionToken(cookie.Value)
		if err != nil || userID == "" || username == "" {
			loggedIn = false
		} else {
			loggedIn = true
			isModerator = IsModerator(userID)
		}
	}

	data := struct {
		Post         models.Post
		Comments     []models.Comment
		LoggedIn     bool
		Username     string
		IsModerator  bool
		Notification string
	}{
		Post:         post,
		Comments:     comments,
		LoggedIn:     loggedIn,
		Username:     username,
		IsModerator:  isModerator,
		Notification: notification,
	}
	tmpl.Execute(w, data)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	postID := r.FormValue("post_id")
	ownerID, err := models.GetPostOwner(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching post owner")
		return
	}

	if ownerID != userID {
		ErrorHandler(w, r, http.StatusForbidden, "You are not allowed to delete this post")
		return
	}

	err = models.DeletePost(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error deleting post")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	postID := r.FormValue("post_id")
	newContent := r.FormValue("content")

	ownerID, err := models.GetPostOwner(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching post owner")
		return
	}

	if ownerID != userID {
		ErrorHandler(w, r, http.StatusForbidden, "You are not allowed to edit this post")
		return
	}

	// Обработка загрузки нового изображения
	var newImagePath string
	file, header, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		ErrorHandler(w, r, http.StatusBadRequest, "Invalid image upload")
		return
	}

	if err == nil {
		defer file.Close()

		// Проверка размера изображения
		if header.Size > 20*1024*1024 {
			ErrorHandler(w, r, http.StatusBadRequest, "Image size exceeds 20MB limit")
			return
		}

		// Проверка типа изображения
		fileType := header.Header.Get("Content-Type")
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
		}
		if !allowedTypes[fileType] {
			ErrorHandler(w, r, http.StatusBadRequest, "Unsupported image type. Allowed types: JPEG, PNG, GIF")
			return
		}

		// Генерация уникального имени файла
		fileExtension := filepath.Ext(header.Filename)
		uniqueFileName, err := uuid.NewV4()
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Unable to generate unique file name")
			return
		}
		uploadPath := "./web/static/uploads/" + uniqueFileName.String() + fileExtension

		// Создание директории, если она не существует
		err = os.MkdirAll(filepath.Dir(uploadPath), os.ModePerm)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Unable to create directory for uploads")
			return
		}

		// Сохранение нового изображения
		out, err := os.Create(uploadPath)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Unable to save the image")
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error while saving the image")
			return
		}

		// Установка нового пути к изображению
		newImagePath = "/static/uploads/" + uniqueFileName.String() + fileExtension

		// Удаление старого изображения, если оно существует
		oldPost, err := models.GetPostByID(postID)
		if err == nil && oldPost.ImagePath != "" {
			oldImagePath := "." + oldPost.ImagePath // Преобразуем путь для удаления
			if _, err := os.Stat(oldImagePath); err == nil {
				os.Remove(oldImagePath)
			}
		}
	}

	// Обновление поста в базе данных
	err = models.UpdatePost(postID, newContent, newImagePath)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating post")
		return
	}

	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}
