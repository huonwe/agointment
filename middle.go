package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		_, err = ParseToken(token)
		if err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		ctx.Next()
	}
}
