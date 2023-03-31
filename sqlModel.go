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

	DepartmentName string
	Department     Department
	IsAdmin        bool `gorm:"default:false"`
	IsSuper        bool `gorm:"default:false"`

	Requests []Request
}

type EquipmentUnit struct {
	ID          uint `gorm:"primarykey"`
	EquipmentID uint
	Equipment   Equipment
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
	User   User
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

	EquipmentID   uint
	EquipmentName string
	Equipment     Equipment
	// 不是UID
	EquipmentUnitID  uint
	EquipmentUnitUID string
	EquipmentUnit    EquipmentUnit
	UserID           uint
	User             User
	Status           string // REQUESTING
	Deleted          gorm.DeletedAt
}

type UnAssigned struct {
	Request
}

type Ongoing struct {
	Request
}

// func (u *User) AfterFind(tx *gorm.DB) (err error) {
// 	dept := Department{}
// 	db.Model(&Department{Name: u.DepartmentName}).Take(&dept)
// 	if dept.Availiable {
// 		return
// 	}
// 	return errors.New("Department Unavailiable")
// }

// func (u *User) AfterTake(tx *gorm.DB) (err error) {
// 	dept := Department{}
// 	db.Model(&Department{Name: u.DepartmentName}).Take(&dept)
// 	if dept.Availiable {
// 		return
// 	}
// 	return errors.New("Department Unavailiable")
// }

// func (eu *EquipmentUnit) AfterCreate(tx *gorm.DB) (err error) {

// }

func (eu *EquipmentUnit) AfterCreate(tx *gorm.DB) (err error) {
	// 如果指定了equipmentID
	// log.Println("1 eu", eu.Name, eu.UID)
	if eu.EquipmentID != 0 {
		// log.Println("!!!!", eu.EquipmentID)
		err = tx.Save(&Equipment{ID: eu.EquipmentID, Name: eu.Name, Type: eu.Type, Class: eu.Class, Brand: eu.Brand, Availiable: true}).Error
		return
	}
	// 如果没指定equipmentID 但 相同equipment已存在，则不需要添加
	exist_equipment := Equipment{}
	tx.Model(&Equipment{}).Where(&Equipment{Class: eu.Class, Name: eu.Name, Type: eu.Type}).Take(&exist_equipment)
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
// 	log.Println("after update:", eu.EquipmentID)
// 	err = tx.Save(&Equipment{ID: eu.EquipmentID, Name: eu.Name, Type: eu.Type, Class: eu.Class, Brand: eu.Brand, Availiable: true}).Error
// 	return
// }

// func (eu *EquipmentUnit) AfterCreate(tx *gorm.DB) (err error) {
// 	log.Println("after update:", eu.EquipmentID)
// 	err = tx.Save(&Equipment{ID: eu.EquipmentID, Name: eu.Name, Type: eu.Type, Class: eu.Class, Brand: eu.Brand, Availiable: true}).Error
// 	return
// }

func initDB(db *gorm.DB) {
	db.Exec("DROP TABLE departments")
	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE equipment_units")

	db.Exec("DROP TABLE equipment") // equipment is uncountable

	db.Exec("DROP TABLE requests")
	db.Exec("DROP TABLE un_assigneds")
	db.Exec("DROP TABLE ongoings")

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Department{})
	db.AutoMigrate(&EquipmentUnit{})
	db.AutoMigrate(&Equipment{})
	db.AutoMigrate(&Request{})
	db.AutoMigrate(&UnAssigned{})
	db.AutoMigrate(&Ongoing{})

	db.Create(&Department{Name: "测试部门", Description: "测试..."})

	// db.Create(&User{Name: "admin", Password: "202303311700", DepartmentName: "测试部门", IsAdmin: true})

	db.Create(&User{Name: "huonwe", Password: "Hhw20020120", DepartmentName: "测试部门", IsAdmin: true})
	db.Create(&User{Name: "jimengxvan", Password: "jimengxvan", DepartmentName: "测试部门", IsAdmin: true})

}
