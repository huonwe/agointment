package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	token, _ := ctx.Cookie("token")
	claim, _ := ParseToken(token)

	user := User{}
	db.Model(&User{}).Joins("Department").Take(&user, claim.UserID)

	page := ctx.Query("page")
	if page == "appoint" {
		ctx.HTML(http.StatusOK, "appoint.html", nil)
	} else if page == "status" {
		ctx.HTML(http.StatusOK, "status.html", nil)
	} else if page == "me" {
		ctx.HTML(http.StatusOK, "me.html", gin.H{
			"greeting": fmt.Sprintf("欢迎您，%s 的 %s", user.Department.Name, user.Name),
			"user":     user,
		})
	} else if page == "index" {
		ctx.HTML(http.StatusOK, "index.html", nil)
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, "/?page=index")
	}
}
