package main

import (
	"net/http"
	"time"

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

	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.LoadHTMLGlob("html/*")
	r.Static("static", "static")

	group_authless := r.Group("/")
	group_authless.GET("/shell", func(ctx *gin.Context) {
		time.Sleep(12 * time.Second)             // 优雅的准备12s
		ctx.String(http.StatusBadRequest, "www") // www
		ctx.Abort()                              // 人生多别离
	})
	group_authless.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})
	group_authless.POST("/login", login)
	group_authless.GET("/signup", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "signup.html", nil)
	})
	group_authless.POST("/signup", signup)

	group_home := r.Group("/home")
	group_home.Use(LoginFilter())
	group_home.GET("/:name", index)
	group_home.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/home/index")
	})

	group_equipment := r.Group("/equipment")
	group_equipment.Use(LoginFilter())
	group_equipment.GET("/availiable", getAvailiable)     // HTML
	group_equipment.GET("/makeRequest", equipmentRequest) // JSON

	group_equipment.Use(MustAdmin())
	group_equipment.POST("/import", adminEquipmentImport)
	group_equipment.GET("/equipmentOp", equipmentOp) // JSON
	group_equipment.GET("/export", adminExportEquipment)
	group_equipment.GET("/getTemplate", equipmentTemplate)

	group_user := r.Group("/user")
	group_user.Use(LoginFilter())
	group_user.GET("/myRequest", myRequest)
	group_user.GET("/myRequestOp", myRequestOp)

	group_admin := r.Group("/admin")
	group_admin.Use(LoginFilter(), MustAdmin())
	group_admin.GET("/requestings", adminRequestings)
	group_admin.GET("/requestingsOp", adminRequestingsOp)
	group_admin.GET("/assignUnits", getAvailiableUnits) // HTML
	group_admin.POST("/assignUnits", assignUnits)       // HTML
	group_admin.GET("/ongoings", adminOngoings)         // HTML
	group_admin.GET("/all", adminAll)
	group_admin.GET("/users", adminUsers)
	group_admin.POST("/users/:op", adminUsersOp) // JSON or HTML

	group_admin.GET("/equipment", adminEquipment)

	group_admin.GET("/exportRequests", adminExportRequests)

	r.NoRoute(redirect2home)

	r.Run("0.0.0.0:5501")
}
