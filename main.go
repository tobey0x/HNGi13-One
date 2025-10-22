package main

import (
	

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	
	db, err = gorm.Open(sqlite.Open("strings.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&StringRecord{})

	r := gin.Default()

	r.POST("/strings", createStringHandler)
	r.GET("/strings/:value", getStringHandler)
	r.GET("strings", getStringsWithFilter)
	r.GET("/strings/filter-by-natural-language", getNaturalLanguageFilter)

	r.DELETE("/strings/:string_value", deleteStringRecord)


	r.Run(":8080")

}

