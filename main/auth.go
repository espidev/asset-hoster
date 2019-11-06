package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

func AuthRoute(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	password := c.PostForm("password")

	var userF *IUser

	for _, user := range db.Users {
		if user.Username == username {
			userF = user
		}
	}

	if userF == nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"incorrectPassword": false,
			"noUser":            true,
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userF.Password), []byte(password))

	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"incorrectPassword": true,
			"noUser":            false,
		})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = userF.Username
	claims["expires"] = time.Now().Add(time.Second * time.Duration(86400)).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(config.Secret))

	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	c.SetCookie("GOSESSID", tokenString, 86400, "/", config.Domain, false, false)

	// successful
	c.Redirect(302, "/")
}

func CheckJWTToken(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { // verify token authenticity
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// get claims of JWT token

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

		for _, userf := range db.Users {
			if userf.Username == claims["id"].(string) {
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
