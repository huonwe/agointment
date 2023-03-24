package main

import (
	"fmt"
	"net/http"
	"strconv"

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
	//&Equipment{Availiable: true, Name: ctx.Query("name")}
	db.Model(&Equipment{}).Where("name LIKE ?", "%"+ctx.Query("name")+"%").Find(&equipment)
	// log.Println(ctx.Query("name"))
	// log.Println(equipment)

	ctx.HTML(http.StatusOK, "availableList.html", gin.H{
		"heads":      []string{"序号", "设备名", "型号", "类别", "操作"},
		"equipments": equipment,
		"total":      len(equipment),
	})
}

func equipmentRequest(ctx *gin.Context) {
	token, _ := ctx.Cookie("token")
	claim, err := ParseToken(token)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}
	// log.Println(ctx.Query("equipmentID"))
	equipmentID, err := strconv.ParseUint(ctx.Query("equipmentID"), 10, 32)
	handle(err)
	request := &Request{
		EquipmentID: uint(equipmentID),
		UserID:      claim.UserID,
		Status:      REQUESTING,
	}

	request_exist := &UnAssigned{}

	db.Model(&UnAssigned{}).Where(&request).Take(&request_exist)
	if request_exist.ID != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "请勿重复申请",
		})
		return
	}

	err = db.Model(&Request{}).Create(request).Error
	handle_resp(err, ctx)
	err = db.Model(&UnAssigned{}).Create(&UnAssigned{*request}).Error
	handle_resp(err, ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "请求完成",
	})

}

func myRequest(ctx *gin.Context) {
	token, _ := ctx.Cookie("token")
	claim, err := ParseToken(token)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}

	// request = &Request{}
	user := &User{}
	db.Model(&User{}).Where(&User{ID: claim.UserID}).Preload("Requests").Preload("Requests.Equipment").Take(&user)
	// log.Println(user.Requests)
	ctx.HTML(http.StatusOK, "myRequestList.html", gin.H{
		"heads":    []string{"序号", "设备名", "型号", "类别", "创建时间", "状态", "操作"},
		"requests": user.Requests,
		"total":    len(user.Requests),
	})

}
