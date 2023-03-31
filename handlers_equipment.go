package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func getAvailiable(ctx *gin.Context) {
	equipment := []Equipment{}
	db.Model(&Equipment{}).Where(&Equipment{Availiable: true}).Where("name LIKE ?", "%"+ctx.Query("name")+"%").Find(&equipment)
	ctx.HTML(http.StatusOK, "availableList.html", gin.H{
		"heads":      []string{"序号", "设备名", "型号", "类别", "操作"},
		"equipments": equipment,
		"total":      len(equipment),
	})
}

func equipmentRequest(ctx *gin.Context) {
	value, exist := ctx.Get("user")
	user, ok := value.(User)
	if !exist || !ok {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	equipmentID := str2uint(ctx.Query("equipmentID"))
	equipment := Equipment{}
	db.Take(&equipment, equipmentID)
	request := Request{
		EquipmentID:   equipment.ID,
		EquipmentName: equipment.Name,
		UserID:        user.ID,
	}

	var count int64
	db.Model(&UnAssigned{}).Where(&UnAssigned{request}).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "请勿重复申请",
		})
		return
	}
	request.CreatedAtStr = now()
	request.Status = REQUESTING

	err := db.Model(&Request{}).Create(&request).Error
	handle_resp(err, ctx)
	err = db.Model(&UnAssigned{}).Create(&UnAssigned{request}).Error
	handle_resp(err, ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "请求完成",
	})

}

// assignUnits GET
func getAvailiableUnits(ctx *gin.Context) {
	equipmentID := str2uint(ctx.Query("equipmentID"))

	equipment := Equipment{}
	db.Take(&equipment, equipmentID)

	units := []EquipmentUnit{}
	db.Model(&EquipmentUnit{}).Order("availiable desc").Where(&EquipmentUnit{Name: equipment.Name}).Find(&units)
	// for _, unit := range units {
	// 	log.Println(unit.Availiable)
	// }
	ctx.HTML(http.StatusOK, "adminAssignUnit.html", gin.H{
		"units": units,
		"total": len(units),
	})
}

// assignUnits POST
func assignUnits(ctx *gin.Context) {
	if ctx.PostForm("unitID") == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "未选择设备实体",
		})
		return
	}
	unitID := ctx.PostForm("unitID")
	requestID := str2uint(ctx.PostForm("requestID"))
	equipmentID := str2uint(ctx.PostForm("equipmentID"))
	requestorID := str2uint(ctx.PostForm("requestorID"))
	// log.Println(unitID)
	unit := EquipmentUnit{
		ID:          unitID,
		EquipmentID: equipmentID,
		UserID:      requestorID,
		Availiable:  false, // 0值, updates不会更新这条
	}

	tx := db.Begin()
	// 更新实体所有者
	err := tx.Model(&unit).Updates(&unit).Error
	if err != nil {
		tx.Rollback()
		handle_resp(err, ctx)
	}
	// 单独更新Availiable
	err = tx.Model(&unit).Update("Availiable", false).Error
	if err != nil {
		tx.Rollback()
		handle_resp(err, ctx)
	}

	request_new_state := Request{
		ID:              requestID,
		BeginAt:         time.Now(),
		BeginAtStr:      now(),
		EquipmentUnitID: unitID,
		Status:          ONGOING,
	}
	// 更新Request状态
	err = tx.Model(&Request{ID: requestID}).Updates(&request_new_state).Take(&request_new_state, requestID).Error
	if err != nil {
		tx.Rollback()
		handle_resp(err, ctx)
	}
	// 更新Request池
	err = tx.Delete(&UnAssigned{Request{ID: requestID}}).Error
	if err != nil {
		tx.Rollback()
		handle_resp(err, ctx)
	}
	err = tx.Create(&Ongoing{request_new_state}).Error
	if err != nil {
		tx.Rollback()
		handle_resp(err, ctx)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		handle_resp(err, ctx)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "分配成功",
	})
}

func equipmentImport(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	if !strings.HasSuffix(file.Filename, ".xlsx") {
		ctx.JSON(http.StatusOK, gin.H{
			"static": "Failed",
			"msg":    "文件格式错误",
		})
		return
	}
	log.Println(file.Filename)
	ctx.SaveUploadedFile(file, "./static/files/equipment.xlsx")

	f, err := excelize.OpenFile("./static/files/equipment.xlsx")
	handle_resp(err, ctx)
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	rows, _ := f.GetRows("Sheet1")

	var units []EquipmentUnit
	for index, row := range rows {
		if index == 0 {
			continue
		}
		// log.Println(row)

		unit := EquipmentUnit{
			Class:        row[1],
			Name:         row[2],
			Type:         row[3],
			Brand:        row[4],
			Factory:      row[5],
			Price:        row[6],
			ID:           row[7],
			SerialNumber: row[8],
			Label:        row[9],
			Status:       row[10],
			Remark:       row[11],
		}
		if row[0] != "" {
			unit.EquipmentID = str2uint(row[0])
		}
		units = append(units, unit)

	}
	err = db.Save(&units).Error
	handle_resp(err, ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"stauts": "Success",
		"msg":    "成功",
	})
}

func equipmentOp(ctx *gin.Context) {
	op := ctx.Query("op")
	id := ctx.Query("id")
	if op == "" || id == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"stauts": "Failed",
			"msg":    "参数错误",
		})
		return
	}
	ID := str2uint(id)
	switch op {
	case "del":
		err := db.Model(&Equipment{ID: ID}).Where(&Equipment{ID: ID}).Delete(ID).Error
		handle_resp(err, ctx)

		ctx.JSON(http.StatusOK, gin.H{
			"stauts": "Success",
			"msg":    "删除成功",
		})
	case "enable":
		db.Model(&Equipment{ID: ID}).Where(&Equipment{ID: ID}).Update("availiable", true)
		ctx.JSON(http.StatusOK, gin.H{
			"stauts": "Success",
			"msg":    "启用成功",
		})
	case "disable":
		db.Model(&Equipment{ID: ID}).Where(&Equipment{ID: ID}).Update("availiable", false)
		ctx.JSON(http.StatusOK, gin.H{
			"stauts": "Success",
			"msg":    "禁用成功",
		})
	}
}
