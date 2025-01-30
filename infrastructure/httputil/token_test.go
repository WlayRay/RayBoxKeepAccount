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
		Audience:    "ray_box",
		Expiration:  time.Now().Add(180 * 24 * time.Hour).Unix(),
		ID:          "1",
		Issue:       "test",
		IssueAt:     0,
		NotBefore:   0,
		Subject:     "sub",
		UserDefined: map[string]any{"username": "ray_box"},
	}
	secret := config.GetConfig("SECRET_KEY")
	if token, err := GenerateToken(header, payload, secret); err == nil {
		fmt.Printf("Token len: %d\n Token value: { %s }\n", len([]rune(token)), token)
	} else {
		fmt.Printf("GenerateToken error: { %s }\n", err)
	}
}
