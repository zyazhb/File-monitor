package model

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
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
	CheckLogin(c, true)
	c.HTML(http.StatusOK, "index.html", nil)
}

//LoginHandler 登录页
func LoginHandler(c *gin.Context) {
	if CheckLogin(c, false) == true {
		log.Println("你已经登录了")
		c.Redirect(http.StatusMovedPermanently, "/")
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

//LoginOutHandler 退出登录
func LoginOutHandler(c *gin.Context) {
	if CheckLogin(c, false) == true {
		session := sessions.Default(c)
		session.Delete("loginuser")
	}
	c.Redirect(http.StatusMovedPermanently, "/login")
}

//Checkin 接收前端数据
func Checkin(c *gin.Context) {
	//接收数据
	email := c.PostForm("email")
	password := c.PostForm("password")
	var user User
	DbSel(&user)
	if email == user.Email && password == user.Password {
		//邮箱和密码验证成功之后设置session
		session := sessions.Default(c)
		session.Set("loginuser", email)
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
}

//ManagerHandler 控制台
func ManagerHandler(c *gin.Context) {
	CheckLogin(c, true)
	c.HTML(http.StatusOK, "manager.html", nil)
}

//Register 注册页
func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

//NotFoundHandle 404页面
func NotFoundHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "404.html", nil)
}
