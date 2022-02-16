package main

import (
	"fmt"
	"strings"

	"log"
	"net/http"
	"os"

	//"time"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Event struct {
	Id          int `json:"Id"`
	Eventtittel string `json:"Eventtittel"`
	Eventtype   string `json:"Eventtype"`
	Description string `json:"Description"`
	Image       string `json:"Image"`
	Date        string `json:"Date"`
}

type EventFilter struct {
	// Id        int `json:"Id"`
	Eventtype string `json:"Eventtype"`
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
	r.StaticFile("manifest.webmanifest", "./manifest.webmanifest")
	r.StaticFile("service-worker.js", "./service-worker.js")

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
			var Eventtittel string 
			var Eventtype string
			var Description string 
			var Image string 
			var Date string 

			
			events := make([]Event, 0)

			for rows.Next() {
				
				if err := rows.Scan(&Eventtittel, &Eventtype, &Description, &Image, &Date); 
				err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error scanning events: %q", err))
					return
				}
				// event := &Event{
				// 	Eventtittel: Eventtittel,
				// 	Eventtype: Eventtype,
				// 	Description: Description,
				// 	Image: Image,
				// 	Date: Date,
				// }
				
				 events = append(events, Event{
					 	Eventtittel: Eventtittel,
						Eventtype: Eventtype,
						Description: Description,
						Image: Image,
						Date: Date,
					})
			}
				
			
			

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

	r.POST("/map", func(c *gin.Context) {
		if err := c.Request.ParseForm();err != nil {
			c.String(http.StatusInternalServerError,
			fmt.Sprintf("ParseForm() err: %v", err))
			return
		}

		// filter := c.Request.FormValue("filter")
		filterSlice := c.Request.Form["filter"] //gets all the eventtypes in a slice string format
		filters := strings.Join(filterSlice, ",") //convert to string format that db can read
		filter := "{" + filters + "}"
		fmt.Printf("filter %v", filter)

		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, image TEXT, location GEOMETRY(POINT,4326), eventdate DATE, eventtime TIME)"); 
				err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}

			
			rows, err := db.Query("SELECT json_build_object( 'type', 'FeatureCollection', 'features', json_agg( json_build_object( 'type', 'Feature', 'properties', to_jsonb( t.* ) - 'location', 'geometry', ST_AsGeoJSON(location)::jsonb ) ) ) AS json FROM events as t(id, eventtittel, eventtype, description, image, location, eventdate, eventtime) WHERE eventtype = ANY($1)", filter)
				if err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}
			
			var featureCollection string

			defer rows.Close()
			
			for rows.Next() {
				if err := rows.Scan(&featureCollection); 
					err != nil {
						c.String(http.StatusInternalServerError,
							fmt.Sprintf("Error scanning events: %q", err))
						return
					}
			}
			fmt.Printf("%v", featureCollection)
			
			
			rowss, err := db.Query("SELECT DISTINCT eventtype FROM events")
				if err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}
			
			var Eventtype string

			filterTypes := make([]EventFilter, 0)

			defer rowss.Close()
			
			for rowss.Next() {
				if err := rowss.Scan(&Eventtype); 
					err != nil {
						c.String(http.StatusInternalServerError,
							fmt.Sprintf("Error scanning events: %q", err))
						return
					}

					filterTypes = append(filterTypes, EventFilter{
					   Eventtype: Eventtype,
				   })
			}
			fmt.Printf("%v", filterTypes)

			
			c.HTML(http.StatusOK, "map.html", gin.H{ "featureCollection": featureCollection, "filterTypes": filterTypes, })
	})

	//INSERT INTO events (id, eventtittel, eventtype, description, image, location, eventdate)
	//VALUES (1, 'LANDSCAPE TRAIL', 'Walk', 'walk around campus visiting the main landcapes', 'landscape.png', 'SRID=4326;POINT(-3.321578 55.910807)', '2022/01/23');

	r.GET("/map", func(c *gin.Context) {
			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, image TEXT, location GEOMETRY(POINT,4326), eventdate DATE, eventtime TIME)"); 
				err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}

			
			rows, err := db.Query("SELECT json_build_object( 'type', 'FeatureCollection', 'features', json_agg( json_build_object( 'type', 'Feature', 'properties', to_jsonb( t.* ) - 'location', 'geometry', ST_AsGeoJSON(location)::jsonb ) ) ) AS json FROM events as t(id, eventtittel, eventtype, description, image, location, eventdate, eventtime)")
				if err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}
			
			var featureCollection string

			defer rows.Close()
			
			for rows.Next() {
				if err := rows.Scan(&featureCollection); 
					err != nil {
						c.String(http.StatusInternalServerError,
							fmt.Sprintf("Error scanning events: %q", err))
						return
					}
			}
			fmt.Printf("%v", featureCollection)
			
			rowss, err := db.Query("SELECT DISTINCT eventtype FROM events")
				if err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}
			
			var Eventtype string

			filterTypes := make([]EventFilter, 0)

			defer rowss.Close()
			
			for rowss.Next() {
				if err := rowss.Scan(&Eventtype); 
					err != nil {
						c.String(http.StatusInternalServerError,
							fmt.Sprintf("Error scanning events: %q", err))
						return
					}

					filterTypes = append(filterTypes, EventFilter{
					   Eventtype: Eventtype,
				   })
			}
			fmt.Printf("%v", filterTypes)

			
			c.HTML(http.StatusOK, "map.html", gin.H{ "featureCollection": featureCollection, "filterTypes": filterTypes, })
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
