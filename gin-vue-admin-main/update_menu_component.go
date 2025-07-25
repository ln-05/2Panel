package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:cd423a2933d92b70@tcp(14.103.168.20:3306)/123?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 更新数据库管理菜单的组件路径
	result := db.Exec("UPDATE sys_base_menus SET component = ? WHERE name = ?", "view/database/index.vue", "Database")
	if result.Error != nil {
		log.Fatal("更新父菜单组件路径失败:", result.Error)
	}
	fmt.Println("更新父菜单组件路径成功")

	// 更新数据库列表菜单的组件路径
	result = db.Exec("UPDATE sys_base_menus SET component = ? WHERE name = ?", "view/database/index.vue", "DatabaseList")
	if result.Error != nil {
		log.Fatal("更新子菜单组件路径失败:", result.Error)
	}
	fmt.Println("更新子菜单组件路径成功")

	fmt.Println("菜单组件路径更新完成！")
	fmt.Println("请重启 gin-vue-admin 服务，然后刷新页面查看效果。")
} 