package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/steppbol/activity-manager/configs"
	"github.com/steppbol/activity-manager/internal/dtos"
	"github.com/steppbol/activity-manager/internal/utils/exception"
)

type TokenDetail struct {
	AccessToken            string
	RefreshToken           string
	AccessID               uuid.UUID
	RefreshID              uuid.UUID
	AccessTokenExpireDate  int64
	RefreshTokenExpireDate int64
}

type AccessDetail struct {
	AccessID string
	UserID   uint
}

type RefreshDetail struct {
	RefreshID string
	UserID    uint
}

type JWTMiddleware struct {
	config    *configs.Security
	jwtSecret []byte
}

type claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}

func NewJWTMiddleware(conf *configs.Security) (*JWTMiddleware, error) {
	js, err := conf.JWTSecret.MarshalBinary()
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
		token := j.extractToken(c)

		if token == "" {
			code = exception.Unauthorized
		} else {
			_, err := j.verifyToken(c)
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

func (j JWTMiddleware) GenerateToken(username, password string, userId uint) (*TokenDetail, error) {
	var td TokenDetail

	err := j.getAccessToken(&td, username, password, userId)
	if err != nil {
		return nil, err
	}

	err = j.getRefreshToken(&td, username, password, userId)
	if err != nil {
		return nil, err
	}

	return &td, err
}

func (j JWTMiddleware) ExtractAccessTokenMetadata(c *gin.Context) (*AccessDetail, error) {
	token, err := j.verifyToken(c)
	if err != nil {
		return nil, err
	}

	cl, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		acId, iOk := cl["access_id"].(string)
		if !iOk {
			return nil, err
		}

		cId, iErr := strconv.Atoi(cl["user_id"].(string))
		if iErr != nil {
			return nil, iErr
		}

		return &AccessDetail{
			AccessID: acId,
			UserID:   uint(cId),
		}, nil
	}

	return nil, err
}

func (j JWTMiddleware) ExtractRefreshTokenMetadata(tk string) (*RefreshDetail, error) {
	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.config.RefreshSecret, nil
	})
	if err != nil {
		return nil, err
	}

	_, ok := token.Claims.(jwt.Claims)
	if !ok && !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", token.Header["alg"])
	}

	cl, ok := token.Claims.(jwt.MapClaims)

	var rd RefreshDetail

	if ok && token.Valid {
		refreshId, iOk := cl["refresh_id"].(string)
		if !iOk {
			return nil, fmt.Errorf("invalid token: %v", token.Header["alg"])
		}

		cId, iErr := strconv.Atoi(cl["user_id"].(string))
		if iErr != nil {
			return nil, iErr
		}

		rd.UserID = uint(cId)
		rd.RefreshID = refreshId
	}

	return &rd, nil
}

func (j JWTMiddleware) verifyToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := j.extractToken(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j JWTMiddleware) extractToken(c *gin.Context) string {
	token := c.Request.Header["Authorization"]
	arg := strings.Split(token[0], " ")

	if len(arg) == 2 {
		return arg[1]
	}

	return ""
}

func (j JWTMiddleware) getAccessToken(td *TokenDetail, username, password string, userId uint) error {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour).Unix()

	c := claims{
		j.encodeMD5(username),
		j.encodeMD5(password),
		fmt.Sprintf("%d", userId),
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    j.config.Issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := tokenClaims.SignedString(j.jwtSecret)

	td.AccessToken = token
	td.AccessID = uuid.New()
	td.AccessTokenExpireDate = expireTime

	return err
}

func (j JWTMiddleware) getRefreshToken(td *TokenDetail, username, password string, userId uint) error {
	nowTime := time.Now()
	expireTime := nowTime.Add(168 * time.Hour).Unix()

	c := claims{
		j.encodeMD5(username),
		j.encodeMD5(password),
		fmt.Sprintf("%d", userId),
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    j.config.Issuer,
			Id:        fmt.Sprintf("%s_%d", td.AccessID.String(), userId),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := tokenClaims.SignedString(j.jwtSecret)

	td.RefreshToken = token
	td.RefreshID = uuid.New()
	td.RefreshTokenExpireDate = expireTime

	return err
}

func (j JWTMiddleware) encodeMD5(value string) string {
	md := md5.New()
	md.Write([]byte(value))

	return hex.EncodeToString(md.Sum(nil))
}
