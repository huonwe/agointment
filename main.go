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

	group_authless := r.Group("/")
	group_authless.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})
	group_authless.POST("/login", login)

	r.Use(LoginFilter())
	r.GET("/", index)

	group_equipment := r.Group("/equipment")
	group_equipment.GET("/availiable", getAvailiable)     // HTML
	group_equipment.GET("/makeRequest", equipmentRequest) // JSON
	// group_equipment.GET("/availiableUnits", getAvailiableUnits) // HTML

	group_user := r.Group("/user")
	group_user.GET("/myRequest", myRequest)
	group_user.GET("/myRequestOp", myRequestOp)

	group_admin := r.Group("/admin")
	group_admin.Use(MustAdmin())
	group_admin.GET("/requestings", adminRequestings)
	group_admin.GET("/requestingsOp", adminRequestingsOp)
	group_admin.GET("/assignUnits", getAvailiableUnits) // HTML
	group_admin.POST("/assignUnits", assignUnits)       // HTML
	group_admin.GET("/ongoings", adminOngoings)         // HTML
	group_admin.GET("/all", adminAll)

	r.Run("0.0.0.0:5501")
}
