package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go_server/base/config"
	"go_server/base/core"
	"go_server/model/common/response"
	"net/http"
	"strings"
	"time"
)

const AuthorizationHeader = "Authorization"

// JwtMiddleware JWT中间件, 强制要求用户登录
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(AuthorizationHeader)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorObjByCode(response.ResponseCodeMissAuthToken))
			return
		}
		// 解析JWT
		member, err := ParseJWT(tokenString)
		if err != nil || member == nil || member.ExpiresAt.Time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorObjByCode(response.ResponseCodeTokenInvalid))
			return
		}
		// token有效，设置用户信息到上下文
		setClaimsToContext(c, member)
		// 继续执行
		c.Next()
	}
}

// setClaimsToContext 设置用户信息到上下文
func setClaimsToContext(c *gin.Context, member *MyClaims) {
	c.Set("userId", member.UserID)
	c.Set("roleId", member.RoleId)
	c.Set("exp", member.ExpiresAt)
	c.Set("issuer", member.Issuer)
}

// OptionalJwtMiddleware 允许用户携带 JWT，但不强制要求登录
func OptionalJwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(AuthorizationHeader)
		if tokenString != "" {
			// 解析JWT
			member, err := ParseJWT(tokenString)
			if err == nil {
				if member != nil {
					setClaimsToContext(c, member)
				}
			} else {
				core.Log.Infof("OptionalJwtMiddleware:%s", err.Error())
			}
		}
		c.Next()
	}
}

type Member struct {
	ID     int64
	RoleId int64
}

// Secret key used to sign the JWT token (In production, use environment variables to keep it secure)
var jwtSecretStr = "admin.secret.1234565"
var jwtSecretByte = []byte(jwtSecretStr)
var jwtExpireDuration = int64(24) // in seconds
var jwtIssuer = "issuer"

// MyClaims 定义了JWT中的自定义声明
type MyClaims struct {
	UserID int64 `json:"userId"`
	RoleId int64 `json:"roleId"`
	jwt.RegisteredClaims
}

func init() {
	jwtSecretStr = config.EnvConf().JWT.SigningKey
	jwtSecretByte = []byte(jwtSecretStr)
	jwtExpireDuration = config.EnvConf().JWT.ExpiresTime
	jwtIssuer = config.EnvConf().JWT.Issuer
}

// GenerateJWT 根据用户信息生成JWT
func GenerateJWT(user Member) (string, error) {
	// 设置JWT的过期时间
	d := time.Duration(jwtExpireDuration) * time.Hour
	expirationTime := time.Now().Add(d)
	// 创建JWT的声明
	claims := MyClaims{
		UserID: user.ID,
		RoleId: user.RoleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprintf("%d", user.ID),         // 设置ID
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 设置过期时间
			Issuer:    jwtIssuer,                          // 设置签发者
		},
	}

	// 使用HS256算法签署JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成并返回签名后的JWT字符串
	tokenString, err := token.SignedString(jwtSecretByte)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT 解析和验证JWT
func ParseJWT(tokenString string) (*MyClaims, error) {
	// 定义一个空的声明对象
	claims := &MyClaims{}
	// 去掉 "Bearer " 部分
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	// 解析JWT
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecretByte, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// 解析gin上下文中的用户信息

func ParseUser(c *gin.Context) *MyClaims {
	userId, _ := c.Get("userId")
	roleId, _ := c.Get("openId")
	issuer, _ := c.Get("issuer")

	claims := &MyClaims{
		UserID: userId.(int64),
		RoleId: roleId.(int64),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuer.(string),
		},
	}
	return claims
}
