package httputil

import (
	"fmt"
	"ray_box/infrastructure/config"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	header := DefaultHeader
	payload := JwtPayload{
		Audience:    "RayBox",
		Expiration:  time.Now().Add(180 * 24 * time.Hour).Unix(),
		ID:          "1",
		Issue:       "Test",
		IssueAt:     0,
		NotBefore:   0,
		Subject:     "Session",
		UserDefined: map[string]any{"username": "User"},
	}
	secret := config.GetConfig("SECRET_KEY")
	if token, err := GenerateToken(header, payload, secret); err == nil {
		fmt.Printf("Token len: %d\nToken value: { %s }\n", len([]rune(token)), token)
	} else {
		fmt.Printf("GenerateToken error: { %s }\n", err)
	}
}
