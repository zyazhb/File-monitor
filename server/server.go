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



func main(){
	DbInit()
	go RpcServer()

	router := gin.Default()
	router.LoadHTMLGlob("static/templates/*")
	router.GET("/", IndexHandler)

	router.GET("/rpc/:key", RpcHandler)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}