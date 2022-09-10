package dbop

import (
	"envoy/redis"
	"errors"
	"strconv"
)

type EndpointRequestJson struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	PortValue int    `json:"port_value"`
}

func AddEds(e *EndpointRequestJson) error {
	db := ConnectPostgresClient()
	db.AutoMigrate(&Eds{})
	err := db.Create(&Eds{Name: e.Name}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}

func DeleteEds(e *EndpointRequestJson) error {
	db := ConnectPostgresClient()
	db.AutoMigrate(&Eds{})
	err := db.Table("eds").Where("name = ?", e.Name).Delete(&EndpointAddress{}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}

func AddEndpointAddress(e *EndpointRequestJson) error {
	db := ConnectPostgresClient()
	db.AutoMigrate(&EndpointAddress{})
	err := db.Table("endpoint_addresses").Where("eds_name = ?", e.Name).Where("port_value = ?", e.PortValue).First(&EndpointAddress{}).Error
	if err == nil {
		return errors.New("Eds: " + e.Name + " is already using " + strconv.Itoa(e.PortValue))
	}
	err = db.Create(&EndpointAddress{EdsName: e.Name, Address: e.Address, PortValue: e.PortValue}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}

func DeleteEndpointAddress(e *EndpointRequestJson) error {
	db := ConnectPostgresClient()
	err := db.Table("endpoint_addresses").Where("eds_name = ?", e.Name).Where("port_value = ?", e.PortValue).First(&EndpointAddress{}).Delete(&EndpointAddress{}).Error
	if err != nil {
		return err
	}
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}
