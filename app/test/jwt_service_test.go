package test

import (
	"testing"
	"time"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/services"
)

func newTestJWT() *services.JWTService {
	cfg := &configs.AuthConfig{
		JWTSecret: "test-secret-key-for-unit-tests",
		JWTExpire: 24 * time.Hour,
	}
	return services.NewJWTService(cfg)
}

func TestJWTService_GenerateTokenPair(t *testing.T) {
	jwt := newTestJWT()
	user := models.User{ID: 1, Email: "test@example.com"}

	pair, err := jwt.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("GenerateTokenPair failed: %v", err)
	}
	if pair.AccessToken == "" {
		t.Fatal("access token is empty")
	}
	if pair.RefreshToken == "" {
		t.Fatal("refresh token is empty")
	}
	if pair.AccessToken == pair.RefreshToken {
		t.Fatal("access and refresh tokens should be different")
	}
}

func TestJWTService_ValidateAccessToken_Success(t *testing.T) {
	jwt := newTestJWT()
	user := models.User{ID: 42, Email: "john@example.com"}

	pair, err := jwt.GenerateTokenPair(user)
	if err != nil {
		t.Fatal(err)
	}

	userID, err := jwt.ValidateAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("ValidateAccessToken failed: %v", err)
	}
	if userID != 42 {
		t.Errorf("userID = %d, want 42", userID)
	}
}

func TestJWTService_ValidateRefreshToken_Success(t *testing.T) {
	jwt := newTestJWT()
	user := models.User{ID: 10}

	pair, _ := jwt.GenerateTokenPair(user)

	userID, err := jwt.ValidateRefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("ValidateRefreshToken failed: %v", err)
	}
	if userID != 10 {
		t.Errorf("userID = %d, want 10", userID)
	}
}

func TestJWTService_AccessToken_CantBeUsedAsRefresh(t *testing.T) {
	jwt := newTestJWT()
	pair, _ := jwt.GenerateTokenPair(models.User{ID: 1})

	_, err := jwt.ValidateRefreshToken(pair.AccessToken)
	if err == nil {
		t.Error("access token should not validate as refresh")
	}
}

func TestJWTService_RefreshToken_CantBeUsedAsAccess(t *testing.T) {
	jwt := newTestJWT()
	pair, _ := jwt.GenerateTokenPair(models.User{ID: 1})

	_, err := jwt.ValidateAccessToken(pair.RefreshToken)
	if err == nil {
		t.Error("refresh token should not validate as access")
	}
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	jwt := newTestJWT()

	_, err := jwt.ValidateAccessToken("invalid.token.string")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestJWTService_ValidateToken_WrongSecret(t *testing.T) {
	jwt1 := newTestJWT()
	jwt2 := services.NewJWTService(&configs.AuthConfig{
		JWTSecret: "different-secret",
		JWTExpire: 24 * time.Hour,
	})

	pair, _ := jwt1.GenerateTokenPair(models.User{ID: 1})

	_, err := jwt2.ValidateAccessToken(pair.AccessToken)
	if err == nil {
		t.Error("expected error for wrong secret")
	}
}

func TestJWTService_ValidateToken_Empty(t *testing.T) {
	jwt := newTestJWT()

	_, err := jwt.ValidateAccessToken("")
	if err == nil {
		t.Error("expected error for empty token")
	}
}

func TestJWTService_BackwardCompat_ValidateToken(t *testing.T) {
	jwt := newTestJWT()
	user := models.User{ID: 99}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := jwt.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.ID != 99 {
		t.Errorf("ID = %d, want 99", parsed.ID)
	}
}
