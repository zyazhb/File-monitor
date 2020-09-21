package model

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
