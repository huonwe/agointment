package main

import (
	"time"

	"gorm.io/gorm"
)

const (
	Requesting int = iota
	Timing
	Completed

	Timeout

	Assigned

	Rejected
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

type Equipment struct {
	ID           uint `gorm:"primarykey"`
	EID          uint
	Name         string
	Brand        string
	SerialNumber string
	Type         string
	Price        float64
	Class        string
	Label        string
	Factory      string
	Remark       string

	Availiable bool
	// 佔用這個設備的人
	UserID uint
	User   User
}

type EquipmentFree struct {
	Name string
}

type Request struct {
	gorm.Model

	BeginAt       time.Time
	SupposeBackAt time.Time
	DoBackAt      time.Time

	EquipmentNumber string `gorm:"unique"`
	EquipmentID     uint
	Equipment       Equipment
	UserID          uint
	User            User
	Status          int
}

type UnAssigned struct {
	gorm.Model

	RequestID uint
	Request   Request
}

type Ongoing struct {
	gorm.Model
	RequestID uint
	Request   Request
}

func initDB(db *gorm.DB) {
	db.Exec("DROP TABLE departments")
	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE equipment") // equipment is uncountable
	db.Exec("DROP TABLE requests")

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Department{})
	db.AutoMigrate(&Equipment{})
	db.AutoMigrate(&Request{})
	db.AutoMigrate(&UnAssigned{})
	db.AutoMigrate(&Ongoing{})

	db.Create(&Department{Name: "設備課", Description: "設備課，測試"})
	db.Create(&User{Name: "huonwe", Password: "huonwe", DepartmentName: "設備課"})
	db.Create(&Equipment{Name: "測試設備", ID: 0001, EID: 1234, Type: "试做型", Brand: "宏偉製造", SerialNumber: "001", Price: 999.9, Class: "醫用設備", Label: "沒有標註", Factory: "宏偉天津製造工廠", Availiable: true})

	// db.Preload("User").Find(&see, 11)
	// fmt.Println(see.UserID)
}
