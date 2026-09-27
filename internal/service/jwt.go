package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/KernyrMindDev/core/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// 不对外暴露
var jwtKey []byte

// generateRandomKey 生成指定字节长度的强随机数
func generateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func InitJWTKey(db *gorm.DB) {
	config := model.AppDatabase{Key: "jwt_secret"}

	// 查询数据库
	err := db.First(&config).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("未找到密钥，正在生成...")

		randomBytes, _ := generateRandomKey(32)
		config.Value = base64.StdEncoding.EncodeToString(randomBytes)

		// 插入数据库
		db.Create(&config)
		jwtKey = randomBytes
	} else {
		// 找到记录，直接解码
		decodedKey, _ := base64.StdEncoding.DecodeString(config.Value)
		jwtKey = decodedKey
	}
}

// GenerateToken 生成JWT Token
func GenerateToken(uid string) (string, error) {
	// 设置过期时间，24小时过期
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建 Claims
	claims := &model.JwtClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
			Issuer:    "KernyrMind",                       // 签发人
		},
	}

	// 使用 HS256 算法创建 Token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥进行签名并获取完整的 JWT 字符串
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenString string) (*model.JwtClaims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &model.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为 HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("错误的签名算法")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 校验 Claims 并返回
	if claims, ok := token.Claims.(*model.JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的 Token")
}
