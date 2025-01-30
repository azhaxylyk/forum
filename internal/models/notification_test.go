package models

import (
	"reflect"
	"testing"
)

func TestCreateNotification(t *testing.T) {
	type args struct {
		userID     string
		actionBy   string
		actionType string
		targetID   string
		targetType string
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
			if err := CreateNotification(tt.args.userID, tt.args.actionBy, tt.args.actionType, tt.args.targetID, tt.args.targetType); (err != nil) != tt.wantErr {
				t.Errorf("CreateNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetNotificationsForUser(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []Notification
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNotificationsForUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotificationsForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNotificationsForUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkNotificationAsRead(t *testing.T) {
	type args struct {
		notificationID string
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
			if err := MarkNotificationAsRead(tt.args.notificationID); (err != nil) != tt.wantErr {
				t.Errorf("MarkNotificationAsRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkAllNotificationsAsRead(t *testing.T) {
	type args struct {
		userID string
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
			if err := MarkAllNotificationsAsRead(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("MarkAllNotificationsAsRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteReadNotifications(t *testing.T) {
	type args struct {
		userID string
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
			if err := DeleteReadNotifications(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteReadNotifications() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotification_GetMessage(t *testing.T) {
	tests := []struct {
		name    string
		n       *Notification
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetMessage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Notification.GetMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Notification.GetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUsernameByID(t *testing.T) {
	type args struct {
		userID string
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
			got, err := GetUsernameByID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsernameByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUsernameByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUnreadNotificationCount(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUnreadNotificationCount(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUnreadNotificationCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUnreadNotificationCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
