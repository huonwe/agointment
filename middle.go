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
			ctx.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
			ctx.Abort()
			return
		}
		user := User{}
		db.Where(&User{ID: claims.UserID}).Take(&user)
		ctx.Set("user", user)
		ctx.Next()
	}
}

func MustAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, _ := ctx.Cookie("token")

		claims, _ := ParseToken(token)
		user := &User{}
		err := db.Where(&User{ID: claims.UserID}).Take(&user).Error
		if err != nil {
			ctx.String(http.StatusBadRequest, "认证失败")
			ctx.Abort()
			return
		}
		if !user.IsAdmin || user.ID == 0 {
			ctx.Redirect(http.StatusTemporaryRedirect, "/home")
			ctx.Abort()
			return
		}
		ctx.Set("userID", user.ID)
		ctx.Next()
	}
}
