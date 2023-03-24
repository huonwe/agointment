package main

import (
	"time"

	"gorm.io/gorm"
)

const (
	REQUESTING string = "请求中"
	TIMING     string = "进行中"
	COMPLETED  string = "已完成"

	TIMEOUT_REQUEST string = "请求超时"

	ASSIGNED string = "已指派"
	REJECTED string = "被拒绝"
)

type Department struct {
	Name        string `gorm:"primarykey"`
	Description string

	Users []User
	// ...
}

type User struct {
	ID             uint `gorm:"primarykey"`
	Name           string
	Password       string
	DepartmentName string
	Department     Department
	IsAdmin        bool

	Requests []Request
}

type EquipmentUnit struct {
	EquipmentName string
	EquipmentID   uint
	Equipment     Equipment
	// Name         string
	// Type         string
	// Class   string
	ID           uint `gorm:"primarykey"`
	Brand        string
	SerialNumber string
	Price        float64
	Label        string
	Factory      string
	Remark       string

	Availiable bool
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
	// 设备个体
	// EquipmentUnits []EquipmentUnit
	Availiable bool
}

type Request struct {
	gorm.Model

	BeginAt       time.Time
	SupposeBackAt time.Time
	DoBackAt      time.Time

	EquipmentID     uint
	Equipment       Equipment
	EquipmentUnitID uint
	EquipmentUnit   EquipmentUnit
	UserID          uint
	User            User
	Status          string // REQUESTING
}

type UnAssigned struct {
	Request
}

type Ongoing struct {
	Request
}

type Finished struct {
	Request
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

	db.Create(&Department{Name: "設備課", Description: "設備課，測試"})
	db.Create(&User{Name: "huonwe", Password: "huonwe", DepartmentName: "設備課"})
	db.Create(&Equipment{Name: "测试设备", Type: "试做型", Class: "醫用設備", Availiable: true})
	db.Create(&Equipment{Name: "测试设备", Type: "试做型", Class: "未来科技", Availiable: true})
	db.Create(&EquipmentUnit{EquipmentName: "测试设备", ID: 0001, Brand: "宏偉製造", SerialNumber: "001", Price: 999.9, Label: "沒有標註", Factory: "宏偉天津製造工廠", Availiable: true})

	// db.Preload("User").Find(&see, 11)
	// fmt.Println(see.UserID)
}
