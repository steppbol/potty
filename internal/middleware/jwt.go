package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/steppbol/activity-manager/configs"
	"github.com/steppbol/activity-manager/internal/dtos"
	"github.com/steppbol/activity-manager/internal/utils/exception"
)

type JWTMiddleware struct {
	config    *configs.Security
	jwtSecret []byte
}

type claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func Setup(conf *configs.Security) (*JWTMiddleware, error) {
	js, err := conf.Secret.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return &JWTMiddleware{
		config:    conf,
		jwtSecret: js,
	}, nil
}

func (j JWTMiddleware) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = exception.Success
		token := c.Request.Header["Token"]

		if token == nil || token[0] == "" {
			code = exception.Unauthorized
		} else {
			_, err := j.parseToken(token[0])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				default:
					code = exception.InternalServerError
				}
			}
		}

		if code != exception.Success {
			dtos.CreateJSONResponse(c, code, code, data)

			c.Abort()
			return
		}

		c.Next()
	}
}

func (j JWTMiddleware) GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	c := claims{
		j.encodeMD5(username),
		j.encodeMD5(password),
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    j.config.Issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := tokenClaims.SignedString(j.jwtSecret)

	return token, err
}

func (j JWTMiddleware) parseToken(token string) (*claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.jwtSecret, nil
	})

	if tokenClaims != nil {
		if c, ok := tokenClaims.Claims.(*claims); ok && tokenClaims.Valid {
			return c, nil
		}
	}

	return nil, err
}

func (j JWTMiddleware) encodeMD5(value string) string {
	md := md5.New()
	md.Write([]byte(value))

	return hex.EncodeToString(md.Sum(nil))
}
