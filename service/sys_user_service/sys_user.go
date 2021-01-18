package sys_user_service

import (
	"gin-demo/models"
	"sort"
)

type SysUser struct {
	ID uint

	Username    string
	DisplayName string
	Password    string
	OldPassword string
	Email       string
	Phone       string
	Status      int
	Remark      string
	RoleIds     []uint

	Page     int
	PageSize int
}

func (u *SysUser) Get() (*models.SysUser, error) {
	sysUser, err := models.GetSysUser(u.ID)
	if err != nil {
		return nil, err
	}
	return sysUser, nil
}

func (u *SysUser) GetAll() ([]*models.SysUser, error) {
	query := make(map[string]interface{})
	query["username"] = u.Username

	pageSize := u.PageSize
	page := u.Page

	return models.GetSysUsers(query, page, pageSize)
}

func (u *SysUser) Count() (int, error) {
	query := make(map[string]interface{})
	query["username"] = u.Username
	query["status"] = u.Status

	return models.GetSysUserCount(query)
}

func (u *SysUser) Add() error {
	data := map[string]interface{}{
		"username":     u.Username,
		"display_name": u.DisplayName,
		"password":     u.Password,
		"email":        u.Email,
		"phone":        u.Phone,
		"status":       u.Status,
		"remark":       u.Remark,
	}

	if err := models.AddSysUser(data); err != nil {
		return err
	}

	return nil
}

func (u *SysUser) Update() error {
	data := map[string]interface{}{
		"display_name": u.DisplayName,
		"email":        u.Email,
		"phone":        u.Phone,
		"status":       u.Status,
		"remark":       u.Remark,
	}
	if len(u.Password) > 0 {
		data["password"] = u.Password
	}

	if err := models.UpdateSysUser(u.ID, data); err != nil {
		return err
	}

	return nil
}

func (u *SysUser) GetSysUserByUsername() (*models.SysUser, error) {
	user, err := models.GetSysUserByUsername(u.Username)
	if err != nil {
		return nil, nil
	}
	return user, nil
}

func (u *SysUser) Delete() error {
	return models.DeleteUser(u.ID)
}

func (u *SysUser) GetRole() ([]uint, error) {
	var data []uint
	menu, err := models.GetSysRoleByUser(u.ID)
	if err != nil {
		return data, err
	}
	for _, item := range menu {
		data = append(data, item.Model.ID)
	}
	if len(data) == 0 {
		return []uint{}, nil
	}
	return data, nil
}

func (u *SysUser) UpdateRole() error {
	data := make(map[string]interface{})
	data["user_id"] = u.ID
	data["role_ids"] = u.RoleIds

	return models.UpdateSysUserRole(data)
}

type SysMenus []*models.SysMenu

func (m SysMenus) Len() int {
	return len(m)
}

func (m SysMenus) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m SysMenus) Less(i, j int) bool {
	return m[i].Sort < m[j].Sort
}

func (u *SysUser) GetSysMenuByUser() ([]*models.SysMenu, []string, error) {
	menus := []*models.SysMenu{}
	buttons := []string{}

	roles, err := models.GetSysRoleByUser(u.ID)
	if err != nil {
		return nil, nil, err
	}
	for _, role := range roles {
		items, err := models.GetSysMenuByRole(role.ID)
		if err != nil {
			return nil, nil, err
		}
		for _, item := range items {
			menu, button, err := getFatherMenuList(item)
			if err != nil {
				return nil, nil, err
			}
			buttons = append(buttons, button...)
			menus = append(menus, menu)
		}
	}

	sort.Sort(SysMenus(menus))
	data := mergeSameFather(menus)
	return data, buttons, nil
}

func getFatherMenuList(menu *models.SysMenu) (*models.SysMenu, []string, error) {
	var parent *models.SysMenu
	var err error

	buttons := []string{}
	parentId := menu.ParentId
	if parentId == 0 {
		return menu, buttons, nil
	}
	for {
		parent, err = models.GetSysMenu(parentId)
		if err != nil {
			return nil, nil, err
		}
		if parent.Type != 3 && menu.Type != 3 {
			parent.Children = append(parent.Children, menu)
		} else if parent.Type == 3 || menu.Type == 3 {
			buttons = append(buttons, menu.Code)
		}
		if parent.ParentId == 0 {
			break
		} else {
			parentId = parent.ParentId
			menu = parent
		}
	}
	return parent, buttons, nil
}

func mergeSameFather(data []*models.SysMenu) []*models.SysMenu {
	menuTree := []*models.SysMenu{}
	for _, item := range data {
		found := false
		for _, menu := range menuTree {
			if menu.ID == item.ID {
				found = true
				makeMergedTree(menu, item)
			}
		}
		if !found {
			menuTree = append(menuTree, item)

		}
	}
	return menuTree
}

func makeMergedTree(list, item *models.SysMenu) {
	found := false
	if len(list.Children) > 0 && len(item.Children) > 0 {
		for _, menu := range list.Children {
			if menu.ID == item.Children[0].ID {
				found = true
				makeMergedTree(menu, item.Children[0])
			}
		}
		if !found {
			list.Children = append(list.Children, item.Children[0])
			sort.Sort(SysMenus(list.Children))
		}
	} else if len(list.Children) > 0 && len(item.Children) == 0 {
		//Do nothing
	} else if len(list.Children) == 0 && len(item.Children) > 0 {
		list.Children = append(list.Children, item.Children[0])
		sort.Sort(SysMenus(list.Children))
	}
}
