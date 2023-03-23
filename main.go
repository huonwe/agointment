package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	handle(err)
	db.AutoMigrate(&User{})
	db.Create(&User{Username: "huonwe", Password: "huonwe"})

	r := gin.Default()

	r.LoadHTMLGlob("html/*")
	r.Static("static", "static")

	r.GET("/user/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/user/login", login)
	r.Use(LoginFilter())
	r.GET("/", index)

	r.Run("0.0.0.0:5500")
}
