package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func redirect2home(ctx *gin.Context) {
	ctx.Redirect(http.StatusPermanentRedirect, "/home/index")
}

func login(ctx *gin.Context) {
	user := User{}
	if ctx.PostForm("username") == "" || ctx.PostForm("password") == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "Username or Password Empty.",
		})
		return
	}
	db.Where(&User{Name: ctx.PostForm("username"), Password: ctx.PostForm("password")}).Take(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "Username or Password Incorrect.",
		})
		return
	}
	token, err := GenerateToken(user.ID, user.Name)
	handle(err)
	ctx.SetCookie("token", token, 0, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "Successfully signed in.",
	})
}

func index(ctx *gin.Context) {
	value, exist := ctx.Get("user")
	user, ok := value.(User)
	if !exist || !ok {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	name := ctx.Param("name")
	if name == "appoint" {
		ctx.HTML(http.StatusOK, "appoint.html", nil)
	} else if name == "status" {
		ctx.HTML(http.StatusOK, "status.html", nil)
	} else if name == "me" {
		ctx.HTML(http.StatusOK, "me.html", gin.H{
			"greeting": fmt.Sprintf("欢迎您，%s 的 %s", user.DepartmentName, user.Name),
			"user":     user,
		})
	} else if name == "index" {
		ctx.HTML(http.StatusOK, "index.html", nil)
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, "/home/index")
	}
}
