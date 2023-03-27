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
	EquipmentID uint
	Equipment   Equipment
	Type        string
	Class       string
	Name        string

	ID           uint `gorm:"primarykey"`
	Brand        string
	SerialNumber string
	Price        float64
	Label        string
	Factory      string
	Remark       string
	Status       string

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
	EquipmentUnitID uint
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

type Finished struct {
	Request
}

func initDB(db *gorm.DB) {
	// db.Exec("DROP TABLE departments")
	// db.Exec("DROP TABLE users")
	// db.Exec("DROP TABLE equipment_units")

	// db.Exec("DROP TABLE equipment") // equipment is uncountable

	// db.Exec("DROP TABLE requests")
	// db.Exec("DROP TABLE un_assigneds")
	// db.Exec("DROP TABLE ongoings")

	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&Department{})
	// db.AutoMigrate(&EquipmentUnit{})
	// db.AutoMigrate(&Equipment{})
	// db.AutoMigrate(&Request{})
	// db.AutoMigrate(&UnAssigned{})
	// db.AutoMigrate(&Ongoing{})

	// db.Create(&Department{Name: "德国骨科", Description: "德国骨科..."})
	// db.Create(&User{Name: "test", Password: "123456", DepartmentName: "德国骨科", IsAdmin: false})
	// db.Create(&Equipment{Name: "测试设备", Type: "试做型", Class: "醫用設備", Availiable: true})
	// db.Create(&Equipment{Name: "测试设备", Type: "试做型", Class: "未来科技", Availiable: true})
	// db.Create(&Equipment{Name: "空想具现", Type: "试做型", Class: "宏伟制造", Availiable: true})

	// db.Create(&EquipmentUnit{Name: "测试设备", Type: "试做型", Class: "醫用設備", ID: 0001, Brand: "宏偉製造", SerialNumber: "001", Price: 999.9, Label: "沒有標註", Factory: "宏偉天津製造工廠", Availiable: true})
	// db.Create(&EquipmentUnit{Name: "测试设备", Type: "试做型", Class: "醫用設備", ID: 0002, Brand: "宏偉製造", SerialNumber: "001", Price: 999.9, Label: "沒有標註", Factory: "宏偉天津製造工廠", Availiable: true})

}
