package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	handle(err)
	initDB(db)

	r := gin.Default()

	r.LoadHTMLGlob("html/*")
	r.Static("static", "static")

	r.GET("/user/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/user/login", login)
	r.Use(LoginFilter())
	r.GET("/", index)
	r.GET("/equipment/getAvailiable", getAvailiable)

	r.Run("0.0.0.0:5500")
}
