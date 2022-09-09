package dbop

import (
	"envoy/redis"
	"errors"
)

type ClusterRequestJson struct {
	Name    string `json:"name"`
	LdsName string `json:"lds_name"`
	EdsName string `json:"eds_name"`
}

func AddCds(c *ClusterRequestJson) error {
	db := ConnectPostgresClient()
	var eds Eds
	err := db.Table("eds").Where("name = ?", c.EdsName).First(&eds).Error
	if err != nil {
		return err
	}
	db.AutoMigrate(&Cds{})
	err = db.Create(&Cds{Name: c.Name, LdsName: c.LdsName, EdsName: c.EdsName, Eds: eds}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("cdsDeployed", "no")
	redis.SetRedisMemcached("ldsDeployed", "no")
	return nil
}

func UpdateCds(c *ClusterRequestJson) error {
	db := ConnectPostgresClient()
	err := db.Model(&Cds{}).Where("name = ?", c.Name).First(&Cds{}).Error
	if err != nil {
		return errors.New("Cds: " + c.Name + " is not found.")
	}
	err = db.Table("cds").Where("eds_name = ?", c.EdsName).First(&Cds{}).Error
	if err == nil {
		return errors.New("Eds: " + c.EdsName + " is already binded.")
	}
	var eds Eds
	err = db.Table("eds").Where("name = ?", c.EdsName).First(&eds).Error
	if err != nil {
		return err
	}
	err = db.Model(&Cds{}).Where("name = ?", c.Name).Updates(map[string]interface{}{"lds_name": c.LdsName, "eds_name": c.EdsName, "Eds": eds}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("cdsDeployed", "no")
	redis.SetRedisMemcached("ldsDeployed", "no")
	return nil
}

func DeleteCds(c *ClusterRequestJson) error {
	db := ConnectPostgresClient()
	var cds Cds
	err := db.Table("cds").Where("name = ?", c.Name).First(&cds).Error
	if err != nil {
		return err
	}
	if cds.LdsName != "" {
		return errors.New("Binded cds cannot be deleted. Lds Name :" + cds.LdsName)
	} else {
		err = db.Table("cds").Where("name = ?", c.Name).Delete(&Cds{}).Error
		if err != nil {
			return err
		}
		redis.SetRedisMemcached("cdsDeployed", "no")
		return nil
	}
}
