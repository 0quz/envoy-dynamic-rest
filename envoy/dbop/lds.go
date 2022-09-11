package dbop

import (
	"envoy/redis"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ListenerRequestJson struct {
	Name      string `json:"name"`
	CdsName   string `json:"cds_name"`
	PortValue int    `json:"port_value"`
}

func sleep() {
	fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
	time.Sleep(2 * time.Second)
}

func AddLds(l *ListenerRequestJson, db *gorm.DB) error {
	var cds Cds
	err := db.Table("cds").Where("name = ?", l.CdsName).Preload("Eds").First(&cds).Error
	if err != nil {
		return err
	}
	db.AutoMigrate(&Lds{})
	err = db.Create(&Lds{Name: l.Name, PortValue: l.PortValue, CdsName: l.CdsName, Cds: cds}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("cdsDeployed", "no")
	sleep() // Suppose you add a listener when the DB is empty. Envoy can't take the configuration of cds properly so I need to use sleep for 2 seconds.
	redis.SetRedisMemcached("ldsDeployed", "no")
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}

func UpdateLds(l *ListenerRequestJson, db *gorm.DB) error {
	err := db.Table("lds").Where("name = ?", l.Name).First(&Lds{}).Error
	if err != nil {
		return err
	}
	var cds Cds
	err = db.Table("cds").Where("name = ?", l.CdsName).Preload("Eds").First(&cds).Error
	if err != nil {
		return err
	}
	err = db.Model(&Lds{}).Where("name = ?", l.Name).Updates(map[string]interface{}{"cds_name": l.CdsName, "port_value": l.PortValue, "Cds": cds}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("ldsDeployed", "no")
	return nil
}

func DeleteLds(l *ListenerRequestJson, db *gorm.DB) error {
	err := db.Table("lds").Where("name = ?", l.Name).Delete(&Lds{}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("ldsDeployed", "no")
	return nil
}
