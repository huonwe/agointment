package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		_, err = ParseToken(token)
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		ctx.Next()
	}
}

func MustAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, _ := ctx.Cookie("token")
		claims, _ := ParseToken(token)
		user := &User{}
		db.Where(&User{ID: claims.UserID}).Take(&user)
		if !user.IsAdmin || user.ID == 0 {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		ctx.Next()
	}
}
