package utils

import (
	"os"
	"testing"
	"time"
)

func TestValidateJWT(t *testing.T) {
	type args struct {
		token     string
		jwtSecret string
	}
	tests := []struct {
		name       string
		args       args
		wantUserID int
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.args.token, tt.args.jwtSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	type args struct {
		subject   string
		secret    string
		expiresIn time.Duration
	}
	t.Setenv("JWT_SECRET", "NqSNAcr98yn7nUmZjFD9iVF4lwv9Wn3N51Kik8y3N4fVZWsD5Tac6okJ9qMWWeByE7o/ETltV2hwWZzqp47w8Q==")
	jwtSecret := os.Getenv("JWT_SECRET")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid JWT",
			args: args{
				subject:   "1",
				secret:    jwtSecret,
				expiresIn: time.Hour,
			},
			wantErr: false,
		},
		{
			name: "Invalid JWT Secret",
			args: args{
				subject:   "1",
				secret:    "",
				expiresIn: time.Hour,
			},
			wantErr: true,
		},
		{
			name: "Invalid JWT Subject",
			args: args{
				subject:   "",
				secret:    jwtSecret,
				expiresIn: time.Hour,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateJWT(tt.args.subject, tt.args.secret, tt.args.expiresIn)
			if (err != nil) && !tt.wantErr {
				t.Errorf("Test %v - error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

			if err == nil && tt.wantErr {
				t.Errorf("Test %v - Expected an error: %v", tt.name, err)
				return
			}
		})
	}
}
