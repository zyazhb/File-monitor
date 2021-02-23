package model

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//GetSession 路由获取session判断是否登录
func GetSession(c *gin.Context) bool {
	session := sessions.Default(c)
	loginuser := session.Get("loginuser")
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
			c.Redirect(http.StatusFound, "/login")
		}
		return false
	}
	return true
}

//GetRandomString 生成随机字符串
func GetRandomString(seed int64) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//GenMD5 生成32位MD5
func GenMD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

//PrintLog 打印日志
func PrintLog(content string) {
	log.Println("\033[1;32m" + content + "\033[0m")
}
