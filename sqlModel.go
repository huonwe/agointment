package main

import (
	"time"

	"gorm.io/gorm"
)

const (
	REQUESTING string = "请求中"
	ONGOING    string = "进行中"
	FINISHED   string = "已完成"

	TIMEOUT_REQUEST string = "请求超时"

	ASSIGNED string = "已指派"
	REJECTED string = "被拒绝"
	CANCELED string = "已取消"
)

type Department struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `gorm:"primarykey"`
	Description string
	Availiable  bool `gorm:"default:true"`
	// ...
}

type User struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Password string
	// 微信的openid
	WeChat string

	// DeptName     string
	DepartmentID uint
	Department   Department
	IsAdmin      bool `gorm:"default:false"`
	IsSuper      bool `gorm:"default:false"`

	Requests []Request
}

type APIUser struct {
	ID           uint
	Name         string
	DeptName     string
	DepartmentID uint
	Department   Department
}

type EquipmentUnit struct {
	ID          uint `gorm:"primarykey"`
	EquipmentID uint
	Type        string
	Class       string
	Name        string

	UID          string
	Brand        string
	SerialNumber string
	Price        string
	Label        string
	Factory      string
	Remark       string
	Status       string

	Availiable bool `gorm:"default:true"`
	// 佔用這個設備的人
	UserID uint
	// User   User
}

type Equipment struct {
	// gorm.Model
	ID uint `gorm:"primarykey"`
	// 设备名称
	Name string
	// 设备型号
	Type string
	// 设备分类
	Class string
	// 品牌
	Brand string

	Availiable bool `gorm:"default:true"`
}

type Request struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedAtStr string
	BeginAt      time.Time
	BeginAtStr   string
	EndAt        time.Time
	EndAtStr     string

	EquipmentID    uint
	EquipmentName  string
	EquipmentType  string
	EquipmentClass string
	EquipmentBrand string
	// Equipment      Equipment

	// 不是UID
	EquipmentUnitID  uint
	UnitUID          string
	UnitSerialNumber string
	UnitPrice        string
	UnitLabel        string
	UnitFactory      string
	UnitRemark       string
	UnitStatus       string

	UserID  uint
	User    User
	Status  string
	Deleted gorm.DeletedAt
}

// func (u *User) AfterFind(tx *gorm.DB) (err error) {
// 	dept := Department{}
// 	err = db.Take(&dept, u.DepartmentID).Error
// 	if err != nil {
// 		return err
// 	}
// 	u.DeptName = dept.Name
// 	return nil
// }

// func (u *User) AfterTake(tx *gorm.DB) (err error) {
// 	dept := Department{}
// 	err = db.Take(&dept, u.DepartmentID).Error
// 	if err != nil {
// 		return err
// 	}
// 	u.DeptName = dept.Name
// 	return nil
// }

// type UnAssigned struct {
// 	Request
// }

// type Ongoing struct {
// 	Request
// }

// func (req *Request) AfterUpdate(tx *gorm.DB) (err error) {
// 	// 维护 UnAssigned 和 Ongoing
// 	switch req.Status {
// 	case ONGOING:
// 		req_new := Request{}
// 		db.Take(&req_new, req.ID)
// 		err = tx.Unscoped().Delete(&UnAssigned{Request{ID: req.ID}}).Error
// 		if err != nil {
// 			return err
// 		}
// 		err = tx.Model(&Ongoing{}).Create(&Ongoing{req_new}).Error
// 		if err != nil {
// 			return err
// 		}
// 		return nil

// 	case CANCELED:
// 		err = tx.Unscoped().Delete(&UnAssigned{Request{ID: req.ID}}).Error
// 		return err

// 	case REJECTED:
// 		err = tx.Unscoped().Delete(&UnAssigned{Request{ID: req.ID}}).Error
// 		return err

// 	case FINISHED:
// 		err = tx.Unscoped().Delete(&Ongoing{Request{ID: req.ID}}).Error
// 		return err

// 	}

// 	return nil
// }

// func (req *Request) AfterCreate(tx *gorm.DB) (err error) {
// 	err = tx.Create(&UnAssigned{*req}).Error
// 	return err
// }

func (eu *EquipmentUnit) AfterCreate(tx *gorm.DB) (err error) {
	// 如果指定了equipmentID
	// log.Println("1 eu", eu.Name, eu.UID)
	// if eu.EquipmentID != 0 {
	// 	err = tx.Save(&Equipment{ID: eu.EquipmentID, Name: eu.Name, Type: eu.Type, Class: eu.Class, Brand: eu.Brand, Availiable: true}).Error
	// 	return
	// }
	// 如果没指定equipmentID 但 相同equipment已存在，则不需要添加
	exist_equipment := Equipment{}
	tx.Model(&Equipment{}).Where(&Equipment{Class: eu.Class, Name: eu.Name, Type: eu.Type, Brand: eu.Brand}).Take(&exist_equipment)
	if exist_equipment.ID > 0 {
		tx.Model(&eu).Update("equipment_id", exist_equipment.ID)
		return
	}
	// 如果没指定equipmentID 且 相同equipment不存在
	// log.Println("2 eu ID", eu.ID)
	new_equipment := Equipment{Name: eu.Name, Type: eu.Type, Class: eu.Class, Brand: eu.Brand, Availiable: true}
	err = tx.Create(&new_equipment).Error
	// log.Println("new equipment ID", new_equipment.ID)

	tx.Model(&eu).Update("equipment_id", new_equipment.ID)
	return
}

// func (eu *EquipmentUnit) AfterUpdate(tx *gorm.DB) (err error) {
// 	// log.Println("after update:", eu.EquipmentID)
// 	err = tx.Save(&Equipment{ID: eu.EquipmentID, Name: eu.Name, Type: eu.Type, Class: eu.Class, Brand: eu.Brand, Availiable: true}).Error
// 	return
// }

func initDB(db *gorm.DB) {
	// db.Exec("SET FOREIGN_KEY_CHECKS=0;")
	// db.Exec("DROP TABLE departments, users, equipment_units, equipment, requests, un_assigneds, ongoings")
	// db.Exec("SET FOREIGN_KEY_CHECKS=1;")

	db.AutoMigrate(&Department{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&EquipmentUnit{})
	db.AutoMigrate(&Equipment{})
	db.AutoMigrate(&Request{})
	// db.AutoMigrate(&UnAssigned{})
	// db.AutoMigrate(&Ongoing{})

	// db.Create(&Department{Name: "智医2002", Description: "智能医学工程", Availiable: true})
	// db.Create(&Department{Name: "智医2102", Description: "智能医学工程", Availiable: true})

	db.Create(&User{Name: "huonwe", Password: "Hhw20020120", DepartmentID: 1, IsAdmin: true, IsSuper: true})
	db.Create(&User{Name: "jimengxvan", Password: "jimengxvan", DepartmentID: 2, IsAdmin: true, IsSuper: true})

}
