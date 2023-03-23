package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(ctx *gin.Context) {

	token, err := ctx.Cookie("token")
	handle(err) // 已經經過了中間件，所以這裏如果出錯就直接報錯
	_, err = ParseToken(token, secret)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}
	page := ctx.Query("page")
	if page == "appoint" {
		ctx.HTML(http.StatusOK, "appoint.html", nil)
	} else if page == "status" {
		ctx.HTML(http.StatusOK, "status.html", nil)
	} else if page == "me" {
		ctx.HTML(http.StatusOK, "me.html", nil)
	} else {
		ctx.HTML(http.StatusOK, "index.html", nil)
	}
}

func login(ctx *gin.Context) {
	user := User{}
	// fmt.Println(ctx.Request.FormValue("username"))
	// fmt.Println(ctx.PostForm("username"))
	db.Where(&User{Username: ctx.PostForm("username"), Password: ctx.PostForm("password")}).First(&user)
	// fmt.Println(user)
	if user.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "Username or Password Incorrect.",
		})
		return
	}
	dict := make(map[string]interface{})
	dict["UserID"] = user.ID
	token, err := GenerateToken(dict, secret)
	handle(err)
	ctx.SetCookie("token", token, 0, "/", "", true, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "Successfully signed in.",
	})
}
