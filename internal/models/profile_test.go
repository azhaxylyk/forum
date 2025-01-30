package models

import (
	"reflect"
	"testing"
)

func TestGetPostsByUser(t *testing.T) {
	type args struct {
		userID string
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
			got, err := GetPostsByUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostsByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPostsByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLikedPostsByUser(t *testing.T) {
	type args struct {
		userID string
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
			got, err := GetLikedPostsByUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLikedPostsByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLikedPostsByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDislikedPostsByUser(t *testing.T) {
	type args struct {
		userID string
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
			got, err := GetDislikedPostsByUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDislikedPostsByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDislikedPostsByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCommentsByUser(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCommentsByUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentsByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCommentsByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
