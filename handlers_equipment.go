package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func getAvailiable(ctx *gin.Context) {
	page := str2int(ctx.Query("page"))
	pageSize := str2int(ctx.Query("pageSize"))
	var total int64 = 0
	equipment := []Equipment{}
	db.Model(&Equipment{}).Where(&Equipment{Availiable: true}).Where("name LIKE ?", "%"+ctx.Query("name")+"%").Count(&total)
	db.Model(&Equipment{}).Where(&Equipment{Availiable: true}).Where("name LIKE ?", "%"+ctx.Query("name")+"%").Find(&equipment).Limit(pageSize).Offset((page - 1) * pageSize)
	ctx.HTML(http.StatusOK, "availableList.html", gin.H{
		"heads":      []string{"序号", "设备名", "型号", "类别", "操作"},
		"equipments": equipment,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"total_page": int(math.Ceil(float64(total) / float64(pageSize))),
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
		EquipmentID:    equipment.ID,
		EquipmentName:  equipment.Name,
		EquipmentType:  equipment.Type,
		EquipmentClass: equipment.Class,
		EquipmentBrand: equipment.Brand,
		UserID:         user.ID,
		// UserName:       user.Name,
		// UserDeptName:   user.DeptName,
	}

	var count int64
	db.Model(&Request{}).Where("status = ?", REQUESTING).Where(&request).Count(&count)
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

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "请求完成",
	})

}

// assignUnits GET
func getAvailiableUnits(ctx *gin.Context) {
	equipmentID := str2uint(ctx.Query("equipmentID"))
	units := []EquipmentUnit{}
	db.Model(&EquipmentUnit{}).Order("availiable desc").Where(&EquipmentUnit{EquipmentID: equipmentID}).Find(&units)
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
	err := tx.Model(&unit).Updates(&unit).Take(&unit).Error
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
		UnitUID:         unit.UID,

		UnitSerialNumber: unit.SerialNumber,
		UnitPrice:        unit.Price,
		UnitLabel:        unit.Label,
		UnitFactory:      unit.Factory,
		UnitRemark:       unit.Remark,
		UnitStatus:       unit.Status,

		UserID: requestorID,

		Status: ONGOING,
	}
	// 更新Request状态
	err = tx.Model(&Request{ID: requestID}).Updates(&request_new_state).Take(&request_new_state, requestID).Error
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

func adminEquipmentImport(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	if !strings.HasSuffix(file.Filename, ".xlsx") {
		ctx.JSON(http.StatusOK, gin.H{
			"static": "Failed",
			"msg":    "文件格式错误",
		})
		return
	}
	// log.Println(file.Filename)
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

	var units_new []EquipmentUnit
	// var units_exist []EquipmentUnit

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
			UID:          row[7],
			SerialNumber: row[8],
			Label:        row[9],
			Status:       row[10],
			Remark:       row[11],
		}
		if row[0] != "" {
			eid, err := str2uint_err(row[0])
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"status": "Failed",
					"msg":    fmt.Sprintf("警告, 第%d行EID %s 有误", index+1, row[0]),
				})
				return
			}
			unit.EquipmentID = eid
			equipment_exist := Equipment{}
			db.Model(&Equipment{}).Take(&equipment_exist, eid)
			equipment_new := Equipment{Name: unit.Name, Class: unit.Class, Type: unit.Type, Brand: unit.Brand}
			if equipment_exist.ID != 0 {
				db.Model(&Equipment{ID: equipment_exist.ID}).Updates(&equipment_new)
			} else {
				db.Create(&equipment_new)
			}
		}

		var count int64
		db.Model(&EquipmentUnit{}).Where(&EquipmentUnit{UID: unit.UID, SerialNumber: unit.SerialNumber}).Count(&count)
		if count == 0 {
			units_new = append(units_new, unit)
		} else {
			// unit_exist := EquipmentUnit{}
			// db.Model(&EquipmentUnit{}).Where(&EquipmentUnit{UID: unit.UID, SerialNumber: unit.SerialNumber}).Take(&unit_exist)
			// if unit.EquipmentID != 0 && unit.EquipmentID != unit_exist.ID {

			// }
			db.Model(&EquipmentUnit{}).Where(&EquipmentUnit{UID: unit.UID, SerialNumber: unit.SerialNumber}).Updates(&unit)
		}
	}
	db.Create(&units_new)

	handle_resp(err, ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "成功",
	})
}

func equipmentOp(ctx *gin.Context) {
	op := ctx.Query("op")
	id := ctx.Query("id")
	if op == "" || id == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "参数错误",
		})
		return
	}
	ID := str2uint(id)
	switch op {
	case "del":
		err := db.Model(&Equipment{ID: ID}).Where(&Equipment{ID: ID}).Delete(ID).Error
		handle_resp(err, ctx)

		db.Model(&EquipmentUnit{}).Where(&EquipmentUnit{EquipmentID: ID}).Delete(&EquipmentUnit{})

		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "删除成功",
		})
	case "enable":
		db.Model(&Equipment{ID: ID}).Where(&Equipment{ID: ID}).Update("availiable", true)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "启用成功",
		})
	case "disable":
		db.Model(&Equipment{ID: ID}).Where(&Equipment{ID: ID}).Update("availiable", false)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Success",
			"msg":    "禁用成功",
		})
	}
}

func equipmentTemplate(ctx *gin.Context) {
	ctx.Redirect(http.StatusTemporaryRedirect, "/static/template/template.xlsx")
}
