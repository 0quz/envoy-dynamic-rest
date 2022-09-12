package dbop

import (
	"envoy/redis"
	"errors"

	"gorm.io/gorm"
)

// Definition of request variables.
type ClusterRequestJson struct {
	Name    string `json:"name"`
	LdsName string `json:"lds_name"`
	EdsName string `json:"eds_name"`
}

func AddCds(c *ClusterRequestJson, db *gorm.DB) error {
	// Get the first matching EDS from DB.
	var eds Eds
	err := db.Table("eds").Where("name = ?", c.EdsName).First(&eds).Error
	if err != nil {
		return err
	}
	// Create a CDS and add the matching EDS to DB.
	err = db.Create(&Cds{Name: c.Name, EdsName: c.EdsName, Eds: eds}).Error
	if err != nil {
		return err
	}
	// Set CDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("cdsDeployed", "no")
	return nil
}

func UpdateCds(c *ClusterRequestJson, db *gorm.DB) error {
	// Get the first matching CDS from DB.
	err := db.Model(&Cds{}).Where("name = ?", c.Name).First(&Cds{}).Error
	if err != nil {
		return errors.New("Cds: " + c.Name + " is not found.")
	}
	// Get the first matching EDS from DB.
	var eds Eds
	err = db.Table("eds").Where("name = ?", c.EdsName).First(&eds).Error
	if err != nil {
		return err
	}
	// Update the CDS table when it matches.
	err = db.Model(&Cds{}).Where("name = ?", c.Name).Updates(map[string]interface{}{"lds_name": c.LdsName, "eds_name": c.EdsName, "Eds": eds}).Error
	if err != nil {
		return err
	}
	// Set CDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("cdsDeployed", "no")
	return nil
}

func DeleteCds(c *ClusterRequestJson, db *gorm.DB) error {
	// Delete the matching CDS.
	err := db.Table("cds").Where("name = ?", c.Name).Delete(&Cds{}).Error
	if err != nil {
		return err
	}
	// Set CDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("cdsDeployed", "no")
	return nil
}
