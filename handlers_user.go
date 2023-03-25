package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func myRequest(ctx *gin.Context) {
	token, _ := ctx.Cookie("token")
	claim, err := ParseToken(token)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}
	user := &User{}
	db.Model(&User{}).Where(&User{ID: claim.UserID}).Preload("Requests").Preload("Requests.Equipment").Take(&user)

	ctx.HTML(http.StatusOK, "myRequestList.html", gin.H{
		"heads":    []string{"序号", "设备名", "型号", "设备ID", "创建时间", "状态", "操作"},
		"requests": user.Requests,
		"total":    len(user.Requests),
	})

}

func myRequestOp(ctx *gin.Context) {
	token, _ := ctx.Cookie("token")
	claim, err := ParseToken(token)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}
	op := ctx.Query("op")
	request := Request{
		ID:     str2uint(ctx.Query("requestID")),
		UserID: claim.UserID,
	}
	switch op {
	case "cancel":
		handle_resp(db.Model(&UnAssigned{}).Delete(&UnAssigned{request}).Error, ctx)
		handle_resp(db.Model(&request).Update("status", CANCELED).Delete(&request).Error, ctx)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "取消成功",
		})
	case "delete":
		handle_resp(db.Model(&request).Delete(&request).Error, ctx)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "删除成功",
		})
	}
}
