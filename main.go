package main

import (
	"os"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set. Cannot connect to Postgres.")
	}
	var err error
	
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to Postgres database")
	}

	db.AutoMigrate(&StringRecord{})

	r := gin.Default()

	r.POST("/strings", createStringHandler)
	r.GET("/strings/:value", getStringHandler)
	r.GET("/strings", getStringsWithFilter)
	r.GET("/strings/filter-by-natural-language", getNaturalLanguageFilter)

	r.DELETE("/strings/:string_value", deleteStringRecord)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)

	fmt.Printf("Server running on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}

