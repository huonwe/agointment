package main

import (
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func adminRequestings(ctx *gin.Context) {
	// token, _ := ctx.Cookie("token")
	// claim, _ := ParseToken(token)

	un_assigneds := []UnAssigned{}
	db.Model(&UnAssigned{}).Order("created_at desc").Preload("Equipment").Preload("User").Find(&un_assigneds)

	ctx.HTML(http.StatusOK, "adminRequestings.html", gin.H{
		"heads":        []string{"序号", "申请人", "设备名", "型号", "创建时间", "操作"},
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
	page := str2int(ctx.Query("page"))
	pageSize := str2int(ctx.Query("pageSize"))
	var total int64 = 0

	requests := []Request{}
	db.Unscoped().Model(&Request{}).Order("created_at desc").Preload("User").Preload("Equipment").Preload("EquipmentUnit").Limit(pageSize).Offset((page - 1) * pageSize).Find(&requests)
	db.Unscoped().Model(&Request{}).Count(&total)

	ctx.HTML(http.StatusOK, "adminAll.html", gin.H{
		"heads":      []string{"记录号", "联系人", "设备ID", "设备名", "型号", "开始时间", "完成时间", "状态", "操作"},
		"requests":   requests,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"total_page": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func adminUsers(ctx *gin.Context) {
	depts := []Department{}
	db.Find(&depts)

	ctx.HTML(http.StatusOK, "adminUsers.html", gin.H{
		"depts": depts,
	})
}

func adminUsersOp(ctx *gin.Context) {
	op := ctx.Param("op")
	deptName := ctx.PostForm("deptName")
	if op == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "参数错误",
		})
		return
	}
	if strings.Contains(op, "dept") {
		if deptName == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "Failed",
				"msg":    "参数错误",
			})
			return
		}
	}
	switch op {
	case "deptDel":
		db.Delete(&Department{Name: deptName})
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "删除成功",
		})
	case "deptBan":
		dept := Department{
			Name:       deptName,
			Availiable: false,
		}
		db.Where(&Department{Name: dept.Name}).Select("Name", "Availiable").Updates(&dept)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "禁用成功",
		})
	case "deptActive":
		dept := Department{
			Name:       deptName,
			Availiable: true,
		}
		// log.Println(deptName)
		db.Where(&Department{Name: dept.Name}).Select("Name", "Availiable").Updates(dept)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "启用成功",
		})

	case "deptNew":
		dept := Department{
			Name:        deptName,
			Description: ctx.PostForm("deptDescpt"),
		}
		var count int64
		db.Model(&dept).Where("name = ?", deptName).Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "Failed",
				"msg":    "部门名重复",
			})
			return
		}
		err := db.Model(&Department{}).Create(&dept).Error
		if err != nil {
			handle_resp(err, ctx)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "创建成功",
		})
	case "userNew":
		userName := ctx.PostForm("user_name")
		userDept := ctx.PostForm("user_dept")
		userPassword := ctx.PostForm("user_password")
		if userDept == "" || userName == "" || userPassword == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "Failed",
				"msg":    "参数错误",
			})
			return
		}

		var count int64
		db.Model(&Department{}).Where("name = ?", userDept).Count(&count)
		if count == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "Failed",
				"msg":    "部门不存在",
			})
			return
		}

		user := User{
			Name:           userName,
			DepartmentName: userDept,
		}
		db.Model(&user).Where(&user).Count(&count)
		if count > 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "Failed",
				"msg":    "用户重复",
			})
			return
		}
		user.Password = userPassword
		err := db.Model(&user).Create(&user).Error
		if err != nil {
			handle_resp(err, ctx)
		}

		log.Println(user)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "创建成功",
		})
	case "userSearch": // HTML
		name := ctx.Query("name")
		dept := ctx.Query("dept")
		users := []User{}
		db.Model(&User{}).Where("name LIKE ? AND department_name LIKE ?", "%"+name+"%", "%"+dept+"%").Find(&users)
		ctx.HTML(http.StatusOK, "adminUsersUsers.html", gin.H{
			"users": users,
		})

	case "userDel":
		// name := ctx.PostForm("name")
		id := str2uint(ctx.PostForm("id"))
		// newPasswd := ctx.PostForm("newPasswd")
		user := User{ID: id}
		err := db.Model(&User{}).Delete(&user).Error
		handle_resp(err, ctx)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "删除用户成功",
		})
	case "userPasswd":
		// name := ctx.PostForm("name")
		id := str2uint(ctx.PostForm("id"))
		newPasswd := ctx.PostForm("newPasswd")
		user := User{ID: id, Password: newPasswd}
		err := db.Model(&user).Updates(&user).Error
		handle_resp(err, ctx)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "更新密码成功",
		})
	}
}

func adminEquipment(ctx *gin.Context) {
	if ctx.Query("page") == "" || ctx.Query("pageSize") == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	page := str2int(ctx.Query("page"))
	pageSize := str2int(ctx.Query("pageSize"))
	var total int64 = 0

	equipments := []Equipment{}
	db.Model(&Equipment{}).Order("id desc").Where("name LIKE ?", "%"+ctx.Query("name")+"%").Limit(pageSize).Offset((page - 1) * pageSize).Find(&equipments)
	db.Model(&Equipment{}).Where("name LIKE ?", "%"+ctx.Query("name")+"%").Count(&total)

	ctx.HTML(http.StatusOK, "adminEquipment.html", gin.H{
		"heads":      []string{"EID", "设备名", "设备型号", "设备分类", "是否可用", "操作"},
		"equipments": equipments,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"total_page": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}
