package dbop

import (
	"envoy/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Lds struct {
	//gorm.Model //for creating automatic id / create / update / delete date
	Name      string `gorm:"primaryKey"`
	CdsName   string `gorm:"unique;not null"`
	Address   string `gorm:"default:0.0.0.0"`
	PortValue int    `gorm:"default:10000;unique"`
	Cds       Cds    //`gorm:"foreignkey:Name"` //[]Cds `gorm:"many2many:lds_cds;"`
}

type Cds struct {
	Name    string `gorm:"primaryKey"`
	LdsName string `gorm:"unique;not null"`
	EdsName string `gorm:"unique;not null"`
	Eds     Eds    //[]Eds `gorm:"many2many:cds_eds;"`
}

type Eds struct {
	Name string `gorm:"primaryKey"`
}

type EndpointAddress struct {
	Id        int    `gorm:"primaryKey;autoIncrement"`
	EdsName   string `gorm:"not null"`
	PortValue int    `gorm:"not null"` //`gorm:"unique"` // I did this purposely
	Address   string `gorm:"default:0.0.0.0"`
}

//var dsn = "host=0.0.0.0 user=oguz dbname=envoy port=5432 sslmode=disable TimeZone=Asia/Istanbul"

//var url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "user", c.DBPass, c.DBHost, c.DBPort, c.DBName)

func ConnectPostgresClient() *gorm.DB {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	//db, err := gorm.Open(sqlite.Open("envoy.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return db
}
