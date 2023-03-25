package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func getAvailiable(ctx *gin.Context) {
	equipment := []Equipment{}
	//&Equipment{Availiable: true, Name: ctx.Query("name")}
	db.Model(&Equipment{}).Where("name LIKE ?", "%"+ctx.Query("name")+"%").Find(&equipment)
	// log.Println(ctx.Query("name"))
	// log.Println(equipment)

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
		ctx.Redirect(http.StatusTemporaryRedirect, "/user/login")
		return
	}
	// log.Println(ctx.Query("equipmentID"))
	// equipmentID, err := strconv.ParseUint(ctx.Query("equipmentID"), 10, 32)
	// handle(err)
	equipmentID := str2uint(ctx.Query("equipmentID"))
	request := Request{
		EquipmentID: uint(equipmentID),
		UserID:      claim.UserID,
		Status:      REQUESTING,
	}

	request.CreatedAtStr = time.Now().Format("2006-01-02 15:04")

	request_exist := UnAssigned{}

	db.Model(&UnAssigned{}).Where(&UnAssigned{request}).Take(&request_exist)
	if request_exist.ID != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "请勿重复申请",
		})
		return
	}

	err = db.Model(&Request{}).Create(&request).Error
	handle_resp(err, ctx)
	err = db.Model(&UnAssigned{}).Create(&UnAssigned{request}).Error
	handle_resp(err, ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "Success",
		"msg":    "请求完成",
	})

}

func getAvailiableUnits(ctx *gin.Context) {
	equipmentID := str2uint(ctx.Query("equipmentID"))

	units := EquipmentUnit{}
	db.Model(&EquipmentUnit{}).Where(&EquipmentUnit{EquipmentID: equipmentID}).Find(&units)
	ctx.HTML(http.StatusOK, "selectEquipmentUnit.html", gin.H{
		"units": units,
	})
}
