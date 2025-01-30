package models

import "testing"

func TestAuthenticateOrRegisterOAuthUser(t *testing.T) {
	type args struct {
		email            string
		username         string
		provider         string
		moderatorRequest bool
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
			got, err := AuthenticateOrRegisterOAuthUser(tt.args.email, tt.args.username, tt.args.provider, tt.args.moderatorRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticateOrRegisterOAuthUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthenticateOrRegisterOAuthUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckEmailExists(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckEmailExists(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckEmailExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckEmailExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckUsernameExists(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckUsernameExists(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUsernameExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckUsernameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	type args struct {
		email              string
		username           string
		password           string
		isModeratorRequest bool
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
			got, err := RegisterUser(tt.args.email, tt.args.username, tt.args.password, tt.args.isModeratorRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegisterUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticateUser(t *testing.T) {
	type args struct {
		email    string
		password string
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
			got, err := AuthenticateUser(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthenticateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIDBySessionToken(t *testing.T) {
	type args struct {
		sessionToken string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetIDBySessionToken(tt.args.sessionToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIDBySessionToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIDBySessionToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetIDBySessionToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
