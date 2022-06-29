package database

import (
	"assignment2/models"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Pras1616"
	dbname   = "postgres"
)

var (
	// db      *gorm.DB
	dbOrder *gorm.DB
	dbItem  *gorm.DB
	err     error
)

func ConnectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	dbOrder, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dbItem, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := dbOrder.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxIdleTime(time.Second * 2)
	sqlDB.SetConnMaxLifetime(time.Second * 2)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	fmt.Println("Success connect to DB using GORM")
	// db.AutoMigrate(&models.Cars{})
	dbOrder.AutoMigrate(&models.Orders{})
	dbItem.AutoMigrate(&models.Items{})
}

// func GetDB() *gorm.DB {
// 	return db
// }

func GetDB_Order() *gorm.DB {
	return dbOrder
}

func GetDB_Item() *gorm.DB {
	return dbItem
}
