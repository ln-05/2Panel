-- 添加数据库管理菜单到 gin-vue-admin
-- 1. 添加父级菜单：数据库管理
INSERT INTO sys_base_menus (created_at, updated_at, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab, transition_type) 
VALUES (
    NOW(), 
    NOW(), 
    0, 
    '/database', 
    'Database', 
    false, 
    'view/database/index.vue', 
    10, 
    '', 
    true, 
    false, 
    '数据库管理', 
    'database', 
    false, 
    'fade-slide'
);

-- 获取刚插入的父菜单ID
SET @parent_id = LAST_INSERT_ID();

-- 2. 添加子菜单：数据库列表
INSERT INTO sys_base_menus (created_at, updated_at, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab, transition_type) 
VALUES (
    NOW(), 
    NOW(), 
    @parent_id, 
    'list', 
    'DatabaseList', 
    false, 
    'view/database/index.vue', 
    1, 
    '', 
    true, 
    false, 
    '数据库列表', 
    'list', 
    false, 
    'fade-slide'
);

-- 3. 为超级管理员角色添加菜单权限
-- 获取超级管理员角色ID（通常是1）
SET @authority_id = 1;

-- 为父菜单添加权限
INSERT INTO sys_authority_menus (created_at, updated_at, sys_authority_authority_id, sys_base_menu_id) 
VALUES (NOW(), NOW(), @authority_id, @parent_id);

-- 为子菜单添加权限
INSERT INTO sys_authority_menus (created_at, updated_at, sys_authority_authority_id, sys_base_menu_id) 
VALUES (NOW(), NOW(), @authority_id, LAST_INSERT_ID()); 