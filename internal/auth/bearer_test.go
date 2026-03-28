package auth

import (
	"net/http"
	"testing"
)

func TestBearer(t *testing.T) {
	authorizedHeader := http.Header{}
	authorizedHeader.Add("Authorization", "token")
	unauthorizedHeader := http.Header{}

	tests := []struct {
		name    string
		headers http.Header
		token   string
		wantErr bool
	}{
		{
			name:    "Authorized Header",
			headers: authorizedHeader,
			token:   "token",
			wantErr: false,
		},
		{
			name:    "Unauthorized Header",
			headers: unauthorizedHeader,
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bearer, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && bearer != tt.token {
				t.Errorf("GetBearerToken() expects = %v, got %v", tt.token, bearer)
			}
		})
	}
}
