package main

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func myRequest(ctx *gin.Context) {
	token, _ := ctx.Cookie("token")
	claim, _ := ParseToken(token)

	page := str2int(ctx.Query("page"))
	pageSize := str2int(ctx.Query("pageSize"))

	var total int64 = 0
	user := &User{
		ID: claim.UserID,
	}
	db.Model(&user).Preload("Requests", order_desc_createdAt, func(db *gorm.DB) *gorm.DB {
		return db.Where("equipment_name LIKE ?", "%"+ctx.Query("name")+"%").Limit(pageSize).Offset((page - 1) * pageSize)
	}).Take(&user)

	db.Model(&Request{}).Where("user_id = ?", claim.UserID).Where("equipment_name LIKE ?", "%"+ctx.Query("name")+"%").Count(&total)
	ctx.HTML(http.StatusOK, "statusTemplate.html", gin.H{
		"heads":      []string{"请求号", "设备名", "型号", "设备ID", "状态", "操作"},
		"requests":   user.Requests,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"total_page": int(math.Ceil(float64(total) / float64(pageSize))),
	})

}

func myRequestOp(ctx *gin.Context) {
	token, _ := ctx.Cookie("token")
	claim, _ := ParseToken(token)

	op := ctx.Query("op")
	request := Request{
		ID:     str2uint(ctx.Query("requestID")),
		UserID: claim.UserID,
	}
	switch op {
	case "cancel":
		// handle_resp(db.Model(&UnAssigned{}).Delete(&UnAssigned{request}).Error, ctx)
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

func changeDept(ctx *gin.Context) {
	user, have := ctx.Get("user")
	if !have {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	new_dept := ctx.PostForm("dept")

	dept_ := Department{}
	db.Model(&Department{}).Where(&Department{Name: new_dept}).First(&dept_)
	if dept_.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "部门修改失败,没有此部门",
		})
		return
	}

	user_ := User{}
	db.First(&user_, user.(User).ID)
	h := time.Until(user_.UpdatedAt).Abs().Hours()
	if h < 7*24 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    fmt.Sprintf("您在一周内只可修改一次个人信息\n距离下次可修改还剩%.1f小时", 7*24-h),
		})
		return
	}

	db.Model(&User{ID: user.(User).ID}).Update("department_id", dept_.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "部门修改成功",
	})
}

func changeName(ctx *gin.Context) {
	user, have := ctx.Get("user")
	if !have {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	new_name := ctx.PostForm("name")
	if len(new_name) < 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "Failed",
			"msg":    "用户名过短",
		})
		return
	}

	user_ := User{}
	db.First(&user_, user.(User).ID)
	h := time.Until(user_.UpdatedAt).Abs().Hours()
	if h < 7*24 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    fmt.Sprintf("您在一周内只可修改一次个人信息\n距离下次可修改还剩%.1f小时", 7*24-h),
		})
		return
	}

	db.Model(&User{ID: user.(User).ID}).Update("name", new_name)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "用户名修改成功",
	})
}

func changePasswd(ctx *gin.Context) {
	user, have := ctx.Get("user")
	if !have {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	if user.(User).Password != md5_str(ctx.PostForm("passwd_old")) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "Failed",
			"msg":    "密码错误",
		})
		return
	}

	new_passwd := ctx.PostForm("passwd_new")
	if len(new_passwd) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "Failed",
			"msg":    "密码过短",
		})
		return
	}

	// user_ := User{}
	// db.First(&user_, user.(User).ID)
	// h := time.Until(user_.UpdatedAt).Abs().Hours()
	// if h < 7*24 {
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"status": "Failed",
	// 		"msg":    fmt.Sprintf("您在一周内只可修改一次个人信息\n距离下次可修改还剩%.1f小时", 7*24-h),
	// 	})
	// 	return
	// }

	db.Model(&User{ID: user.(User).ID}).Update("password", md5_str(new_passwd))
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "密码修改成功",
	})
}
