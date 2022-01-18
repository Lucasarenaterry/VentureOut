package main

import (
	//"fmt"

	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	//"database/sql"

	_ "github.com/lib/pq"
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

	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//if err != nil {
		//log.Fatalf("Error opening database: %q", err)
	//}
	
	r := gin.New()
	r.Use(gin.Logger())
	r.Static("/css", "./static/css")
	r.LoadHTMLGlob("static/templates/*.html")
	r.Static("/static", "static")
	r.Static("/js", "./static/js")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"navtitle": "VentureOut"})
	})

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"navtitle": "VentureOut"})
	})

	r.GET("/map", func(c *gin.Context) {
		//if _, err := db.Exec("CREAT TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtitle varchar(45) NOT NULL)");
			//err != nil {
		//c.String(http.StatusInternalServerError,
				//fmt.Sprintf("Error creating database table: %q", err))
			//return
			//}


		c.HTML(http.StatusOK, "map.html", gin.H{"navtitle": "VentureOut"})
	})

	r.Run(":" + port) // listen and serve on 0.0.0.0:5000
}
