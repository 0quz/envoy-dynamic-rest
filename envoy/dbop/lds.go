package dbop

import (
	"envoy/redis"
	"errors"
)

type ListenerRequestJson struct {
	Name      string `json:"name"`
	CdsName   string `json:"cds_name"`
	PortValue int    `json:"port_value"`
}

func AddLds(l *ListenerRequestJson) error {
	db := ConnectPostgresClient()
	var cds Cds
	err := db.Table("cds").Where("name = ?", l.CdsName).Preload("Eds").First(&cds).Error
	if err != nil {
		return err
	}
	if cds.LdsName == l.Name {
		db.AutoMigrate(&Lds{})
		err = db.Create(&Lds{Name: l.Name, PortValue: l.PortValue, CdsName: l.CdsName, Cds: cds}).Error
		if err != nil {
			return err
		}
		redis.SetRedisMemcached("ldsDeployed", "no")
		return nil
	} else {
		return errors.New("Lds cannot be created. Because cds is not binded with lds: " + cds.LdsName)
	}
}

func UpdateLds(l *ListenerRequestJson) error {
	db := ConnectPostgresClient()
	err := db.Table("lds").Where("name = ?", l.Name).First(&Lds{}).Error
	if err != nil {
		return err
	}
	var cds Cds
	err = db.Table("cds").Where("name = ?", l.CdsName).Preload("Eds").First(&cds).Error
	if err != nil {
		return err
	}
	if cds.LdsName == l.Name {
		err = db.Model(&Lds{}).Where("name = ?", l.Name).Updates(map[string]interface{}{"cds_name": l.CdsName, "port_value": l.PortValue, "Cds": cds}).Error // I have to use interface becase of boolean field update
		if err != nil {
			return err
		}
		redis.SetRedisMemcached("ldsDeployed", "no")
		return nil
	} else {
		return errors.New("Lds cannot be updated. Because cds is not binded with lds: " + cds.LdsName)
	}
}

func DeleteLds(l *ListenerRequestJson) error {
	db := ConnectPostgresClient()
	err := db.Table("lds").Where("name = ?", l.Name).Delete(&Lds{}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("ldsDeployed", "no")
	return nil
}
