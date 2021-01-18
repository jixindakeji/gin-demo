USE general ;
INSERT INTO `general_sys_menu`(`id`, `parent_id`, `name`, `url`, `type`, `code`, `icon`, `sort`, `remark`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1,  0, '系统管理',          '/system',     1, 'system',               'el-icon-setting',   1,'',-1, now(), now(), NULL),
(2,  1, '用户管理',          '/system/user',2, 'system-user',          'el-icon-user-solid',1,'',-1, now(), now(), NULL),
(3,  1, '角色管理',          '/system/role',2, 'system-role',          'el-icon-coin',      2,'',-1, now(), now(), NULL),
(4,  1, '菜单管理',          '/system/menu',2, 'system-menu',          'el-icon-menu',      3,'',-1, now(), now(), NULL),
(5,  1, '操作日志',          '/system/audit',2, 'system-audit',        'el-icon-tickets',   4,'',-1, now(), now(), NULL),
(6,  2, '查询用户',          '',            3, 'system:user:get',      '',                  1,'',-1, now(), now(), NULL),
(7,  2, '增加用户',          '',            3, 'system:user:post',     '',                  2,'',-1, now(), now(), NULL),
(8,  2, '修改用户',          '',            3, 'system:user:put',      '',                  3,'',-1, now(), now(), NULL),
(9,  2, '删除用户',          '',            3, 'system:user:delete',   '',                  4,'',-1, now(), now(), NULL),
(10, 2, '查询用户角色的映射',  '',            3, 'system:user:role:get', '',                  5,'',-1, now(), now(), NULL),
(11, 2, '绑定用户与角色',     '',            3, 'system:user:role:post','',                  6,'',-1, now(), now(), NULL),
(12, 3, '查询角色',          '',            3, 'system:role:get',      '',                  1,'',-1, now(), now(), NULL),
(13, 3, '增加角色',          '',            3, 'system:role:post',     '',                  2,'',-1, now(), now(), NULL),
(14, 3, '修改角色',          '',            3, 'system:role:put',      '',                  3,'',-1, now(), now(), NULL),
(15, 3, '删除角色',          '',            3, 'system:role:delete',   '',                  4,'',-1, now(), now(), NULL),
(16, 3, '获取角色与菜单映射',  '',            3, 'system:role:menu:get', '',                  5,'',-1, now(), now(), NULL),
(17, 3, '绑定角色与菜单',     '',            3, 'system:role:menu:post','',                  6,'',-1, now(), now(), NULL),
(18, 4, '查询菜单',          '',            3, 'system:menu:get',      '',                  1,'',-1, now(), now(), NULL),
(19, 4, '新增菜单',          '',            3, 'system:menu:post',     '',                  2,'',-1, now(), now(), NULL),
(20, 4, '修改菜单',          '',            3, 'system:menu:put',      '',                  3,'',-1, now(), now(), NULL),
(21, 4, '删除菜单',          '',            3, 'system:menu:delete',   '',                  4,'',-1, now(), now(), NULL),
(22, 5, '查询操作日志',       '',            3, 'system:audit:get',     '',                  4,'',-1, now(), now(), NULL);
INSERT INTO `general_sys_role`(`id`, `name`, `remark`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1,'系统管理员','系统管理员',-1, now(), now(), NULl);

INSERT INTO `general_role_menu`(`sys_role_id`, `sys_menu_id`) VALUES
(1,6),
(1,7),
(1,8),
(1,9),
(1,10),
(1,11),
(1,12),
(1,13),
(1,14),
(1,15),
(1,16),
(1,17),
(1,18),
(1,19),
(1,20),
(1,21),
(1,22);
INSERT INTO `general_user_role`(`sys_user_id`, `sys_role_id`) VALUES (1,1);
