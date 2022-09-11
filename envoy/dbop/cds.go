package dbop

import (
	"envoy/redis"
	"errors"

	"gorm.io/gorm"
)

type ClusterRequestJson struct {
	Name    string `json:"name"`
	LdsName string `json:"lds_name"`
	EdsName string `json:"eds_name"`
}

func AddCds(c *ClusterRequestJson, db *gorm.DB) error {
	var eds Eds
	err := db.Table("eds").Where("name = ?", c.EdsName).First(&eds).Error
	if err != nil {
		return err
	}
	db.AutoMigrate(&Cds{})
	err = db.Create(&Cds{Name: c.Name, EdsName: c.EdsName, Eds: eds}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("cdsDeployed", "no")
	return nil
}

func UpdateCds(c *ClusterRequestJson, db *gorm.DB) error {
	err := db.Model(&Cds{}).Where("name = ?", c.Name).First(&Cds{}).Error
	if err != nil {
		return errors.New("Cds: " + c.Name + " is not found.")
	}
	err = db.Model(&Cds{}).Where("name = ?", c.Name).Updates(map[string]interface{}{"lds_name": c.LdsName, "eds_name": c.EdsName}).Error
	if err != nil {
		return err
	}
	err = db.Table("cds").Where("eds_name = ?", c.EdsName).First(&Cds{}).Error
	if err == nil {
		redis.SetRedisMemcached("cdsDeployed", "no")
		return errors.New("Eds: " + c.EdsName + " is already binded. Rest of all updated.")
	}
	var eds Eds
	err = db.Table("eds").Where("name = ?", c.EdsName).First(&eds).Error
	if err != nil {
		return err
	}
	err = db.Model(&Cds{}).Where("name = ?", c.Name).Update("Eds", eds).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("cdsDeployed", "no")
	return nil
}

func DeleteCds(c *ClusterRequestJson, db *gorm.DB) error {
	err := db.Table("cds").Where("name = ?", c.Name).Delete(&Cds{}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("cdsDeployed", "no")
	return nil
}
