package httputil

var (
	DefaultHeader = JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
)

type JwtHeader struct {
	Algo string `json:"alg"` //哈希算法，默认为HMAC SHA256(写为 HS256)
	Type string `json:"typ"` //令牌(token)类型，统一写为JWT
}

type JwtPayload struct {
	ID          string         `json:"jti"` //JWT ID用于标识该JWT
	Issue       string         `json:"iss"` //发行人。比如微信
	Audience    string         `json:"aud"` //受众人。比如王者荣耀
	Subject     string         `json:"sub"` //主题
	IssueAt     int64          `json:"iat"` //发布时间,精确到秒
	NotBefore   int64          `json:"nbf"` //在此之前不可用,精确到秒
	Expiration  int64          `json:"exp"` //到期时间,精确到秒
	UserDefined map[string]any `json:"ud"`  //用户自定义的其他字段
}


