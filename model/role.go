package model

import (
	"errors"
	"fmt"
	"jk-common/dao"
	"strings"
)

//角色表
type Role struct {
	Id          int
	Name        string
	Comment     string
	Org_list    string
	Module_list string
	Action_list string
	Status      int
	Create_time string
	Update_time string
}

type Role_Arr struct {
	Id          int
	Name        string
	Comment     string
	Org_list    string
	Org_name    string
	Module_list string
	Module_name string
	Action_list string
	Status      int
	Create_time string
	Update_time string
}

func Role_List(key string, pageindex, pagesize int, order_field, desc string, user_id int) ([]Role_Arr, int, int) {
	role_id := 0
	role_name := ""
	sql_user := `select b.id,b.comment from tb_sys_users a 
				left join tb_sys_role b on a.role_id = b.id where a.id=?`
	dao.GetOrm().Raw(sql_user, user_id).QueryRow(&role_id, &role_name)

	sql := `select * from tb_sys_role where 1=1 `
	list := []Role_Arr{}
	if key != "" {
		sql += fmt.Sprintf(" and name like '%s%%'", key)
	}

	if role_name != "" && role_name != "admin" {
		sql += fmt.Sprintf(" and comment!='admin'")
	}
	total_page, total_rd, _ := dao.QueryPage(sql, pageindex, pagesize, order_field, desc, &list, Role{})
	for k, v := range list {
		module_name := ""
		if v.Module_list != "" {
			sql_module := "select GROUP_CONCAT(b.name) as module_name from tb_sys_modules b  where b.id in (" + strings.Trim(v.Module_list, "[^\"]+") + ") group by b.status"
			dao.GetOrm().Raw(sql_module).QueryRow(&module_name)
			list[k].Module_name = module_name
		} else {
			list[k].Module_name = ""
		}
		org_name := ""
		if v.Org_list != "" {
			sql_org := "select GROUP_CONCAT(b.name) as org_name from tb_sys_org b  where b.id in (" + strings.Trim(v.Org_list, "[^\"]+") + ") group by b.status"
			dao.GetOrm().Raw(sql_org).QueryRow(&org_name)
			list[k].Org_name = org_name
		} else {
			list[k].Org_name = ""
		}

	}
	return list, total_page, total_rd
}
func Role_Add(rd *Role) (int64, error) {
	sql := `select id from tb_sys_role where name=?`
	id := 0
	dao.GetOrm().Raw(sql, rd.Name).QueryRow(&id)
	if id > 0 {
		rd.Id = id
		return 0, errors.New("NAME_EXSIT")
	}
	return dao.Insert("tb_sys_role", rd)
}
func Role_Update(rd *Role) error {
	search := dao.SearchMap{}
	search.Put("id", rd.Id)
	_, e := dao.Update("tb_sys_role", search, rd, "")
	return e
}

func Role_Del(id int) error {
	sql := `select id from tb_sys_users where role_id=?`
	user_id := 0
	dao.GetOrm().Raw(sql, id).QueryRow(&user_id)
	if user_id > 0 {
		return errors.New("ROLE_USER_EXIST")
	}
	_, e := dao.GetOrm().Raw(`delete from tb_sys_role where id=?`, id).Exec()
	return e
}

//更新权限列表
func Role_Right_Update(role_id int, org_list, module_list string) {
	sql := `update tb_sys_role set org_list=? ,module_list = ? where id=?`
	dao.GetOrm().Raw(sql, org_list, module_list, role_id).Exec()
}
