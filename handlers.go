package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(ctx *gin.Context) {
	token, _ := ctx.Cookie("token")
	claim, err := ParseToken(token)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}

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
		})
	} else {
		ctx.HTML(http.StatusOK, "index.html", nil)
	}
}

func login(ctx *gin.Context) {
	user := User{}
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
	ctx.SetCookie("token", token, 0, "/", "", true, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "Successfully signed in.",
	})
}

func getAvailiable(ctx *gin.Context) {
	equipment := []Equipment{}
	db.Model(&Equipment{}).Where(&Equipment{Availiable: true, Name: ctx.Query("name")}).Joins("User").Find(&equipment)

	ctx.HTML(http.StatusOK, "availableList.html", gin.H{
		"heads":      []string{"序号", "设备名", "品牌", "型号", "操作"},
		"equipments": equipment,
		"total":      len(equipment),
	})
}
