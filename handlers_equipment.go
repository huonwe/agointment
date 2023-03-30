package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func getAvailiable(ctx *gin.Context) {
	equipment := []Equipment{}
	db.Model(&Equipment{}).Where("name LIKE ?", "%"+ctx.Query("name")+"%").Find(&equipment)
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
	log.Println(file.Filename)
	ctx.SaveUploadedFile(file, "./static/files/equipment.xlsx")

	f, err := excelize.OpenFile("./static/files/equipment.xlsx")
	handle_resp(err, ctx)
	rows, _ := f.GetRows("Sheet1")

	var units []EquipmentUnit
	for index, row := range rows {
		if index == 0 {
			continue
		}
		// log.Println(row)

		unit := EquipmentUnit{
			Class:        row[0],
			Name:         row[1],
			Type:         row[2],
			Brand:        row[3],
			Factory:      row[4],
			Price:        row[5],
			ID:           row[6],
			SerialNumber: row[7],
			Label:        row[8],
			Status:       row[9],
			Remark:       row[10],
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
