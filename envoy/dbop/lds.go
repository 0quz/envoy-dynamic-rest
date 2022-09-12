package dbop

import (
	"envoy/redis"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Definition of request variables.
type ListenerRequestJson struct {
	Name      string `json:"name"`
	CdsName   string `json:"cds_name"`
	PortValue int    `json:"port_value"`
}

// Sleep func for the deployment process
func waitForNewConfigInit() {
	// Set CDS, LDS, EDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("cdsDeployed", "no")
	fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
	time.Sleep(2 * time.Second)
	redis.SetRedisMemcached("ldsDeployed", "no")
	redis.SetRedisMemcached("edsDeployed", "no")
}

func AddLds(l *ListenerRequestJson, db *gorm.DB) error {
	// Get the first matching CDS from DB.
	var cds Cds
	err := db.Table("cds").Where("name = ?", l.CdsName).Preload("Eds").First(&cds).Error
	if err != nil {
		return err
	}
	// Create a LDS and add the matching CDS to DB.
	err = db.Create(&Lds{Name: l.Name, PortValue: l.PortValue, CdsName: l.CdsName, Cds: cds}).Error
	if err != nil {
		return err
	}
	go waitForNewConfigInit() // Suppose you add a listener when the DB is empty. Envoy can't take the configuration of cds properly so I need to use sleep for 2 seconds.
	return nil
}

func UpdateLds(l *ListenerRequestJson, db *gorm.DB) error {
	// Get the first matching LDS from DB.
	err := db.Table("lds").Where("name = ?", l.Name).First(&Lds{}).Error
	if err != nil {
		return err
	}
	// Get the first matching CDS from DB.
	var cds Cds
	err = db.Table("cds").Where("name = ?", l.CdsName).Preload("Eds").First(&cds).Error
	if err != nil {
		return err
	}
	// Update the LDS table when it matches.
	err = db.Model(&Lds{}).Where("name = ?", l.Name).Updates(map[string]interface{}{"cds_name": l.CdsName, "port_value": l.PortValue, "Cds": cds}).Error
	if err != nil {
		return err
	}
	// Set LDS deployed status no to let Envoy new configuration.
	go waitForNewConfigInit()
	return nil
}

func DeleteLds(l *ListenerRequestJson, db *gorm.DB) error {
	// Delete the matching LDS.
	err := db.Table("lds").Where("name = ?", l.Name).Delete(&Lds{}).Error
	if err != nil {
		return err
	}
	// Set LDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("ldsDeployed", "no")
	return nil
}
