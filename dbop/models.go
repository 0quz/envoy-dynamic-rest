package dbop

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Lds struct {
	//gorm.Model //for creating automatic id / create / update / delete date
	Name      string `gorm:"primaryKey"`
	CdsName   string `gorm:"unique"`
	Address   string `gorm:"default:0.0.0.0"`
	PortValue int    `gorm:"default:10000;unique"`
	Deployed  bool   `gorm:"default:false"`
	Cds       Cds    //`gorm:"foreignkey:Name"` //[]Cds `gorm:"many2many:lds_cds;"`
}

type Cds struct {
	Name    string `gorm:"primaryKey"`
	EdsName string `gorm:"unique"`
	Eds     Eds    //[]Eds `gorm:"many2many:cds_eds;"`
}

type Eds struct {
	Name string `gorm:"primaryKey"`
}

type EndpointAddress struct {
	Id        int `gorm:"primaryKey;autoIncrement"`
	EdsName   string
	PortValue int    //`gorm:"unique"` // I did this purposely
	Address   string `gorm:"default:0.0.0.0"`
}

var dsn = "host=localhost user=oguz dbname=oguz port=5432 sslmode=disable TimeZone=Asia/Istanbul"

func ConnectPostgresClient() *gorm.DB {
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open("envoy.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return db
}
