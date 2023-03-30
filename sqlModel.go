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

	Requests []Request
}

type EquipmentUnit struct {
	EquipmentID uint
	Equipment   Equipment
	Type        string
	Class       string
	Name        string

	ID           string `gorm:"primarykey"`
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
	ID uint `gorm:"primarykey auto_increment:true"`
	// 设备名称
	Name string `gorm:"primarykey auto_increment:false"`
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

	EquipmentID     uint
	EquipmentName   string
	Equipment       Equipment
	EquipmentUnitID string
	EquipmentUnit   EquipmentUnit
	UserID          uint
	User            User
	Status          string // REQUESTING
	Deleted         gorm.DeletedAt
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

func (eu *EquipmentUnit) AfterCreate(tx *gorm.DB) (err error) {
	var count int64
	tx.Model(&Equipment{}).Where(&Equipment{Class: eu.Class, Name: eu.Name, Type: eu.Type}).Count(&count)
	if count > 0 {
		return
	}
	err = tx.Model(&Equipment{}).Create(&Equipment{Name: eu.Name, Type: eu.Type, Class: eu.Class}).Error
	return
}

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

	db.Create(&Department{Name: "德国骨科", Description: "德国骨科..."})
	db.Create(&Department{Name: "测试部门", Description: "测试..."})

	db.Create(&User{Name: "test", Password: "123456", DepartmentName: "德国骨科", IsAdmin: false})
	db.Create(&User{Name: "huonwe", Password: "huonwe", DepartmentName: "测试部门", IsAdmin: true})

	db.Create(&Equipment{Name: "测试设备", Type: "试做型", Class: "醫用設備", Availiable: true})
	db.Create(&Equipment{Name: "测试设备", Type: "试做型", Class: "未来科技", Availiable: true})
	db.Create(&Equipment{Name: "空想具现", Type: "试做型", Class: "宏伟制造", Availiable: true})

	db.Create(&EquipmentUnit{Name: "测试设备", Type: "试做型", Class: "醫用設備", ID: "0001", Brand: "宏偉製造", SerialNumber: "001", Price: "999.9", Label: "沒有標註", Factory: "宏偉天津製造工廠", Availiable: true})
	db.Create(&EquipmentUnit{Name: "测试设备", Type: "试做型", Class: "醫用設備", ID: "0002", Brand: "宏偉製造", SerialNumber: "001", Price: "999.9", Label: "沒有標註", Factory: "宏偉天津製造工廠", Availiable: true})

}
