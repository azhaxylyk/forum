package models

import (
	"reflect"
	"testing"
)

func TestCreateComment(t *testing.T) {
	type args struct {
		postID  string
		userID  string
		content string
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
			if err := CreateComment(tt.args.postID, tt.args.userID, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLikeComment(t *testing.T) {
	type args struct {
		userID    string
		commentID string
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
			if err := LikeComment(tt.args.userID, tt.args.commentID); (err != nil) != tt.wantErr {
				t.Errorf("LikeComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDislikeComment(t *testing.T) {
	type args struct {
		userID    string
		commentID string
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
			if err := DislikeComment(tt.args.userID, tt.args.commentID); (err != nil) != tt.wantErr {
				t.Errorf("DislikeComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCommentLikesDislikes(t *testing.T) {
	type args struct {
		commentID string
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
			if err := UpdateCommentLikesDislikes(tt.args.commentID); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCommentLikesDislikes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCommentsForPost(t *testing.T) {
	type args struct {
		postID string
	}
	tests := []struct {
		name    string
		args    args
		want    []Comment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCommentsForPost(tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentsForPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCommentsForPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeInput(tt.args.input); got != tt.want {
				t.Errorf("SanitizeInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidContent(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidContent(tt.args.content); got != tt.want {
				t.Errorf("IsValidContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCommentOwner(t *testing.T) {
	type args struct {
		commentID string
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
			got, err := GetCommentOwner(tt.args.commentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentOwner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCommentOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}
