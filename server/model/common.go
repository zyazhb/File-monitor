package model

import (
	"log"
	"net/http"
	"strconv"

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
	if CheckLogin(c, false) == true {
		c.Redirect(http.StatusFound, "/manager")
	}
	c.Redirect(http.StatusFound, "/login")
}

//ManagerHandler 控制台
func ManagerHandler(c *gin.Context) {
	CheckLogin(c, true)
	c.HTML(http.StatusOK, "manager.html", nil)
}

//GetReport 取得报告
func GetReport(c *gin.Context) {
	CheckLogin(c, true)
	page, _ := strconv.Atoi(c.Param("page"))
	result := RPCDbSel(page)
	c.JSON(200, result)
}

//DeleteReport 删除报告
func DeleteReport(c *gin.Context) {
	CheckLogin(c, true)
	rid := c.Param("rid")
	RPCDbDel(rid)
	c.JSON(200, gin.H{
		"msg": "Success",
	})
}

//NotFoundHandle 404页面
func NotFoundHandle(c *gin.Context) {
	c.HTML(http.StatusOK, "404.html", nil)
}
