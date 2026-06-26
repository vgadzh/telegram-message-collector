package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vgadzh/telegram-message-collector/internal/domain/user"
)

func TestGenerateTokenInvalidUserID(t *testing.T) {
	s := NewJWTService("secret1", time.Hour)

	userID := int64(0)
	role := string(user.RoleAdmin)

	token, err := s.GenerateToken(userID, role)
	if err == nil {
		t.Fatal("expected error")
	}
	if token != "" {
		t.Error("expected empty token")
	}
}

func TestGenerateTokenAndParseToken(t *testing.T) {
	s := NewJWTService("secret1", time.Hour)

	userID := int64(123)
	role := string(user.RoleAdmin)

	token, err := s.GenerateToken(userID, role)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	claims, err := s.ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != userID {
		t.Errorf("user id mismatch, got %d, want %d", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Errorf("role mismatch, got %s, want %s", claims.Role, role)
	}
	if claims.ExpiresAt == nil {
		t.Error("claims.ExpiresAt is nil")
	}
	if !claims.ExpiresAt.After(claims.IssuedAt.Time) {
		t.Error("claims.ExpiresAt is not after issuedAt")
	}
}

func TestParseToken(t *testing.T) {
	s := NewJWTService("secret1", time.Hour)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid format",
			token: "hello",
		},
	}
	for _, tt := range tests {
		claims, err := s.ParseToken(tt.token)
		if err == nil {
			t.Errorf("%s: parseToken(%s) should return error", tt.name, tt.token)
		}
		if claims != nil {
			t.Errorf("%s: parseToken(%s) should return nil", tt.name, tt.token)
		}
	}
}

func TestParseTokenAnotherSecret(t *testing.T) {
	s1 := NewJWTService("secret1", time.Hour)
	s2 := NewJWTService("secret2", time.Hour)

	token, err := s1.GenerateToken(123, string(user.RoleAdmin))
	if err != nil {
		t.Fatal(err)
	}

	claims, err := s2.ParseToken(token)
	if err == nil {
		t.Error("expected error")
	}
	if claims != nil {
		t.Error("expected nil claims")
	}
}

func TestParseTokenExpired(t *testing.T) {
	s := NewJWTService("secret1", -1*time.Hour)

	token, err := s.GenerateToken(123, string(user.RoleAdmin))
	if err != nil {
		t.Fatal(err)
	}
	claims, err := s.ParseToken(token)
	if err == nil {
		t.Error("expected error")
	}
	if claims != nil {
		t.Error("expected nil claims")
	}
}

func TestParseTokenWrongIssuer(t *testing.T) {
	s := NewJWTService("secret", time.Hour)
	claims := Claims{
		UserID: 123,
		Role:   string(user.RoleAdmin),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "another-service",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.ParseToken(tokenString)
	if err == nil {
		t.Fatal("expected issuer validation error")
	}
}
