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
	Address   string `gorm:"default:0.0.0.0"` // 192.168.65.2 #host.docker.internal
	PortValue int    `gorm:"default:10000;unique"`
	Cds       Cds    //`gorm:"foreignkey:Name"` //[]Cds `gorm:"many2many:lds_cds;"`
}

type Cds struct {
	Name    string `gorm:"primaryKey"`
	EdsName string `gorm:"unique;not null"`
	Eds     Eds    //[]Eds `gorm:"many2many:cds_eds;"`
}

type Eds struct {
	Name string `gorm:"primaryKey"`
}

type EndpointAddress struct {
	Id        int    `gorm:"primaryKey;autoIncrement"`
	EdsName   string `gorm:"not null"`
	PortValue int    `gorm:"not null"` // 192.168.65.2 #host.docker.internal // I did this purposely
	Address   string `gorm:"default:0.0.0.0"`
}

func Init(c *config.Config) *gorm.DB {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	//db, err := gorm.Open(sqlite.Open("envoy.db"), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	return db
}
