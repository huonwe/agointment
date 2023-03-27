package main

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func myRequest(ctx *gin.Context) {
	token, _ := ctx.Cookie("token")
	claim, err := ParseToken(token)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}
	// log.Println(ctx.Query("name"))
	page := str2int(ctx.Query("page"))
	pageSize := str2int(ctx.Query("pageSize"))
	var total int64 = 0
	user := &User{
		ID: claim.UserID,
	}
	db.Model(&user).Preload("Requests", order_desc_createdAt, func(db *gorm.DB) *gorm.DB {
		return db.Where("equipment_name LIKE ?", "%"+ctx.Query("name")+"%").Limit(pageSize).Offset((page - 1) * pageSize)
	}).Preload("Requests.Equipment").Take(&user)

	db.Model(&Request{}).Where("user_id = ?", claim.UserID).Where("equipment_name LIKE ?", "%"+ctx.Query("name")+"%").Count(&total)
	ctx.HTML(http.StatusOK, "myRequestList.html", gin.H{
		"heads":      []string{"请求序号", "设备名", "型号", "设备ID", "状态", "操作"},
		"requests":   user.Requests,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"total_page": int(math.Ceil(float64(total) / float64(pageSize))),
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
