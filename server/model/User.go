package model

import (
	"net/http"
	"strconv"

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
	uid, role := DbSel(&user, email, GenMD5(password))
	if uid > 0 {
		//邮箱和密码验证成功之后设置session
		session := sessions.Default(c)
		session.Set("loginuser", email)
		session.Set("uid", uid)
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

//ShowInfo 展示可修改信息
func ShowInfo(c *gin.Context) {
	CheckLogin(c, true)
	param := c.Param("uid")
	var uid int
	if param != "" {
		uid, _ = strconv.Atoi(param)
	} else {
		session := sessions.Default(c)
		uid = session.Get("uid").(int)
	}
	var user User
	user = DbGetByuid(&user, uid)
	c.HTML(http.StatusOK, "showinfo.html", gin.H{"User": user})
}

//Editor 修改用户信息
func Editor(c *gin.Context) {
	CheckLogin(c, true)
	uid, _ := strconv.Atoi(c.PostForm("uid"))
	email := c.PostForm("email")
	pass := c.PostForm("pass")
	repass := c.PostForm("repass")
	role, _ := strconv.Atoi(c.PostForm("role"))
	var user User
	md5pass := pass
	user = DbGetByuid(&user, uid)
	if uid != user.UID {
		c.HTML(http.StatusOK, "showinfo.html", gin.H{"User": user, "err": "can not edit uid"})
	}
	if pass != repass {
		c.HTML(http.StatusOK, "showinfo.html", gin.H{"User": user, "err": "password don't match"})
	}
	if pass != user.Password {
		md5pass = GenMD5(pass)
	}
	UserEditor(uid, role, email, md5pass)
	c.Redirect(http.StatusFound, "/usermanager")
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
		err := DbInsert(email, GenMD5(password))
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
