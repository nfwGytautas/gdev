package jwt

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nfwGytautas/mstk/lib/gdev/array"
)

// PUBLIC TYPES
// ========================================================================

/*
Struct for containing token info
*/
type TokenInfo struct {
	Valid bool
	ID    uint
	Role  string
}

/*
API Secret for parsing JWT tokens
*/
var APISecret string

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Middleware for authenticating

Usage:
r.Use(common.JwtAuthenticationMiddleware())
*/
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		info, err := ParseToken(c)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Token")
			c.Abort()
			return
		}

		if !info.Valid {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

/*
Middleware for authorization

Usage:
r.Use(common.AuthorizationMiddleware([]string{"role"}))
*/
func AuthorizationMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := ParseToken(c)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, "Token")
			c.Abort()
		}

		if !array.IsElementInArray(roles, info.Role) {
			c.String(http.StatusUnauthorized, "Access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

/*
Parse a token from gin context
*/
func ParseToken(c *gin.Context) (TokenInfo, error) {
	result := TokenInfo{}
	result.Valid = false

	tokenString := c.Query("token")

	if tokenString == "" {
		// Token empty check if it is inside Authorization header
		tokenString = c.Request.Header.Get("Authorization")

		// Since this is bearer token we need to parse the token out
		if len(strings.Split(tokenString, " ")) == 2 {
			tokenString = strings.Split(tokenString, " ")[1]
		} else {
			return result, errors.New("invalid request")
		}
	}

	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(APISecret), nil
	})

	if err != nil {
		return result, err
	}

	// Token valid fill token information
	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	if !ok || !jwtToken.Valid {
		return result, nil
	}

	// User id
	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
	if err != nil {
		return result, nil
	}

	result.ID = uint(uid)

	// Role
	result.Role = claims["role"].(string)

	result.Valid = true
	return result, nil
}

// PRIVATE FUNCTIONS
// ========================================================================
