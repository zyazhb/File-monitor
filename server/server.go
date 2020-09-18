package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//RPCHandler 暂时保留也许没用
func RPCHandler(c *gin.Context) {
	vars := c.Param("key")
	log.Println(vars)
	c.JSON(200, gin.H{
		"msg": vars,
	})
}

//IndexHandler 首页
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": "This is a content",
	})
}

//LoginHandler 登录页
func LoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

//ManagerHandler 控制台
func ManagerHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "manager.html", nil)
}

//NotFoundHandle 404页面
func NotFoundHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "404.html", nil)
}

func main() {
	//初始化数据库
	DbInit()
	//并发Rpc服务器
	go RpcServer()

	// 初始化Gin
	router := gin.Default()
	router.LoadHTMLGlob("static/templates/*")
	router.StaticFS("/js", http.Dir("static/js"))
	router.StaticFS("/css", http.Dir("static/css"))
	router.StaticFS("/img", http.Dir("static/img"))

	router.GET("/", IndexHandler)
	router.GET("/login", LoginHandler)
	router.GET("/manager", ManagerHandler)
	router.GET("/rpc/:key", RPCHandler)

	router.NoRoute(NotFoundHandle)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
