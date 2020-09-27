package model

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

//Checkin 接收前端数据
func Checkin(c *gin.Context) {
	//接收数据
	email := c.PostForm("email")
	password := c.PostForm("password")
	//连接数据库
	db, err := gorm.Open(sqlite.Open("./foo.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("", err)
	}
	var user User
	db.Select("email,password").Find(&user)
	if email == user.Email && password == user.Password {
		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		fmt.Printf("登录失败")
	}
}

//ManagerHandler 控制台
func ManagerHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "manager.html", nil)
}

//NotFoundHandle 404页面
func NotFoundHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "404.html", nil)
}
