package main

import (
	//"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type event struct {
	id int
	name string
	eventtype string 
	description string 
	img string 
}

func main() {
	
	r := gin.New()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("static/templates/*.html")
	
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"navtitle": "VentureOut."})
	})

	r.Run(":8080")
}
