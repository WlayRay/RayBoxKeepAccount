package test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"ray_box/infrastructure/config"
	"ray_box/infrastructure/httputil"
	"runtime"
	"testing"
	"time"

	"github.com/bytedance/sonic"
)

type Session struct {
	UserId    int64  `json:"UserID"`
	SessionId string `json:"SessionID"`
}

func GenerateToken() {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Println("Error generating random bytes:", err)
		return
	}
	sessionId := fmt.Sprintf("%s:%d", base64.URLEncoding.EncodeToString(randomBytes), time.Now().UnixMilli())
	session := Session{
		SessionId: sessionId,
		UserId:    1,
	}

	var err error
	sessionBytes, err := sonic.Marshal(session)
	if err != nil {
		log.Printf("Error marshaling session: %v\n", err)
		return
	}
	// log.Printf("Session JSON: %s\n", string(sessionBytes))

	secret := config.GetConfig("SECRET_KEY")
	// 创建HMAC-SHA256签名
	hc := hmac.New(sha256.New, []byte(secret))
	hc.Write(sessionBytes)
	sign := base64.URLEncoding.EncodeToString(hc.Sum(nil))

	encodedSession := base64.URLEncoding.EncodeToString(sessionBytes)
	token := fmt.Sprintf("%s|%s", encodedSession, sign)
	// log.Printf("Token: %s\n", token)

	// 使用正确的AES密钥长度
	key := []byte(secret)[67:99]
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Error creating cipher: %v\n", err)
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("Error creating GCM: %v\n", err)
		return
	}

	// 生成随机的Nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		log.Printf("Error generating nonce: %v\n", err)
		return
	}

	// 加密token，并将nonce包含在密文前面
	encryptedToken := gcm.Seal(nonce, nonce, []byte(token), nil)
	encodedToken := base64.URLEncoding.EncodeToString(encryptedToken)
	DoNothing(encodedToken)
	// log.Printf("Token len: %d\nToken value: { %s }\n", len([]rune(encodedToken)), encodedToken)
}

func GenerateTokenByJWT() {
	header := httputil.DefaultHeader
	payload := httputil.JwtPayload{
		Audience:   "RayBox",
		Expiration: time.Now().Add(180 * 24 * time.Hour).Unix(),
		ID:         fmt.Sprintf("%d:RayBox", time.Now().Unix()),
		Issue:      "Test",
		IssueAt:    time.Now().Unix(),
		// NotBefore:   0,
		// Subject:     "Session",
		UserDefined: map[string]any{"username": "User"},
	}
	secret := config.GetConfig("SECRET_KEY")
	token, err := httputil.GenerateToken(header, payload, secret)
	DoNothing(token)
	if err == nil {
		// log.Printf("Token len: %d\nToken value: { %s }\n", len([]rune(token)), token)
	} else {
		log.Printf("GenerateToken error: { %s }\n", err)
	}
}

// 服务端生产Token
// 1.创建用户数据
// 2.对用户数据进行HMAC签名
// 3. 将用户数据和签名数据拼在一起（用"|"区分)
// 4. 对拼接字符穿用SHA256加密，然后用base64编码得到Token
func TestGenerateToken(t *testing.T) {
	fmt.Println("GenerateToke:")
	GenerateToken()
	fmt.Println("GenerateTokenByJWT:")
	GenerateTokenByJWT()
}

func DoNothing(any) {
	// do nothing
}

// goos: windows
// goarch: amd64
// pkg: ray_box/common/test
// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// === RUN   BenchmarkGenerateToken
// BenchmarkGenerateToken
// BenchmarkGenerateToken-16
//
//	313456              3738 ns/op            3932 B/op         35 allocs/op
func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateToken()
	}
}

// goos: windows
// goarch: amd64
// pkg: ray_box/common/test
// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// === RUN   BenchmarkGenerateTokenByJWT
// BenchmarkGenerateTokenByJWT
// BenchmarkGenerateTokenByJWT-16
//   393970              2946 ns/op            3111 B/op         27 allocs/op
func BenchmarkGenerateTokenByJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateTokenByJWT()
	}
}

func generateRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func TestGenerateKey(t *testing.T) {
	secretKey, err := generateRandomString(256)
	if err != nil {
		log.Fatalf("Failed to generate SECRET_KEY: %v", err)
	}

	tokenYmlContent := fmt.Sprintf(`SECRET_KEY: "${SECRET_KEY}|%s"`, secretKey)

	_, filename, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(filename), "../..")
	configDir := filepath.Join(rootPath, "config")
	err = os.WriteFile(filepath.Join(configDir, "token.yml"), []byte(tokenYmlContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write to token.yml: %v", err)
	}

	log.Println("token.yml updated successfully")
}
