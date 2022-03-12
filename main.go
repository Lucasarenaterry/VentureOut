package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"strings"
	//"time"

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
	Id          string `json:"Id"`
	Eventtittel string `json:"Eventtittel"`
	Eventtype   string `json:"Eventtype"`
	Description string `json:"Description"`
	OrganizedBy string `json:"OrganizedBy"`
	Image       string `json:"Image"`
	Date        string `json:"Date"`
	EventStartdDate string `json:"EventStartdDate"`
	EventEndDate string `json:"EventEndDate"`
	EventStartTime string `json:"EventStartTime"`
	EventEndTime string `json:"EventEndTime"` 
	ContactEmail string `json:"ContactEmail"`
	EventLink string `json:"EventLink"`
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
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT, latitude NUMERIC(10,8), longitude NUMERIC(11,8))"); 
				err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}

			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS archive (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT, latitude NUMERIC(10,8), longitude NUMERIC(11,8))"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

			if _, err := db.Exec("WITH moved_rows AS ( DELETE FROM events WHERE displaytill < now() RETURNING * ) INSERT INTO archive SELECT * FROM moved_rows"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

			rows, err := db.Query("SELECT id, eventtittel, eventtype, description, organizedby, image, TO_CHAR(eventstartdate, 'DD Mon YYYY'), TO_CHAR(eventenddate , 'DD Mon YYYY'), TO_CHAR(eventstarttime, 'HH24:MI'), TO_CHAR(eventendtime, 'HH24:MI'), contactemail, eventlink FROM events WHERE displaytill > now() LIMIT 25")
			if err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}
   
			defer rows.Close()
			var Id string
			var Eventtittel string 
			var Eventtype string
			var Description string 
			var Image string 
			var OrganizedBy string 
			var EventStartdDate string
			var EventEndDate string
			var EventStartTime string
			var EventEndTime string
			var ContactEmail string
			var EventLink string

			
			events := make([]Event, 0)

			for rows.Next() {
				
				if err := rows.Scan(&Id, &Eventtittel, &Eventtype, &Description, &OrganizedBy, &Image, &EventStartdDate, &EventEndDate, &EventStartTime, &EventEndTime, &ContactEmail, &EventLink); 
				err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error scanning events: %q", err))
					return
				}
				
				 events = append(events, Event{
					 	Id: Id,
					 	Eventtittel: Eventtittel,
						Eventtype: Eventtype,
						Description: Description,
						OrganizedBy: OrganizedBy,
						Image: Image,
						EventStartdDate: EventStartdDate,
						EventEndDate: EventEndDate,
						EventStartTime: EventStartTime,
						EventEndTime: EventEndTime,
						ContactEmail: ContactEmail,
						EventLink: EventLink,
					})
			}
	
		c.HTML(http.StatusOK, "index.html", gin.H{
			"events": events,
		})
	})


	r.POST("/home", func(c *gin.Context) {
		filterSlice := c.Request.FormValue("searchEvent") //gets all the eventtypes in a slice string format
		split := strings.Fields(filterSlice)
		converttodb := strings.Join(split, "%,") //convert to string format that db can read
		searchevents := "{" + converttodb + "%}"
		fmt.Printf("filter %v", searchevents)
		
		
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT)"); 
				err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}

			rows, err := db.Query("SELECT id, eventtittel, eventtype, description, organizedby, image, TO_CHAR(eventstartdate, 'DD Mon YYYY'), TO_CHAR(eventenddate , 'DD Mon YYYY'), TO_CHAR(eventstarttime, 'HH24:MI'), TO_CHAR(eventendtime, 'HH24:MI'), contactemail, eventlink FROM events WHERE eventtittel ILIKE ANY ($1)", searchevents)
			if err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}
   
			defer rows.Close()
			var Id string
			var Eventtittel string 
			var Eventtype string
			var Description string 
			var Image string 
			var OrganizedBy string 
			var EventStartdDate string
			var EventEndDate string
			var EventStartTime string
			var EventEndTime string
			var ContactEmail string
			var EventLink string

			
			events := make([]Event, 0)

			for rows.Next() {
				
				if err := rows.Scan(&Id, &Eventtittel, &Eventtype, &Description, &OrganizedBy, &Image, &EventStartdDate, &EventEndDate, &EventStartTime, &EventEndTime, &ContactEmail, &EventLink); 
				err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error scanning events: %q", err))
					return
				}
				
				 events = append(events, Event{
					 	Id: Id,
					 	Eventtittel: Eventtittel,
						Eventtype: Eventtype,
						Description: Description,
						OrganizedBy: OrganizedBy,
						Image: Image,
						EventStartdDate: EventStartdDate,
						EventEndDate: EventEndDate,
						EventStartTime: EventStartTime,
						EventEndTime: EventEndTime,
						ContactEmail: ContactEmail,
						EventLink: EventLink,
					})
					fmt.Printf("filter %v", Eventtittel)
			}
			
	
		c.HTML(http.StatusOK, "index.html", gin.H{
			"events": events,
		})

	})

	r.GET("/home", func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT, latitude NUMERIC(10,8), longitude NUMERIC(11,8))"); 
				err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}

			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS archive (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT, latitude NUMERIC(10,8), longitude NUMERIC(11,8))"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

			if _, err := db.Exec("WITH moved_rows AS ( DELETE FROM events WHERE displaytill < now() RETURNING * ) INSERT INTO archive SELECT * FROM moved_rows"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

			rows, err := db.Query("SELECT id, eventtittel, eventtype, description, organizedby, image, TO_CHAR(eventstartdate, 'DD Mon YYYY'), TO_CHAR(eventenddate , 'DD Mon YYYY'), TO_CHAR(eventstarttime, 'HH24:MI'), TO_CHAR(eventendtime, 'HH24:MI'), contactemail, eventlink FROM events WHERE displaytill > now() LIMIT 25")
			if err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}
   
			defer rows.Close()
			var Id string
			var Eventtittel string 
			var Eventtype string
			var Description string 
			var Image string 
			var OrganizedBy string 
			var EventStartdDate string
			var EventEndDate string
			var EventStartTime string
			var EventEndTime string
			var ContactEmail string
			var EventLink string

			
			events := make([]Event, 0)

			for rows.Next() {
				
				if err := rows.Scan(&Id, &Eventtittel, &Eventtype, &Description, &OrganizedBy, &Image, &EventStartdDate, &EventEndDate, &EventStartTime, &EventEndTime, &ContactEmail, &EventLink); 
				err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error scanning events: %q", err))
					return
				}
				
				 events = append(events, Event{
					 	Id: Id,
					 	Eventtittel: Eventtittel,
						Eventtype: Eventtype,
						Description: Description,
						OrganizedBy: OrganizedBy,
						Image: Image,
						EventStartdDate: EventStartdDate,
						EventEndDate: EventEndDate,
						EventStartTime: EventStartTime,
						EventEndTime: EventEndTime,
						ContactEmail: ContactEmail,
						EventLink: EventLink,
					})
			}
	
		c.HTML(http.StatusOK, "index.html", gin.H{
			"events": events,
		})
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

		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT)"); 
				err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}

			var featureCollection string

			if filterSlice != nil {
				rows, err := db.Query("SELECT json_build_object( 'type', 'FeatureCollection', 'features', json_agg( json_build_object( 'type', 'Feature', 'properties', to_jsonb( t.* ) - 'location' - 'geofence', 'geometry', ST_AsGeoJSON(location)::jsonb ) ) ) AS json FROM events as t(id, eventtittel, eventtype, description, organizedby, image, location, geofence, eventstartdate, eventenddate, eventstarttime, eventendtime, contactemail, eventlink) WHERE eventtype = ANY($1)", filter)
					if err != nil {
						c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error reading Events: %q", err))
					return
				}

				defer rows.Close()
			
				for rows.Next() {
					if err := rows.Scan(&featureCollection); 
						err != nil {
							c.String(http.StatusInternalServerError,
								fmt.Sprintf("Error scanning events: %q", err))
							return
						}
				}
				//fmt.Printf("%v", featureCollection)
			} else { //this if else case is if no filter is chosen then all event types will be shown
				rows, err := db.Query("SELECT json_build_object( 'type', 'FeatureCollection', 'features', json_agg( json_build_object( 'type', 'Feature', 'properties', to_jsonb( t.* ) - 'location' - 'geofence', 'geometry', ST_AsGeoJSON(location)::jsonb ) ) ) AS json FROM events as t(id, eventtittel, eventtype, description, organizedby, image, location, geofence, eventstartdate, eventenddate, eventstarttime, eventendtime, contactemail, eventlink)")
				if err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
				}
				defer rows.Close()
			
				for rows.Next() {
					if err := rows.Scan(&featureCollection); 
						err != nil {
							c.String(http.StatusInternalServerError,
								fmt.Sprintf("Error scanning events: %q", err))
							return
						}
				}
				//fmt.Printf("%v", featureCollection)
			}
			
			
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
			//fmt.Printf("%v", filterTypes)

			
			c.HTML(http.StatusOK, "map.html", gin.H{ "featureCollection": featureCollection, "filterTypes": filterTypes, })
	})

	//INSERT INTO events (eventtittel, eventtype, description, image, location, eventdate)
	//VALUES ('LANDSCAPE TRAIL', 'Walk', 'walk around campus visiting the main landcapes', 'landscape.png', 'SRID=4326;POINT(-3.321578 55.910807)', '2022/01/23');


	// INSERT INTO events (eventtittel, eventtype, description, organizedby, image, location, geofence, displayfrom, displaytill, eventstartdate, eventenddate, eventstarttime, eventendtime, contactemail, eventlink)
	// VALUES ('LANDSCAPE TRAIL', 'Walk', 'walk around campus visiting the main landcapes', 'Heriot-Watt University', 'landscape.png', 'SRID=4326;POINT(-3.321578 55.910807)', ST_Buffer(geography(ST_POINT(-3.2138, 55.9406)), 3), '2022/02/26', '2022/03/26', '2022/03/01', '2022/03/4', '08:00', '13:00', 'hw@hw.ac.uk', 'https://www.hw.ac.uk/uk/campus-trails.htm');

	r.GET("/map", func(c *gin.Context) {
			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS archive (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT, latitude NUMERIC(10,8), longitude NUMERIC(11,8))"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

			if _, err := db.Exec("WITH moved_rows AS ( DELETE FROM events WHERE displaytill < now() RETURNING * ) INSERT INTO archive SELECT * FROM moved_rows"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT)"); 
				err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}
			OnMapId := c.Query("OnMapId")
			eventid := c.Query("id")
			fmt.Println("Event id is", eventid)
			var qrscanned bool
			qrscanned = false
			var featureCollection string

			var Eventtittel string 
			var Eventtype string
			var Description string 
			var Image string 
			var OrganizedBy string 
			var EventStartdDate string
			var EventEndDate string
			var EventStartTime string
			var EventEndTime string
			var ContactEmail string
			var EventLink string

			if eventid != "" {
				
				rows, err := db.Query("SELECT json_build_object( 'type', 'FeatureCollection', 'features', json_agg( json_build_object( 'type', 'Feature', 'properties', to_jsonb( t.* ) - 'location' - 'geofence', 'geometry', ST_AsGeoJSON(location)::jsonb ) ) ) AS json FROM events as t(id, eventtittel, eventtype, description, organizedby, image, location, geofence, eventstartdate, eventenddate, eventstarttime, eventendtime, contactemail, eventlink) WHERE id = $1", eventid)
					if err != nil {
						c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error reading Events: %q", err))
					return
				}
				
				defer rows.Close()
				
				for rows.Next() {
					if err := rows.Scan(&featureCollection); 
						err != nil {
							c.String(http.StatusInternalServerError,
								fmt.Sprintf("Error scanning events: %q", err))
							return
						}
				}
				//fmt.Printf("%v", featureCollection)

				rowss, err := db.Query("SELECT eventtittel, eventtype, description, organizedby, image, TO_CHAR(eventstartdate, 'DD Mon YYYY'), TO_CHAR(eventenddate , 'DD Mon YYYY'), TO_CHAR(eventstarttime, 'HH24:MI'), TO_CHAR(eventendtime, 'HH24:MI'), contactemail, eventlink FROM events WHERE id = $1", eventid)
				if err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error reading Events: %q", err))
					return
				}
	
				defer rowss.Close()
				

				for rowss.Next() {
				
					if err := rowss.Scan(&Eventtittel, &Eventtype, &Description, &OrganizedBy, &Image, &EventStartdDate, &EventEndDate, &EventStartTime, &EventEndTime, &ContactEmail, &EventLink); 
					err != nil {
						c.String(http.StatusInternalServerError,
							fmt.Sprintf("Error scanning events: %q", err))
						return
					}
				}

				qrscanned = true
			} else if OnMapId != "" {
				rows, err := db.Query("SELECT json_build_object( 'type', 'FeatureCollection', 'features', json_agg( json_build_object( 'type', 'Feature', 'properties', to_jsonb( t.* ) - 'location' - 'geofence', 'geometry', ST_AsGeoJSON(location)::jsonb ) ) ) AS json FROM events as t(id, eventtittel, eventtype, description, organizedby, image, location, geofence, eventstartdate, eventenddate, eventstarttime, eventendtime, contactemail, eventlink) WHERE id = $1", OnMapId)
					if err != nil {
						c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error reading Events: %q", err))
					return
				}
				
				defer rows.Close()
				
				for rows.Next() {
					if err := rows.Scan(&featureCollection); 
						err != nil {
							c.String(http.StatusInternalServerError,
								fmt.Sprintf("Error scanning events: %q", err))
							return
						}
				}
				//fmt.Printf("%v", featureCollection)

			
			} else {
				rows, err := db.Query("SELECT json_build_object( 'type', 'FeatureCollection', 'features', json_agg( json_build_object( 'type', 'Feature', 'properties', to_jsonb( t.* ) - 'location' - 'geofence', 'geometry', ST_AsGeoJSON(location)::jsonb ) ) ) AS json FROM events as t(id, eventtittel, eventtype, description, organizedby, image, location, geofence, eventstartdate, eventenddate, eventstarttime, eventendtime, contactemail, eventlink) WHERE location && ST_MakeEnvelope(-3.3188149475097695, 55.88051417897946, -3.2062050842285195, 55.947473076492535, 4326)")
					if err != nil {
						c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error reading Events: %q", err))
					return
				}
			
				defer rows.Close()
				
				for rows.Next() {
					if err := rows.Scan(&featureCollection); 
						err != nil {
							c.String(http.StatusInternalServerError,
								fmt.Sprintf("Error scanning events: %q", err))
							return
						}
				}
				//fmt.Printf("%v", featureCollection)
			}
			
			rowss, err := db.Query("SELECT DISTINCT eventtype FROM events")
				if err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error reading Events: %q", err))
				return
			}

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
			//fmt.Printf("%v", filterTypes)

			c.HTML(http.StatusOK, "map.html", gin.H{ 
				"featureCollection": featureCollection, 
				"filterTypes": filterTypes, 
				"qrscanned": qrscanned,
			
				"Eventtittel": Eventtittel,
				"Eventtype": Eventtype,
				"Description": Description,
				"OrganizedBy": OrganizedBy,
				"Image": Image,
				"EventStartdDate": EventStartdDate,
				"EventEndDate": EventEndDate,
				"EventStartTime": EventStartTime,
				"EventEndTime": EventEndTime,
				"ContactEmail": ContactEmail,
				"EventLink": EventLink,
			})
		})

		r.POST("/ingeofence/:lat/:lng/:geofencediscovered", func(c *gin.Context) {
			lat := c.Param("lat")
			lng := c.Param("lng")
			GeofenceDiscovered := c.Param("geofencediscovered")
			//fmt.Printf("%v", GeofenceDiscovered)
			
			rows, err := db.Query("SELECT id, eventtittel, eventtype, description, organizedby, image, eventstartdate, eventenddate, TO_CHAR(eventstarttime, 'HH24:MI'), TO_CHAR(eventendtime, 'HH24:MI'), contactemail, eventlink FROM events WHERE ST_Dwithin ( geography (ST_Point(longitude,latitude)), geography (ST_Point($1, $2)), 60) limit 1", lng, lat)
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
			var OrganizedBy string 
			var EventStartdDate string
			var EventEndDate string
			var EventStartTime string
			var EventEndTime string
			var ContactEmail string
			var EventLink string
			var Id string
			var FoundGeofence bool = false

			for rows.Next() {
				
				if err := rows.Scan(&Id, &Eventtittel, &Eventtype, &Description, &OrganizedBy, &Image, &EventStartdDate, &EventEndDate, &EventStartTime, &EventEndTime, &ContactEmail, &EventLink); 
				err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error scanning events: %q", err))
					return
				}
				FoundGeofence = true
			}

			eventJson :=  Event{
			   Id: Id,	
			   Eventtittel: Eventtittel,
			   Eventtype: Eventtype,
			   Description: Description,
			   OrganizedBy: OrganizedBy,
			   Image: Image,
			   EventStartdDate: EventStartdDate,
			   EventEndDate: EventEndDate,
			   EventStartTime: EventStartTime,
			   EventEndTime: EventEndTime,
			   ContactEmail: ContactEmail,
			   EventLink: EventLink,
		   }

		   eventdata, err := json.Marshal(eventJson)
		   
		   if FoundGeofence && Id != GeofenceDiscovered {
		   	c.JSON(200, string(eventdata))
		   }
		})

		r.POST("/searcharea/:SWlng/:SWlat/:NElng/:NElat", func(c *gin.Context) {
			SWlng := c.Param("SWlng")
			SWlat := c.Param("SWlat")
			NElng := c.Param("NElng")
			NElat := c.Param("NElat")
			fmt.Printf("%v", SWlng)
			
			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT)"); 
				err != nil {
					c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error creating database table: %q", err))
				return
			}
			
			var featureCollection string

			rows, err := db.Query("SELECT json_build_object( 'type', 'FeatureCollection', 'features', json_agg( json_build_object( 'type', 'Feature', 'properties', to_jsonb( t.* ) - 'location' - 'geofence', 'geometry', ST_AsGeoJSON(location)::jsonb ) ) ) AS json FROM events as t(id, eventtittel, eventtype, description, organizedby, image, location, geofence, eventstartdate, eventenddate, eventstarttime, eventendtime, contactemail, eventlink) WHERE location && ST_MakeEnvelope($1, $2, $3, $4, 4326)", SWlng, SWlat, NElng, NElat)
					if err != nil {
						c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error reading Events: %q", err))
					return
				}
				
				

				defer rows.Close()
				
				for rows.Next() {
					if err := rows.Scan(&featureCollection); 
						err != nil {
							c.String(http.StatusInternalServerError,
								fmt.Sprintf("Error scanning events: %q", err))
							return
						}
				}

				

				if (len(string(featureCollection)) > 49 ){
				//fmt.Printf("%v", featureCollection)
				c.JSON(200, featureCollection)
				} else {
					c.JSON(200, "No events found")
				}
		    
		})

		r.POST("/calender", func(c *gin.Context) {
				datefrom := c.Request.FormValue("datefrom")
				
				dateto := c.Request.FormValue("dateto")

				if _, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT)"); 
					err != nil {
						c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error creating database table: %q", err))
					return
				}
				
				rows, err := db.Query("SELECT id, eventtittel, eventtype, description, organizedby, image, TO_CHAR(eventstartdate, 'DD Mon YYYY'), TO_CHAR(eventenddate , 'DD Mon YYYY'), TO_CHAR(eventstarttime, 'HH24:MI'), TO_CHAR(eventendtime, 'HH24:MI'), contactemail, eventlink FROM events WHERE eventstartdate BETWEEN $1 AND $2 ORDER BY eventstartdate ASC, eventstarttime ASC", datefrom, dateto)
				if err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error reading Events: %q", err))
					return
				}
	
				defer rows.Close()
				var Id string
				var Eventtittel string 
				var Eventtype string
				var Description string 
				var Image string 
				var OrganizedBy string 
				var EventStartdDate string
				var EventEndDate string
				var EventStartTime string
				var EventEndTime string
				var ContactEmail string
				var EventLink string

				
				events := make([]Event, 0)

				for rows.Next() {
					
					if err := rows.Scan(&Id, &Eventtittel, &Eventtype, &Description, &OrganizedBy, &Image, &EventStartdDate, &EventEndDate, &EventStartTime, &EventEndTime, &ContactEmail, &EventLink); 
					err != nil {
						c.String(http.StatusInternalServerError,
							fmt.Sprintf("Error scanning events: %q", err))
						return
					}
					
					events = append(events, Event{
							Id: Id,
							Eventtittel: Eventtittel,
							Eventtype: Eventtype,
							Description: Description,
							OrganizedBy: OrganizedBy,
							Image: Image,
							EventStartdDate: EventStartdDate,
							EventEndDate: EventEndDate,
							EventStartTime: EventStartTime,
							EventEndTime: EventEndTime,
							ContactEmail: ContactEmail,
							EventLink: EventLink,
						})
				}

			c.HTML(http.StatusOK, "calender.html", gin.H{
				"events": events,
			})

		})

		r.GET("/calender", func(c *gin.Context) {
			

			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT)"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

			if _, err := db.Exec("CREATE TABLE IF NOT EXISTS archive (id SERIAL PRIMARY KEY, eventtittel varchar(45) NOT NULL, eventtype varchar(45) NOT NULL, description varchar(255) NOT NULL, organizedby varchar(45), image TEXT, location GEOMETRY(POINT,4326), geofence GEOGRAPHY, displayfrom DATE, displaytill DATE, eventstartdate DATE, eventenddate DATE, eventstarttime TIME, eventendtime TIME, contactemail TEXT, eventlink TEXT, latitude NUMERIC(10,8), longitude NUMERIC(11,8))"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

			if _, err := db.Exec("WITH moved_rows AS ( DELETE FROM events WHERE displaytill < now() RETURNING * ) INSERT INTO archive SELECT * FROM moved_rows"); 
			err != nil {
				c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
			}

				rows, err := db.Query("SELECT id, eventtittel, eventtype, description, organizedby, image, TO_CHAR(eventstartdate, 'DD Mon YYYY'), TO_CHAR(eventenddate , 'DD Mon YYYY'), TO_CHAR(eventstarttime, 'HH24:MI'), TO_CHAR(eventendtime, 'HH24:MI'), contactemail, eventlink FROM events ORDER BY eventstartdate ASC, eventstarttime ASC")
				if err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error reading Events: %q", err))
					return
				}
	
				defer rows.Close()
				var Id string
				var Eventtittel string 
				var Eventtype string
				var Description string 
				var Image string 
				var OrganizedBy string 
				var EventStartdDate string
				var EventEndDate string
				var EventStartTime string
				var EventEndTime string
				var ContactEmail string
				var EventLink string

				
				events := make([]Event, 0)

				for rows.Next() {
					
					if err := rows.Scan(&Id, &Eventtittel, &Eventtype, &Description, &OrganizedBy, &Image, &EventStartdDate, &EventEndDate, &EventStartTime, &EventEndTime, &ContactEmail, &EventLink); 
					err != nil {
						c.String(http.StatusInternalServerError,
							fmt.Sprintf("Error scanning events: %q", err))
						return
					}
					
					events = append(events, Event{
							Id: Id,
							Eventtittel: Eventtittel,
							Eventtype: Eventtype,
							Description: Description,
							OrganizedBy: OrganizedBy,
							Image: Image,
							EventStartdDate: EventStartdDate,
							EventEndDate: EventEndDate,
							EventStartTime: EventStartTime,
							EventEndTime: EventEndTime,
							ContactEmail: ContactEmail,
							EventLink: EventLink,
						})
				}

			c.HTML(http.StatusOK, "calender.html", gin.H{
				"events": events,
			})
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
