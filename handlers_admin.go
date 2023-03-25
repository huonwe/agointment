package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func adminRequestings(ctx *gin.Context) {
	// token, _ := ctx.Cookie("token")
	// claim, _ := ParseToken(token)

	un_assigneds := []UnAssigned{}
	db.Model(&UnAssigned{}).Preload("Equipment").Find(&un_assigneds)

	ctx.HTML(http.StatusOK, "adminRequestingList.html", gin.H{
		"heads":        []string{"序号", "设备名", "型号", "创建时间", "操作"},
		"un_assigneds": un_assigneds,
		"total":        len(un_assigneds),
	})
}

func adminRequestingsOp(ctx *gin.Context) {
	// token, _ := ctx.Cookie("token")
	// claim, _ := ParseToken(token)

	switch ctx.Query("op") {
	case "reject":

	}

	un_assigneds := []UnAssigned{}
	db.Model(&UnAssigned{}).Preload("Equipment").Find(&un_assigneds)

	ctx.HTML(http.StatusOK, "adminRequestingList.html", gin.H{
		"heads":        []string{"序号", "设备名", "型号", "创建时间", "操作"},
		"un_assigneds": un_assigneds,
		"total":        len(un_assigneds),
	})
}
