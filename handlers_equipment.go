package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	token, _ := ctx.Cookie("token")
	claim, err := ParseToken(token)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	// log.Println(ctx.Query("equipmentID"))
	// equipmentID, err := strconv.ParseUint(ctx.Query("equipmentID"), 10, 32)
	// handle(err)
	equipmentID := str2uint(ctx.Query("equipmentID"))
	equipment := Equipment{}
	db.Take(&equipment, equipmentID)
	request := Request{
		EquipmentID:   equipment.ID,
		EquipmentName: equipment.Name,
		UserID:        claim.UserID,
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

	err = db.Model(&Request{}).Create(&request).Error
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
	unitID := str2uint(ctx.PostForm("unitID"))
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
