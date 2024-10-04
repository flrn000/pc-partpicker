package utils

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestValidateJWT(t *testing.T) {
	type args struct {
		token     string
		jwtSecret string
	}
	t.Setenv("JWT_SECRET", "NqSNAcr98yn7nUmZjFD9iVF4lwv9Wn3N51Kik8y3N4fVZWsD5Tac6okJ9qMWWeByE7o/ETltV2hwWZzqp47w8Q==")
	jwtSecret := os.Getenv("JWT_SECRET")

	token, _ := GenerateJWT("1", jwtSecret, time.Hour)
	tokenShort, _ := GenerateJWT("2", jwtSecret, 5*time.Millisecond)

	tests := []struct {
		name       string
		args       args
		wantUserID int
		wantErr    bool
	}{
		{
			name: "Is valid",
			args: args{
				token:     token,
				jwtSecret: jwtSecret,
			},
			wantUserID: 1,
			wantErr:    false,
		},
		{
			name: "Invalid secret",
			args: args{
				token:     token,
				jwtSecret: "",
			},
			wantUserID: 0,
			wantErr:    true,
		},
		{
			name: "Short expiry time",
			args: args{
				token:     tokenShort,
				jwtSecret: jwtSecret,
			},
			wantUserID: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.args.token, tt.args.jwtSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %v - error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("Test %v - got = %v, want %v", tt.name, gotUserID, tt.wantUserID)
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

func TestGetAuthToken(t *testing.T) {
	tests := []struct {
		name      string
		input     http.Header
		want      string
		wantError bool
	}{
		{
			name: "Valid auth header",
			input: http.Header{
				"Authorization": []string{"Bearer ftest12395"},
			},
			want: "ftest12395",
		},
		{
			name: "No header",
			input: http.Header{
				"Content-Type": []string{"text/html"},
			},
			want:      "",
			wantError: true,
		},
		{
			name: "No bearer",
			input: http.Header{
				"Authorization": []string{"ftest12395"},
			},
			want:      "",
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAuthToken(tt.input)
			if err == nil && tt.wantError {
				t.Errorf("Test %v - expected error, got: %v", tt.name, err)
				return
			}
			if got != tt.want {
				t.Errorf("Test %v - got: %v, want: %v", tt.name, got, tt.want)
				return
			}
		})
	}
}
