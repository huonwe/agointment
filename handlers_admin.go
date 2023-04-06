package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func adminRequestings(ctx *gin.Context) {
	un_assigneds := []Request{}
	db.Model(&Request{}).Where("status = ?", REQUESTING).Order("created_at desc").Joins("User").Preload("User.Department").Find(&un_assigneds)

	ctx.HTML(http.StatusOK, "adminRequestings.html", gin.H{
		"heads":        []string{"序号", "申请人", "设备名", "型号", "创建时间", "操作"},
		"un_assigneds": un_assigneds,
		"total":        len(un_assigneds),
	})
}

func adminOngoings(ctx *gin.Context) {
	ongoings := []Request{}
	db.Model(&Request{}).Where("status = ?", ONGOING).Order("begin_at desc").Joins("User").Preload("User.Department").Find(&ongoings)
	ctx.HTML(http.StatusOK, "adminOngoings.html", gin.H{
		"heads":    []string{"序号", "联系人", "设备ID", "设备名", "型号", "开始时间", "操作"},
		"ongoings": ongoings,
		"total":    len(ongoings),
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
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "",
		})
	case "assign":
		ctx.JSON(http.StatusBadGateway, nil)
		// ctx.Redirect(http.StatusPermanentRedirect, "/admin/assignUnits")
	case "finish":
		request := Request{
			ID:       requestID,
			Status:   FINISHED,
			EndAt:    time.Now(),
			EndAtStr: now(),
		}

		tx := db.Begin()
		// 更新request状态
		tx.Model(&request).Updates(&request)
		tx.Model(&Request{}).Take(&request, requestID)
		// 更新unit状态
		tx.Model(&EquipmentUnit{ID: request.EquipmentUnitID}).Update("availiable", true)

		err := tx.Commit().Error
		handle_resp(err, ctx)

		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "确认成功",
		})

	case "detail":
		request := Request{}
		db.Unscoped().Where(&Request{ID: requestID}).Joins("User").Preload("User.Department").Take(&request)
		request.User.Password = ""
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
	db.Unscoped().Model(&Request{}).Order("created_at desc").Joins("User").Preload("User.Department").Limit(pageSize).Offset((page - 1) * pageSize).Find(&requests)
	db.Unscoped().Model(&Request{}).Count(&total)

	ctx.HTML(http.StatusOK, "adminAll.html", gin.H{
		"heads":      []string{"联系人", "设备ID", "设备名", "开始时间", "完成时间", "状态", "操作"},
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
		db.Model(&Department{}).Where(&Department{Name: deptName}).Delete(&Department{Name: deptName})
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "删除成功",
		})
	case "deptBan":
		dept := Department{
			Name:       deptName,
			Availiable: false,
		}
		db.Model(&dept).Where(&Department{Name: dept.Name}).Select("Name", "Availiable").Updates(&dept)
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

		dept := Department{}
		db.Model(&Department{}).Where("name = ?", userDept).Take(&dept)
		if dept.ID == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "Failed",
				"msg":    "部门不存在",
			})
			return
		}

		user := User{
			Name:         userName,
			DepartmentID: dept.ID,
		}
		var count int64
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

		// log.Println(user)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "创建成功",
		})
	case "userSearch": // HTML
		name := ctx.Query("name")
		dept := ctx.Query("dept")
		users := []User{}

		db.Model(&User{}).InnerJoins("Department", db.Where(&Department{Name: dept})).Where("users.name LIKE ?", "%"+name+"%").Find(&users)
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
	case "userSetAdmin":
		id := str2uint(ctx.PostForm("id"))
		db.Model(&User{ID: id}).Update("is_admin", true)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "设置管理员成功",
		})
	case "userUnsetAdmin":
		id := str2uint(ctx.PostForm("id"))
		db.Model(&User{ID: id}).Update("is_admin", false)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "取消管理员成功",
		})
	}
}

func adminEquipment(ctx *gin.Context) {
	page, err := str2int_err(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := str2int_err(ctx.Query("pageSize"))
	if err != nil {
		pageSize = 20
	}
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

func adminExportRequests(ctx *gin.Context) {
	f := excelize.NewFile()

	requests := []Request{}
	db.Unscoped().Model(&Request{}).Joins("User").Preload("User.Department").Find(&requests)
	//
	// for _, request := range requests {
	// 	log.Println(request.User.DeptName)
	// }
	//
	heads := []string{"请求ID", "请求人", "请求人部门", "创建时间", "开始时间", "结束时间", "请求状态", "EID", "设备类别", "设备名", "设备型号", "设备品牌", "设备实体ID", "设备序列号", "工厂", "标签", "状态", "备注"}
	for i, v := range heads {
		f.SetCellValue("Sheet1", Num2Col(i)+"1", v)
	}
	for i, request := range requests {
		rowIdx := strconv.FormatInt(int64(i)+2, 10)
		f.SetCellValue("Sheet1", "A"+rowIdx, request.ID)
		f.SetCellValue("Sheet1", "B"+rowIdx, request.User.Name)
		f.SetCellValue("Sheet1", "C"+rowIdx, request.User.Department.Name)
		f.SetCellValue("Sheet1", "D"+rowIdx, request.CreatedAt)
		f.SetCellValue("Sheet1", "E"+rowIdx, request.BeginAt)
		f.SetCellValue("Sheet1", "F"+rowIdx, request.EndAt)
		f.SetCellValue("Sheet1", "G"+rowIdx, request.Status)

		f.SetCellValue("Sheet1", "H"+rowIdx, request.EquipmentID)
		f.SetCellValue("Sheet1", "I"+rowIdx, request.EquipmentClass)
		f.SetCellValue("Sheet1", "J"+rowIdx, request.EquipmentName)
		f.SetCellValue("Sheet1", "K"+rowIdx, request.EquipmentType)
		f.SetCellValue("Sheet1", "L"+rowIdx, request.EquipmentBrand)

		f.SetCellValue("Sheet1", "M"+rowIdx, request.UnitUID)
		f.SetCellValue("Sheet1", "N"+rowIdx, request.UnitSerialNumber)
		f.SetCellValue("Sheet1", "O"+rowIdx, request.UnitFactory)
		f.SetCellValue("Sheet1", "P"+rowIdx, request.UnitLabel)
		f.SetCellValue("Sheet1", "Q"+rowIdx, request.UnitStatus)
		f.SetCellValue("Sheet1", "R"+rowIdx, request.UnitRemark)
	}

	filename := "export_requests_" + time.Now().Format("20060102150405") + ".xlsx"
	err := f.SaveAs("./static/files/" + filename)
	handle_resp(err, ctx)
	ctx.Redirect(http.StatusTemporaryRedirect, "/static/files/"+filename)
}

func adminExportEquipment(ctx *gin.Context) {
	f := excelize.NewFile()

	units := []EquipmentUnit{}
	db.Model(&EquipmentUnit{}).Order("equipment_id asc").Find(&units)
	heads := []string{"EID", "类别", "设备名", "型号", "品牌", "工厂", "价格", "设备实体ID", "序列号", "标签", "状态", "备注"}
	for i, v := range heads {
		f.SetCellValue("Sheet1", Num2Col(i)+"1", v)
		// log.Println(Num2Col(i))
	}
	for i, unit := range units {
		rowIdx := strconv.FormatInt(int64(i)+2, 10)
		f.SetCellValue("Sheet1", "A"+rowIdx, unit.EquipmentID)
		f.SetCellValue("Sheet1", "B"+rowIdx, unit.Class)
		f.SetCellValue("Sheet1", "C"+rowIdx, unit.Name)
		f.SetCellValue("Sheet1", "D"+rowIdx, unit.Type)
		f.SetCellValue("Sheet1", "E"+rowIdx, unit.Brand)
		f.SetCellValue("Sheet1", "F"+rowIdx, unit.Factory)
		f.SetCellValue("Sheet1", "G"+rowIdx, unit.Price)

		f.SetCellValue("Sheet1", "H"+rowIdx, unit.UID)
		f.SetCellValue("Sheet1", "I"+rowIdx, unit.SerialNumber)
		f.SetCellValue("Sheet1", "J"+rowIdx, unit.Label)
		f.SetCellValue("Sheet1", "K"+rowIdx, unit.Status)
		f.SetCellValue("Sheet1", "L"+rowIdx, unit.Remark)
	}
	filename := "export_equipment_" + time.Now().Format("20060102150405") + ".xlsx"
	err := f.SaveAs("./static/files/" + filename)
	handle_resp(err, ctx)
	ctx.Redirect(http.StatusTemporaryRedirect, "/static/files/"+filename)
}

func emptyEnded(ctx *gin.Context) {
	err := db.Unscoped().Model(&Request{}).Where("status IN ?", []string{CANCELED, FINISHED, REJECTED}).Delete(&Request{}).Error
	handle_resp(err, ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "清除成功",
	})
}

func maintain(ctx *gin.Context) {
	pageSize, err := str2int_err(ctx.Query("pageSize"))
	if err != nil {
		pageSize = 20
	}
	page, err := str2int_err(ctx.Query("page"))
	if err != nil {
		page = 1
	}

	units := []UnitMaintainAPI{}
	var count int64
	db.Model(&EquipmentUnit{}).Count(&count)
	db.Model(&EquipmentUnit{}).Order("last_maintained asc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&units)
	for i, unit := range units {
		unit.LastMaintainedStr = unit.LastMaintained.Format("2006-01-02 15:04")
		units[i] = unit
	}
	ctx.HTML(http.StatusOK, "maintain.html", gin.H{
		"heads":      []string{"设备UID", "设备名", "上次维护", "操作"},
		"units":      units,
		"total":      count,
		"page":       page,
		"pageSize":   pageSize,
		"total_page": int(math.Ceil(float64(count) / float64(pageSize))),
	})
}

func doMaintain(ctx *gin.Context) {
	id := ctx.PostForm("id")
	ID, err := str2uint_err(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	now := time.Now()

	unit := EquipmentUnit{}
	db.First(&unit, ID)
	// log.Println(time.Duration(now.Sub(unit.LastMaintained)))
	// log.Println(time.Duration(time.Duration.Minutes(5)))
	if time.Duration(now.Sub(unit.LastMaintained)) < 5*time.Minute {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "Failed",
			"msg":    "您在之前5分钟内已经维护过这个设备啦, 休息一下吧~",
		})
		return
	}

	db.Model(&EquipmentUnit{ID: ID}).Update("last_maintained", now.Local())
	maintain := Maintain{
		UnitID:   ID,
		DoAt:     now,
		DoAtStr:  now.Format("2006-01-02 15:04"),
		UID:      unit.UID,
		UnitName: unit.Name,
		Type:     unit.Type,
		Serial:   unit.SerialNumber,
	}
	db.Model(&Maintain{}).Create(&maintain)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "维护成功",
	})
}

func exportMaintain(ctx *gin.Context) {
	f := excelize.NewFile()
	heads := []string{"维护时间", "设备名", "型号", "设备实体ID", "序列号"}
	for i, v := range heads {
		f.SetCellValue("Sheet1", Num2Col(i)+"1", v)
	}

	maintains := []Maintain{}
	db.Model(&Maintain{}).Order("do_at desc").Find(&maintains)
	for i, maintain := range maintains {
		rowIdx := strconv.FormatInt(int64(i)+2, 10)
		f.SetCellValue("Sheet1", "A"+rowIdx, maintain.DoAt.Format("2006-01-02 15:04"))
		f.SetCellValue("Sheet1", "B"+rowIdx, maintain.UnitName)
		f.SetCellValue("Sheet1", "C"+rowIdx, maintain.Type)
		f.SetCellValue("Sheet1", "D"+rowIdx, maintain.UID)
		f.SetCellValue("Sheet1", "E"+rowIdx, maintain.Serial)
	}

	filename := "export_" + time.Now().Format("20060102150405") + ".xlsx"
	err := f.SaveAs("./static/files/" + filename)
	handle_resp(err, ctx)
	ctx.Redirect(http.StatusTemporaryRedirect, "/static/files/"+filename)
}

func attentionUpdate(ctx *gin.Context) {
	content := ctx.PostForm("content")
	db.Model(&Attention{}).Create(&Attention{Content: content})
	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "提交成功",
	})
}
