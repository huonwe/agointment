package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func adminRequestings(ctx *gin.Context) {
	// token, _ := ctx.Cookie("token")
	// claim, _ := ParseToken(token)

	un_assigneds := []UnAssigned{}
	db.Model(&UnAssigned{}).Order("created_at desc").Preload("Equipment").Preload("User").Find(&un_assigneds)

	ctx.HTML(http.StatusOK, "adminRequestings.html", gin.H{
		"heads":        []string{"序号", "申请人", "设备ID", "设备名", "型号", "创建时间", "操作"},
		"un_assigneds": un_assigneds,
		"total":        len(un_assigneds),
	})
}

func adminOngoings(ctx *gin.Context) {
	// token, _ := ctx.Cookie("token")
	// claim, _ := ParseToken(token)

	ongoing := []Ongoing{}
	db.Model(&Ongoing{}).Order("begin_at desc").Preload("Equipment").Preload("User").Find(&ongoing)

	ctx.HTML(http.StatusOK, "adminOngoings.html", gin.H{
		"heads":    []string{"序号", "联系人", "设备ID", "设备名", "型号", "开始时间", "操作"},
		"ongoings": ongoing,
		"total":    len(ongoing),
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
		ctx.Redirect(http.StatusPermanentRedirect, "/admin/assignUnits")
	case "finish":
		request := Request{
			ID:       requestID,
			Status:   FINISHED,
			EndAt:    time.Now(),
			EndAtStr: now(),
		}

		tx := db.Begin()
		// 更新request状态
		tx.Model(&request).Updates(&request).Take(&request, requestID)
		// 更新unit状态
		tx.Model(&EquipmentUnit{ID: request.EquipmentUnitID}).Updates(&EquipmentUnit{
			Availiable: true,
		})

		// 更新request池
		tx.Model(&Ongoing{}).Delete(&Ongoing{Request{ID: requestID}})
		err := tx.Commit().Error
		handle_resp(err, ctx)

		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "确认成功",
		})

	case "detail":
		request := Request{}
		db.Where(&Request{ID: requestID}).Preload("User").Preload("Equipment").Preload("EquipmentUnit").Take(&request)
		// log.Println(request)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"detail": request,
		})
	}
}

func adminAll(ctx *gin.Context) {
	requests := []Request{}
	db.Unscoped().Order("created_at desc").Preload("User").Preload("Equipment").Preload("EquipmentUnit").Find(&requests)

	ctx.HTML(http.StatusOK, "adminAll.html", gin.H{
		"heads":    []string{"记录号", "联系人", "设备ID", "设备名", "型号", "开始时间", "完成时间", "状态", "操作"},
		"requests": requests,
		"total":    len(requests),
	})
}
