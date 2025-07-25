package model

import "gorm.io/gorm"

type Container struct {
	gorm.Model
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
	IP     string `json:"ip"`
	Ports  string `json:"ports"`
	Remark string `json:"remark"`
}
