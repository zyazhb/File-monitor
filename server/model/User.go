package model

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//LoginHandler 登录页
func LoginHandler(c *gin.Context) {
	if CheckLogin(c, false) == true {
		c.Redirect(http.StatusFound, "/manager")
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

//LogoutHandler 退出登录
func LogoutHandler(c *gin.Context) {
	if CheckLogin(c, false) == true {
		session := sessions.Default(c)
		session.Delete("loginuser")
		session.Save()
	}
	c.Redirect(http.StatusFound, "/login")
}

//Checkin 检查登录
func Checkin(c *gin.Context) {
	//接收数据
	email := c.PostForm("email")
	password := c.PostForm("password")
	var user User
	uid, role := DbSel(&user, email, password)
	if uid > 0 {
		//邮箱和密码验证成功之后设置session
		session := sessions.Default(c)
		session.Set("loginuser", email)
		session.Set("role", role)
		session.Save()
		c.Redirect(http.StatusFound, "/manager")
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
}

//UserManager 用户管理器页面
func UserManager(c *gin.Context) {
	CheckLogin(c, true)
	c.HTML(http.StatusOK, "usermanager.html", nil)
}

//UserManage 取得所有用户信息
func UserManage(c *gin.Context) {
	CheckLogin(c, true)
	session := sessions.Default(c)
	role := session.Get("role")
	if role == 0 {
		result := AllUserInfo()
		c.JSON(200, result)
	} else {
		c.JSON(200, gin.H{"No Access Permission": "No Access Permission"})
	}
}

//showinfo 展示可修改信息
func ShowInfo(c *gin.Context) {
	CheckLogin(c, true)
	c.HTML(http.StatusOK, "showinfo.html", nil)
}

//Register 注册页
func Register(c *gin.Context) {
	if CheckLogin(c, false) == true {
		c.Redirect(http.StatusFound, "/")
	}
	c.HTML(http.StatusOK, "register.html", nil)
}

//RegisterForm 接收注册数据
func RegisterForm(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	if password == repassword {
		err := DbInsert(email, password)
		if err != nil {
			c.Redirect(http.StatusFound, "/register")
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
	} else {
		//留给js写弹窗 两次密码不匹配
		c.Redirect(http.StatusFound, "/register")
	}
}
