package main

import (
	"main/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)

	//初始化数据库
	model.DbInit()
	go RPCServer()

	// 初始化Gin
	router := gin.Default()
	router.LoadHTMLGlob("static/templates/*")
	router.StaticFS("/js", http.Dir("static/js"))
	router.StaticFS("/css", http.Dir("static/css"))
	router.StaticFS("/img", http.Dir("static/img"))

	router.GET("/", model.IndexHandler)

	UserRouter := router.Group("admin")
	{
		UserRouter.POST("/register", model.Register)
		UserRouter.POST("/login", model.Login)
	}
	router.GET("/login", model.LoginHandler)
	router.GET("/manager", model.ManagerHandler)
	router.GET("/rpc/:key", model.RPCHandler)
	router.POST("/login", model.Checkin)

	router.NoRoute(model.NotFoundHandle)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
