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
	islogin := GetSession(c)
	if islogin == true {
		c.HTML(http.StatusOK, "index.html", nil)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
}

//LoginHandler 登录页
func LoginHandler(c *gin.Context) {
	islogin := GetSession(c)
	log.Println("islogin", islogin)
	if islogin == false {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		log.Println("你已经登录了")
	}
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
		//设置cookie
		// cookie, err := c.Cookie("gin_cookie")
		// if err != nil {
		// 	cookie = "NotSet"
		// 	c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		// } else {
		// 	log.Println("cookie value: ", cookie)
		// }
	} else {
		log.Printf("登录失败")
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

//GetSession 路由获取session判断是否登录
func GetSession(c *gin.Context) bool {
	session := sessions.Default(c)
	loginuser := session.Get("loginuser")
	log.Println("loginuser:", loginuser)
	if loginuser != nil {
		return true
	}
	return false
}
