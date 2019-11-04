package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

func CheckJWTToken(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { // Verify token authenticity
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// Get claims of JWT token

func GetJWTClaims(token string, secret string) (jwt.MapClaims, error) {
	tok, err := CheckJWTToken(token, secret)
	if err != nil {
		return nil, err
	}
	if !tok.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("not ok")
	}
	return claims, nil
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, err := c.Cookie("GOSESSID")
		if err != nil {
			if config.Debug {
				log.Println("[Token] " + err.Error())
			}
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		claims, err := GetJWTClaims(sess, config.Secret)
		if err != nil {
			if config.Debug {
				log.Println("[Token]" + err.Error())
			}
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		exp := claims["expires"].(float64)

		if exp < float64(time.Now().Unix()) {
			if config.Debug {
				log.Println("[Token] Expired token.")
			}
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		var user *IUser

		for _, userf := range users {
			if userf.UserName == claims["id"].(string) {
				user = userf
			}
		}

		if user == nil {
			if config.Debug {
				log.Println("[Token] IUser not found.")
			}
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}