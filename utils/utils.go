package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JwtSecret = []byte("ABAB")

type Claims struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateToken 签发token
func GenerateToken(id uint, username string, password string) (string, error) {
	notTime := time.Now()
	expireTime := notTime.Add(24 * time.Hour)
	claims := Claims{
		ID:       id,
		UserName: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "todo_list",
		},
	}
	// 使用HS256算法计算结构体哈希
	tokeClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 根据密钥对哈希签名
	token, err := tokeClaims.SignedString(JwtSecret)
	return token, err
}

// ParseToken 验证token
func ParseToken(token string) (*Claims, error) {
	// 根据私钥解密Token为哈希，并把值填充到Claims类型的结构体中
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	// 如果内容不为空
	if tokenClaims != nil {
		// 判断是否合法并且 是否可以转换为Claims类型
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
