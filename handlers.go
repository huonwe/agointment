package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func redirect2home(ctx *gin.Context) {
	ctx.Redirect(http.StatusTemporaryRedirect, "/home/index")
}

func login(ctx *gin.Context) {
	user := User{}
	username_ := ctx.PostForm("username")
	password_ := ctx.PostForm("password")
	if username_ == "" || password_ == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "参数错误",
		})
		return
	}
	db.Model(&User{}).Where(&User{Name: username_, Password: md5_str(password_)}).Preload("Department").First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "用户名或密码错误",
		})
		return
	}

	if !user.Department.Availiable && !user.IsAdmin {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "用户不可用",
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

func signup(ctx *gin.Context) {
	if ctx.PostForm("username") == "" || ctx.PostForm("pass") == "" || ctx.PostForm("dept") == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "参数不足",
		})
		return
	}

	user := User{
		Name:     ctx.PostForm("username"),
		Password: ctx.PostForm("pass"),
		// DepartmentName: ctx.PostForm("dept"),
	}

	dept := Department{}
	db.Model(&Department{}).Where(&Department{Name: ctx.PostForm("dept"), Availiable: true}).Take(&dept)
	if dept.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "没有这个部门或部门不可用",
		})
		return
	}
	user.DepartmentID = dept.ID
	// user.DeptName = dept.Name

	var count int64
	db.Model(&User{}).Where(&User{Name: user.Name, DepartmentID: dept.ID}).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "用户名已存在",
		})
		return
	}

	err := db.Create(&user).Error
	if err != nil {
		handle_resp(err, ctx)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "注册成功",
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
		dept := Department{}
		db.Take(&dept, user.DepartmentID)
		ctx.HTML(http.StatusOK, "me.html", gin.H{
			"greeting": fmt.Sprintf("欢迎您，%s 的 %s", dept.Name, user.Name),
			"user":     user,
		})
	} else if name == "index" {
		ctx.HTML(http.StatusOK, "index.html", nil)
	} else {
		ctx.Redirect(http.StatusTemporaryRedirect, "/home/index")
	}
}
