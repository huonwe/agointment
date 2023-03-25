package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func adminRequestings(ctx *gin.Context) {
	// token, _ := ctx.Cookie("token")
	// claim, _ := ParseToken(token)

	un_assigneds := []UnAssigned{}
	db.Model(&UnAssigned{}).Preload("Equipment").Preload("User").Find(&un_assigneds)

	ctx.HTML(http.StatusOK, "adminRequestingList.html", gin.H{
		"heads":        []string{"序号", "申请人", "设备名", "型号", "创建时间", "操作"},
		"un_assigneds": un_assigneds,
		"total":        len(un_assigneds),
	})
}

func adminRequestingsOp(ctx *gin.Context) {
	// token, _ := ctx.Cookie("token")
	// claim, _ := ParseToken(token)
	if ctx.Query("requestID") == "" || ctx.Query("op") == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "参数错误",
		})
		return
	}
	requestID := str2uint(ctx.Query("requestID"))
	switch ctx.Query("op") {
	case "reject":
		handle_resp(db.Model(&Request{}).Where(&Request{ID: requestID}).Update("Status", REJECTED).Error, ctx)
		handle_resp(db.Model(&UnAssigned{}).Delete(&UnAssigned{Request{ID: requestID}}).Error, ctx)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "",
		})
	case "assign":

	}

}
