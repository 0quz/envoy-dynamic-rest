package dbop

import "errors"

type ClusterRequestJson struct {
	Name    string `json:"name"`
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
	err = db.Create(&Cds{Name: c.Name, EdsName: c.EdsName, Eds: eds}).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateCds(c *ClusterRequestJson) error {
	db := ConnectPostgresClient()
	err := db.Table("cds").Where("eds_name = ?", c.EdsName).First(&Cds{}).Error
	if err == nil {
		return errors.New("Eds: " + c.EdsName + " is already binded.")
	}
	var eds Eds
	err = db.Table("eds").Where("name = ?", c.EdsName).First(&eds).Error
	if err != nil {
		return err
	}
	err = db.Model(&Cds{}).Where("name = ?", c.Name).Updates(map[string]interface{}{"eds_name": c.EdsName, "Eds": eds}).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCds(c *ClusterRequestJson) error {
	db := ConnectPostgresClient()
	err := db.Table("cds").Where("name = ?", c.Name).Delete(&Cds{}).Error
	if err != nil {
		return err
	}
	return nil
}
