package models

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestSetDB(t *testing.T) {
	type args struct {
		database *sql.DB
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDB(tt.args.database)
		})
	}
}

func TestCreatePost(t *testing.T) {
	type args struct {
		userID    string
		content   string
		imagePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreatePost(tt.args.userID, tt.args.content, tt.args.imagePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreatePost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddCategoryToPost(t *testing.T) {
	type args struct {
		postID     string
		categoryID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddCategoryToPost(tt.args.postID, tt.args.categoryID); (err != nil) != tt.wantErr {
				t.Errorf("AddCategoryToPost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCategoriesForPost(t *testing.T) {
	type args struct {
		postID string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCategoriesForPost(tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategoriesForPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCategoriesForPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLikePost(t *testing.T) {
	type args struct {
		userID string
		postID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LikePost(tt.args.userID, tt.args.postID); (err != nil) != tt.wantErr {
				t.Errorf("LikePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDislikePost(t *testing.T) {
	type args struct {
		userID string
		postID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DislikePost(tt.args.userID, tt.args.postID); (err != nil) != tt.wantErr {
				t.Errorf("DislikePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdatePostLikesDislikes(t *testing.T) {
	type args struct {
		postID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdatePostLikesDislikes(tt.args.postID); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePostLikesDislikes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFilteredPosts(t *testing.T) {
	type args struct {
		loggedIn   bool
		userID     string
		categoryID string
	}
	tests := []struct {
		name    string
		args    args
		want    []Post
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFilteredPosts(tt.args.loggedIn, tt.args.userID, tt.args.categoryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilteredPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFilteredPosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllCategories(t *testing.T) {
	tests := []struct {
		name    string
		want    []Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllCategories()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPostByID(t *testing.T) {
	type args struct {
		postID string
	}
	tests := []struct {
		name    string
		args    args
		want    Post
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPostByID(tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPostByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPostOwner(t *testing.T) {
	type args struct {
		postID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPostOwner(tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostOwner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPostOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	type args struct {
		postID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeletePost(tt.args.postID); (err != nil) != tt.wantErr {
				t.Errorf("DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	type args struct {
		postID    string
		content   string
		imagePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdatePost(tt.args.postID, tt.args.content, tt.args.imagePath); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
