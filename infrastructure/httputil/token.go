package httputil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/bytedance/sonic"
)

func GenerateToken(header JwtHeader, payload JwtPayload, secret string) (string, error) {
	var part1, part2, signature string
	if headerInfo, err := sonic.Marshal(header); err != nil {
		return "", err
	} else {
		part1 = base64.RawURLEncoding.EncodeToString(headerInfo)
	}

	if payloadInfo, err := sonic.Marshal(payload); err != nil {
		return "", err
	} else {
		part2 = base64.RawURLEncoding.EncodeToString(payloadInfo)
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(part1 + "." + part2))
	signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return part1 + "." + part2 + "." + signature, nil
}

func VerifyToken(Token, secret string) (header *JwtHeader, payload *JwtPayload, err error) {
	parts := strings.Split(Token, ".")
	if len(parts) != 3 {
		return nil, nil, errors.New("非法的Token长度")
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(parts[0] + "." + parts[1]))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if signature != parts[2] {
		return nil, nil, errors.New("非法的Token")
	}

	var (
		part1, part2 []byte
	)
	if part1, err = base64.RawURLEncoding.DecodeString(parts[0]); err != nil {
		return nil, nil, err
	}
	if part2, err = base64.RawURLEncoding.DecodeString(parts[1]); err != nil {
		return nil, nil, err
	}

	if err := sonic.Unmarshal(part1, header); err != nil {
		return nil, nil, err
	}
	if err := sonic.Unmarshal(part2, payload); err != nil {
		return nil, nil, err
	}
	return header, payload, nil
}
