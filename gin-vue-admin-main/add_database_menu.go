package main

import (
	"fmt"
	"log"

	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func main() {
	// 初始化配置和数据库连接
	global.GVA_VP = core.Viper()
	initialize.OtherInit()
	global.GVA_LOG = core.Zap()
	global.GVA_DB = initialize.Gorm()
	
	db := global.GVA_DB

	// 1. 添加父级菜单：数据库管理
	parentMenu := system.SysBaseMenu{
		ParentId:  0,
		Path:      "/database",
		Name:      "Database",
		Hidden:    false,
		Component: "view/database/index.vue",
		Sort:      10,
		Meta: system.Meta{
			ActiveName:  "",
			KeepAlive:   true,
			DefaultMenu: false,
			Title:       "数据库管理",
			Icon:        "database",
			CloseTab:    false,
		},
	}

	result := db.Create(&parentMenu)
	if result.Error != nil {
		log.Fatal("创建父菜单失败:", result.Error)
	}
	fmt.Printf("创建父菜单成功，ID: %d\n", parentMenu.ID)

	// 2. 添加子菜单：数据库列表
	childMenu := system.SysBaseMenu{
		ParentId:  parentMenu.ID,
		Path:      "list",
		Name:      "DatabaseList",
		Hidden:    false,
		Component: "view/database/index.vue",
		Sort:      1,
		Meta: system.Meta{
			ActiveName:  "",
			KeepAlive:   true,
			DefaultMenu: false,
			Title:       "数据库列表",
			Icon:        "list",
			CloseTab:    false,
		},
	}

	result = db.Create(&childMenu)
	if result.Error != nil {
		log.Fatal("创建子菜单失败:", result.Error)
	}
	fmt.Printf("创建子菜单成功，ID: %d\n", childMenu.ID)

	// 3. 为超级管理员角色添加菜单权限
	// 假设超级管理员角色ID为888
	authorityId := "888"

	// 为父菜单添加权限
	parentAuthorityMenu := system.SysAuthorityMenu{
		AuthorityId: authorityId,
		MenuId:      fmt.Sprintf("%d", parentMenu.ID),
	}
	result = db.Create(&parentAuthorityMenu)
	if result.Error != nil {
		log.Fatal("为父菜单添加权限失败:", result.Error)
	}
	fmt.Println("为父菜单添加权限成功")

	// 为子菜单添加权限
	childAuthorityMenu := system.SysAuthorityMenu{
		AuthorityId: authorityId,
		MenuId:      fmt.Sprintf("%d", childMenu.ID),
	}
	result = db.Create(&childAuthorityMenu)
	if result.Error != nil {
		log.Fatal("为子菜单添加权限失败:", result.Error)
	}
	fmt.Println("为子菜单添加权限成功")

	fmt.Println("数据库管理菜单添加完成！")
	fmt.Println("请重启 gin-vue-admin 服务，然后登录查看菜单。")
} 