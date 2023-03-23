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

type User struct {
	gorm.Model
	Name     string
	Password string

	Department string
	IsAdmin    bool

	Requests []Request
}

type Equipment struct {
	gorm.Model
	Name         string
	Band         string
	Type         string
	SerialNumber string
	Price        float64
	Department   string
	Contract     string
	Status       string
	Class        string
	Label        string
	Factory      string
	Remark       string

	Availiable bool
}

type Request struct {
	gorm.Model

	BeginAt       time.Time
	SupposeBackAt time.Time
	DoBackAt      time.Time

	EquipmentID uint
	Equipment   Equipment
	UserID      uint
	User        User
	Process     int
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
