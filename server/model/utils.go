package model

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//GetSession 路由获取session判断是否登录
func GetSession(c *gin.Context) bool {
	session := sessions.Default(c)
	loginuser := session.Get("loginuser")
	// log.Println("loginuser:", loginuser)
	if loginuser != nil {
		return true
	}
	return false
}

//CheckLogin 检查登录态
func CheckLogin(c *gin.Context, RedirectFlag bool) bool {
	islogin := GetSession(c)
	if islogin == false {
		if RedirectFlag == true {
			c.Redirect(http.StatusMovedPermanently, "/login")
		}
		return false
	}
	return true
}
