package model

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//LoginHandler 登录页
func LoginHandler(c *gin.Context) {
	if CheckLogin(c, false) {
		c.Redirect(http.StatusFound, "/manager")
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

//LogoutHandler 退出登录
func LogoutHandler(c *gin.Context) {
	if CheckLogin(c, false) {
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
		SetSession(c, uid, email, role)
		c.Redirect(http.StatusFound, "/manager")
	} else {
		c.HTML(http.StatusFound, "login.html", gin.H{"err": "Fail to login !", "errshow": "show"})
	}
}

//UserManager 用户管理器页面
func UserManager(c *gin.Context) {
	CheckLogin(c, true)
	c.HTML(http.StatusOK, "usermanager.html", gin.H{"userdata": AllUserInfo()})
}

//UserManage 取得所有用户信息
func UserManage(c *gin.Context) {
	CheckLogin(c, true)
	if CheckAdmin(c) {
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
	//拿到要修改用户基本信息
	var user User
	user = DbGetByuid(&user, uid)
	if pass == repass {
		//用户想要修改角色吗
		newrole, _ := strconv.Atoi(c.PostForm("role"))
		if newrole != user.Role && !CheckAdmin(c) {
			newrole = user.Role
			c.HTML(http.StatusOK, "showinfo.html", gin.H{"User": user, "err": "No permission", "errshow": "show"})
		}
		//用户想要修改密码吗
		if pass != user.Password {
			pass = GenMD5(pass)
		}
		EditUser(uid, newrole, email, pass)
		//修改成功后更新信息
		user = DbGetByuid(&user, uid)
		c.HTML(http.StatusOK, "showinfo.html", gin.H{"User": user, "success": "Saved!", "sucshow": "show"})
	} else {
		c.HTML(http.StatusOK, "showinfo.html", gin.H{"User": user, "err": "password isn't match", "errshow": "show"})
	}
}

//DelUser 删除用户
func DelUser(c *gin.Context) {
	CheckLogin(c, true)
	if CheckAdmin(c) {
		DbDelUser(c.Param("uid"))
	}
}

//AddUser 添加用户
func AddUser(c *gin.Context) {
	CheckLogin(c, true)
	if CheckAdmin(c) {
		role, _ := strconv.Atoi(c.PostForm("role"))
		DbAddUser(c.PostForm("email"), c.PostForm("password"), role)
	}
}

//Register 注册页
func Register(c *gin.Context) {
	if !CheckLogin(c, false) || CheckAdmin(c) {
		c.HTML(http.StatusOK, "register.html", nil)
	}
	c.Redirect(http.StatusFound, "/")

}

//RegisterForm 接收注册数据
func RegisterForm(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	if password == repassword {
		err := DbAddUser(email, GenMD5(password), 999)
		if err != nil {
			c.HTML(http.StatusFound, "register.html", gin.H{"err": "User name already exist!", "errshow": "show"})
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
	} else {
		//留给js写弹窗 两次密码不匹配
		c.HTML(http.StatusFound, "register.html", gin.H{"err": "Repeat password isn't match!", "errshow": "show"})
	}
}
