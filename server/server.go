package main

import (
	"main/model"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)

	//初始化数据库
	exist, err := IsExists("./user.db")
	if !exist && err == nil {
		model.DbInit()
	}
	go RPCServer()

	// 初始化Gin
	router := gin.Default()
	router.LoadHTMLGlob("static/templates/*")
	router.StaticFS("/js", http.Dir("static/js"))
	router.StaticFS("/css", http.Dir("static/css"))
	router.StaticFS("/img", http.Dir("static/img"))
	router.StaticFS("/include", http.Dir("static/include"))

	router.GET("/rpc/:key", model.RPCHandler)
	store := cookie.NewStore([]byte("loginuser"))
	router.Use(sessions.Sessions("sessionid", store))
	{
		router.GET("/", model.IndexHandler)
		router.GET("/register", model.Register)
		router.POST("/register", model.RegisterForm)
		router.GET("/login", model.LoginHandler)
		router.POST("/login", model.Checkin)
		router.GET("/logout", model.LogoutHandler)
		router.GET("/manager", model.ManagerHandler)
		router.GET("/getresult/:page", model.GetReport)
		router.GET("/getresult/", model.GetReportCount)
		router.GET("/delete/:rid", model.DeleteReport)
		router.GET("/usermanager", model.UserManager)
		router.GET("/getalluser", model.UserManage)
		router.GET("/showinfo/", model.ShowInfo)
		router.GET("/showinfo/:uid", model.ShowInfo)
		router.POST("/editor", model.Editor)
	}

	router.NoRoute(model.NotFoundHandle)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
