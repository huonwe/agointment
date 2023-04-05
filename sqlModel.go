package main

import (
	"errors"
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
	Name        string `gorm:"unique"`
	Description string
	Availiable  bool `gorm:"default:true"`
	// ...
}

type User struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Password string
	// 微信的openid 待用
	WeChat string

	// DeptName     string
	DepartmentID uint
	Department   Department

	UpdatedAt time.Time

	IsAdmin bool `gorm:"default:false"`
	IsSuper bool `gorm:"default:false"`

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
	EquipmentID uint `gorm:"index"`
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

	LastMaintained time.Time

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
	EndAt        time.Time `gorm:"index"`
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
	Status  string `gorm:"index"`
	Deleted gorm.DeletedAt
}

type Maintain struct {
	ID uint `gorm:"primaryKey"`

	UnitID   uint
	UnitName string
	UID      string
	Type     string
	Serial   string
	// Unit   EquipmentUnit
	DoAt    time.Time `gorm:"index"`
	DoAtStr string
}

type UnitMaintainAPI struct {
	ID                uint
	Name              string
	UID               string
	LastMaintained    time.Time
	LastMaintainedStr string `gorm:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password == "" {
		return errors.New("创建用户时密码为空")
	}
	u.Password = md5_str(u.Password)
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Password == "" {
		return
	}
	u.Password = md5_str(u.Password)
	return nil
}

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

func initDB(db *gorm.DB) {
	// db.Exec("SET FOREIGN_KEY_CHECKS=0;")
	// db.Exec("DROP TABLE departments, users, equipment_units, equipment, requests, maintains")
	// db.Exec("SET FOREIGN_KEY_CHECKS=1;")

	db.AutoMigrate(&Department{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&EquipmentUnit{})
	db.AutoMigrate(&Equipment{})
	db.AutoMigrate(&Request{})
	db.AutoMigrate(&Maintain{})

	var count int64
	var dept Department
	db.First(&dept, 1)
	if dept.ID == 0 {
		dept.Name = "保留部门"
		dept.Availiable = false
		db.Create(&dept)
		// db.Where(&Department{Name: dept.Name}).Select("Name", "Availiable").Updates(&dept)
	}

	db.Where(&User{Name: "admin", IsAdmin: true}).Count(&count)
	if count == 0 {
		db.Create(&User{ID: 1, Name: "admin", Password: "password123456", DepartmentID: 1, IsAdmin: true, IsSuper: true})
	}
}
