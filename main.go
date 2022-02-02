package main

import (
	"fmt"

	"log"
	"net/http"
	"os"
	//"time"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	Id          int `json:"Id"`
	Eventtittel string `json:"Eventtittel"`
	Eventtype   string `json:"Eventtype"`
	Description string `json:"Description"`
	Image         string `json:"Image"`
	Date         string `json:"Date"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	} 

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	
	r := gin.New()
	r.Use(gin.Logger())
	r.Static("/css", "./static/css")
	r.LoadHTMLGlob("static/templates/*.html")
	r.Static("/static", "static")
	r.Static("/js", "./static/js")

	r.GET("/", func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, image TEXT, location GEOMETRY(POINT,4326), eventdate DATE, eventtime TIME)"); 
					err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}

			rows, err := db.Query("SELECT eventtittel, eventtype, description, image, eventdate FROM Events")
			if err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}
   
			defer rows.Close()
			// var Eventtittel string 
			// var Eventtype string
			// var Description string 
			// var Image string 
			// var Date string 

			
			events := make([]Event, 0)

			for rows.Next() {
				event := Event{}
				if err := rows.Scan(&event.Eventtittel, &event.Eventtype, &event.Description, &event.Image, &event.Date); 
				err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error scanning events: %q", err))
					return
				}
				
				 events = append(events, event)

			}
				// 	fmt.Println(events)
			
			

		c.HTML(http.StatusOK, "index.html", gin.H{
			// "eventtittel": Eventtittel,
			// "eventtype": Eventtype,
			// "description": Description,
			// "image": Image,
			// "date": Date,
			"events": events,
		})
	})

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	//INSERT INTO events (id, eventtittel, eventtype, description, image, location, eventdate)
	//VALUES (1, 'LANDSCAPE TRAIL', 'Walk', 'walk around campus visiting the main landcapes', 'landscape.png', 'SRID=4326;POINT(-3.321578 55.910807)', '2022/01/23');

	r.GET("/map", func(c *gin.Context) {
			// if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Events ticks (tick timestamp)"); 
			// 		err != nil {
			// 	c.String(http.StatusInternalServerError,
			// 		fmt.Sprintf("Error creating database table: %q", err))
			// 	return
			// }
   
			// if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); 
			// 	err != nil {
			// 	c.String(http.StatusInternalServerError,
			// 		fmt.Sprintf("Error incrementing tick: %q", err))
			// 	return
			// }
   
			// rows, err := db.Query("SELECT tick FROM ticks")
			// if err != nil {
			// 	c.String(http.StatusInternalServerError,
			// 		fmt.Sprintf("Error reading ticks: %q", err))
			// 	return
			// }
   
			// defer rows.Close()
			// for rows.Next() {
			// 	var tick time.Time
			// 	if err := rows.Scan(&tick); err != nil {
			// 		c.String(http.StatusInternalServerError,
			// 			fmt.Sprintf("Error scanning ticks: %q", err))
			// 		return
			// 	}
			// 	c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
			// }
			c.HTML(http.StatusOK, "map.html", gin.H{})
		})

		r.GET("/addevent", func(c *gin.Context) {
			c.HTML(http.StatusOK, "addevent.html", gin.H{})
		})

		r.GET("/settings", func(c *gin.Context) {
			c.HTML(http.StatusOK, "settings.html", gin.H{})
		})

		r.GET("/scan", func(c *gin.Context) {
			c.HTML(http.StatusOK, "scan.html", gin.H{})
		})
	

	r.Run(":" + port) // listen and serve on 0.0.0.0:5000
}
