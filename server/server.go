package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func RpcHandler(c *gin.Context) {
	vars := c.Param("key")
	log.Println(vars)
	c.JSON(200,gin.H{
            "msg":vars,
        })
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": "This is a content",
	})
}

func LoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"static_js": "/js",
		"static_css": "/css",
	})
}



func main(){
	DbInit()
	go RpcServer()

	router := gin.Default()
	router.LoadHTMLGlob("static/templates/*")
	router.StaticFS("/js", http.Dir("static/js"))
	router.StaticFS("/css", http.Dir("static/css"))

	router.GET("/", IndexHandler)

	router.GET("/login", LoginHandler)


	router.GET("/rpc/:key", RpcHandler)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}