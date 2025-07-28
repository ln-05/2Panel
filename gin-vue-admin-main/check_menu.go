package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SysBaseMenu struct {
	ID            uint   `gorm:"primarykey"`
	ParentId      uint   `json:"parentId"`
	Path          string `json:"path"`
	Name          string `json:"name"`
	Hidden        bool   `json:"hidden"`
	Component     string `json:"component"`
	Sort          int    `json:"sort"`
	Title         string `json:"title"`
	Icon          string `json:"icon"`
}

func main() {
	// 连接数据库
	dsn := "root:cd423a2933d92b70@tcp(14.103.168.20:3306)/123?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 查询数据库管理相关的菜单
	var menus []SysBaseMenu
	result := db.Where("name LIKE ? OR title LIKE ?", "%Database%", "%数据库%").Find(&menus)
	if result.Error != nil {
		log.Fatal("查询菜单失败:", result.Error)
	}

	fmt.Println("找到的数据库相关菜单:")
	fmt.Println("================================")
	for _, menu := range menus {
		fmt.Printf("ID: %d\n", menu.ID)
		fmt.Printf("父菜单ID: %d\n", menu.ParentId)
		fmt.Printf("路径: %s\n", menu.Path)
		fmt.Printf("名称: %s\n", menu.Name)
		fmt.Printf("标题: %s\n", menu.Title)
		fmt.Printf("组件: %s\n", menu.Component)
		fmt.Printf("图标: %s\n", menu.Icon)
		fmt.Printf("是否隐藏: %t\n", menu.Hidden)
		fmt.Printf("排序: %d\n", menu.Sort)
		fmt.Println("--------------------------------")
	}

	// 查询所有菜单
	var allMenus []SysBaseMenu
	result = db.Order("sort").Find(&allMenus)
	if result.Error != nil {
		log.Fatal("查询所有菜单失败:", result.Error)
	}

	fmt.Println("\n所有菜单列表:")
	fmt.Println("================================")
	for _, menu := range allMenus {
		fmt.Printf("ID: %d | 父ID: %d | 标题: %s | 路径: %s | 组件: %s\n", 
			menu.ID, menu.ParentId, menu.Title, menu.Path, menu.Component)
	}
} 