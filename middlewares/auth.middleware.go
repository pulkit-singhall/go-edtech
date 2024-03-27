package middlewares

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pulkit-singhall/go-edtech/utils"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		token, err := c.Cookie("token")
		refresh, refErr := c.Cookie("refresh_token")
		if err != nil || refErr != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": utils.TokenError.Error(), "detail": err.Error()})
			return
		}
		tok, verErr := utils.VerifyToken(token)
		refTok, verRefErr := utils.VerifyRefreshToken(refresh)
		if verErr != nil {
			c.AbortWithStatusJSON(410, gin.H{"error": utils.TokenError.Error(), "detail": verErr.Error()})
			return
		}
		if verRefErr != nil {
			c.AbortWithStatusJSON(410, gin.H{"error": utils.TokenError.Error(), "detail": verRefErr.Error()})
			return
		}
		if !tok.Valid { // access token expired
			if !refTok.Valid { // refresh token expired
				c.AbortWithStatusJSON(415, gin.H{"message": "tokens expired. pls login again!"})
				return
			} else {
				// refresh the tokens again
				claims := refTok.Claims.(jwt.MapClaims)
				email := claims["email"].(string)
				sameErr := utils.SameRefreshToken(email, refTok.Raw)
				if sameErr != nil {
					c.AbortWithStatusJSON(500, gin.H{"error": utils.TokenError.Error(), "detail": sameErr.Error()})
					return
				}
				// generate new tokens
				token, refresh, err := utils.GenerateNewTokens(email)
				if err != nil {
					c.AbortWithStatusJSON(500, gin.H{"error": utils.InternalServerError.Error(), "detail": err.Error()})
					return
				}
				c.SetCookie("token", token, int(time.Hour*48), "/", "localhost", true, true)
				c.SetCookie("refresh_token", refresh, int(time.Hour*240), "/", "localhost", true, true)
				c.Set("email", email)
				c.Next()
			}
		} else {
			claims := tok.Claims.(jwt.MapClaims)
			email := claims["email"].(string)
			c.Set("email", email)
			c.Next()
		}
	}
}
