package dbop

import (
	"envoy/redis"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

// Definition of request variables.
type EndpointRequestJson struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	PortValue int    `json:"port_value"`
}

func AddEds(e *EndpointRequestJson, db *gorm.DB) error {
	// Create an EDS.
	err := db.Create(&Eds{Name: e.Name}).Error
	if err != nil {
		return err
	}
	// Set CDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}

func DeleteEds(e *EndpointRequestJson, db *gorm.DB) error {
	// Delete the matching EDS.
	err := db.Table("eds").Where("name = ?", e.Name).Delete(&EndpointAddress{}).Error
	if err != nil {
		return err
	}
	// Set CDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}

func AddEndpointAddress(e *EndpointRequestJson, db *gorm.DB) error {
	// Throw an error If the port already exists.
	err := db.Table("endpoint_addresses").Where("eds_name = ?", e.Name).Where("port_value = ?", e.PortValue).First(&EndpointAddress{}).Error
	if err == nil {
		return errors.New("Eds: " + e.Name + " is already using " + strconv.Itoa(e.PortValue))
	}
	// Create an EndpointAddress.
	err = db.Create(&EndpointAddress{EdsName: e.Name, Address: e.Address, PortValue: e.PortValue}).Error
	if err != nil {
		return err
	}
	// Set CDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}

func DeleteEndpointAddress(e *EndpointRequestJson, db *gorm.DB) error {
	// Delete the matching DeleteEndpointAddress.
	err := db.Table("endpoint_addresses").Where("eds_name = ?", e.Name).Where("port_value = ?", e.PortValue).First(&EndpointAddress{}).Delete(&EndpointAddress{}).Error
	if err != nil {
		return err
	}
	// Set EDS deployed status no to let Envoy new configuration.
	redis.SetRedisMemcached("edsDeployed", "no")
	return nil
}
