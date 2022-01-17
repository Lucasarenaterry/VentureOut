package main

import (
	//"fmt"
	"os"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

type event struct {
	id          int `json:"id"`
	name        string `json:"name"`
	eventtype   string `json:"eventtype"`
	description string `json:"description"`
	img         string `json:"img"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	} 
	
	r := gin.New()
	r.Use(gin.Logger())
	r.Static("/css", "./static/css")
	r.LoadHTMLGlob("static/templates/*.html")
	r.Static("/static", "static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"navtitle": "VentureOut"})
	})

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"navtitle": "VentureOut"})
	})

	r.Run(":" + port) // listen and serve on 0.0.0.0:5000
}
