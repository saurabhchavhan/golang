package main

import (
	//"net/http"
	//"github.com/jinzhu/gorm"
	"github.com/Pallinder/go-randomdata"
	"time"
	"math/rand"
	"database/sql"
	"gopkg.in/gorp.v1"
	"log"
	"fmt"
	 "github.com/umahmood/haversine"
	"strconv"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type NearCabDrivers struct {
	Id        int64  `db:"id" json:"id"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
	Latitude float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
	Distance *float64  `db:"distance" json:"distance"`


}
type CabBookModel struct {
	Id    int `id:"id" json:"id"`
 	PhoneNumber  int `db:"phonenumber" json:"phonenumber" `
 	Firstname string `db:"firstname" json:"firstname" `
 	Lastname  string `db:"lastname" json:"lastname"`
 	PickUpLatitude float64 `db:"pickuplatitude" json:"pickuplatitude"`
  	PickUpLongitude float64 `db:"pickuplongitude" json:"pickuplongitude"`
 	DropLatitude float64 `db:"droplatitude" json:"droplatitude"`
 	DropLongitude float64 `db:"droplongitude" json:"droplongitude"`
 	Distance float64 `db:"distance" "json:"distance"`
 	Date string `db:"date" json:"date"`
 	Driverid int `db:"driverid" json:"driverid"`
 	Status string `db:"status" json:"status"`
 	BookingId int `db:"bookingid" json:"bookingid"`

}

type CabBookModelResponse struct {
 	PhoneNumber  int `db:"phonenumber" json:"phonenumber" `
 	Firstname string `db:"firstname" json:"firstname" `
 	Lastname  string `db:"lastname" json:"lastname"`
 	Distance float64 `db:"distance" "json:"distance"`
 	Date string `db:"date" json:"date"`
 	Driverid int `db:"driverid" json:"driverid"`
 	Status string `db:"status" json:"status"`
 	BookingId int `db:"bookingid" json:"bookingid"`

}
type BookingHistory struct{
	Id    int `id:"id" json:"id"`
 	PhoneNumber  int `db:"phonenumber" json:"phonenumber" `
 	Firstname string `db:"firstname" json:"firstname" `
 	Lastname  string `db:"lastname" json:"lastname"`
 	Distance float64 `db:"distance" "json:"distance"`
 	Date string `db:"date" json:"date"`
 	Driverid int `db:"driverid" json:"driverid"`
 	Status string `db:"status" json:"status"`
 	BookingId int `db:"bookingid" json:"bookingid"`
}
var dbmap = initDb()

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/taxi")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(NearCabDrivers{}, "nearcabdrivers").SetKeys(true, "Id")
	dbmap.AddTableWithName(CabBookModel{},"cabbookmodel").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		v1.GET("/cabs/:longitude/:latitude", GetNearbyCabs)
		v1.GET("/History/:ph", GetHistory)
		v1.POST("/bookcab", BookCab)
		
	}

	r.Run(":8080")
}
func GetNearbyCabs(c *gin.Context) {
	var longtitudeOfUser float64
	longtitudeOfUser,_=strconv.ParseFloat(c.Params.ByName("longitude"), 64)

	var latitudeOfUser float64
	latitudeOfUser,_=strconv.ParseFloat(c.Params.ByName("latitude"), 64)
	//fmt.Println("lat",c.Params.ByName("latitude"))
	//fmt.Println(latitudeOfUser)
	User := haversine.Coord{Lat:latitudeOfUser , Lon:longtitudeOfUser }
    
    DriversDetails,err:=dbmap.Prepare(`insert into nearcabdrivers(firstname,lastname,latitude,longitude) VALUES(?,?,?,?)`)
		
		if err != nil {
		fmt.Println(err.Error())
		} else {
			fmt.Println("Error")
		}
       for i := 0; i < 5; i++ {
       	LLat:=51 +i
       	LLon:=0+i
		_, err = DriversDetails.Exec(randomdata.SillyName(),randomdata.SillyName(),LLat,LLon) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error())
		}
	}
	var Nb []NearCabDrivers
	 _, err = dbmap.Select(&Nb, "SELECT * FROM nearcabdrivers")
     if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
      }
            
        for _,service := range Nb {
           Driver := haversine.Coord{Lat: service.Latitude, Lon: service.Longitude}  // Turin, Italy
           _, km := haversine.Distance(User, Driver)
         
          	result,err:= dbmap.Exec(`UPDATE nearcabdrivers set distance=? where id=?`, km, service.Id);
         	 if err != nil {
				panic(err.Error())
			} else {
				fmt.Println(result.RowsAffected())
					}
        }
            Nb=nil
           _, errx := dbmap.Select(&Nb, "select * from nearcabdrivers order by distance asc limit 2")
        if err != nil {
			panic(err.Error())
		} 
     

	if errx == nil {
		c.JSON(200, Nb)
	} else {
		c.JSON(404, gin.H{"error": "no user(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/users
}

func GetHistory(c *gin.Context) {
	Ph := c.Params.ByName("ph")
	var CBMH []BookingHistory

	_,err := dbmap.Select(&CBMH, "SELECT id, phonenumber,firstname, lastname,distance,date,driverid,status,bookingid  FROM cabbookmodel WHERE phonenumber=?", Ph)
	if err == nil {
		
		c.JSON(200, CBMH)
	
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	// curl -i http://localhost:8080/api/v1/users/1
} 
func BookCab(c *gin.Context) {
	var CBM CabBookModel
	c.Bind(&CBM)

	log.Println(CBM)
	currentTime := time.Now()
	status:="successful"
	driverid:=rand.Intn(10000000)
	bookingid:=rand.Intn(10000000)
	datetoday:=currentTime.String()
	PickupCooridinate:=haversine.Coord{Lat:CBM.PickUpLatitude,Lon:CBM.PickUpLongitude}
	DropCooridinate:=haversine.Coord{Lat:CBM.DropLatitude,Lon:CBM.DropLongitude}
           _, km := haversine.Distance(PickupCooridinate, DropCooridinate)

	if CBM.Firstname != "" && CBM.Lastname != "" {
		if insert, _ := dbmap.Exec(`INSERT INTO CabBookModel (phonenumber,firstname, lastname,pickuplatitude,pickuplongitude,droplatitude,droplongitude,distance,date,driverid,status,bookingid) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, CBM.PhoneNumber,CBM.Firstname, CBM.Lastname,CBM.PickUpLatitude,CBM.PickUpLongitude,CBM.DropLatitude,CBM.DropLongitude,km,datetoday,driverid ,status,bookingid); insert != nil {
			_, err := insert.LastInsertId()
			if err == nil {
				content := &CabBookModelResponse{
					PhoneNumber:CBM.PhoneNumber,
					Firstname: CBM.Firstname,
					Lastname:  CBM.Lastname,
					Distance:km,
					Date :datetoday,
					Driverid:driverid,
					BookingId:bookingid,
					Status:status,
				}
				c.JSON(201, content)
			} else {
				checkErr(err, "Insert failed")
			}
		}

	} else {
		c.JSON(400, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func OptionsUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}