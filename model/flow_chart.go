package model

import (
	"errors"
	"fmt"
	"jk-common/dao"
	"jk-common/util"
)

type Flow_Chart struct {
	Id          int
	Name        string
	Content     string
	User_id     int
	Create_time string
	Update_time string
}

type Flow_Chart_Ext struct {
	Id          int
	Name        string
	Content     string
	User_id     int
	User_name   string
	Create_time string
	Update_time string
}

func Flow_Chart_List(key string, pageindex, pagesize int, order_field, desc string, user_id int) ([]Flow_Chart_Ext, int, int) {
	role_name := ""
	sql_user := `select b.org_list,b.comment from tb_sys_users a 
				left join tb_sys_role b on a.role_id = b.id where a.id=?`
	dao.GetOrm().Raw(sql_user, user_id).QueryRow(&role_name)

	sql := `select a.*,b.name as user_name from tb_sys_flow_chart a 
			left join tb_sys_users b on b.id=a.user_id where 1=1`
	if role_name != "admin" && role_name != "" {
		sql += fmt.Sprintf(" and a.user_id=%d ", user_id)
	}
	if key != "" {
		sql += fmt.Sprintf(" and a.name like '%s%%'", key)
	}
	sql += "  order by id desc"
	list := []Flow_Chart_Ext{}
	dao.GetOrm().Raw(sql).QueryRows(&list)
	total_page, total_rd, _ := dao.QueryPage(sql, pageindex, pagesize, order_field, desc, &list, Users_Ext{})
	return list, total_page, total_rd
}

func Flow_chart_Add(rd *Flow_Chart) (int64, error) {
	sql := `select id from tb_sys_flow_chart where name=?`
	id := 0
	dao.GetOrm().Raw(sql, rd.Name).QueryRow(&id)
	if id > 0 {
		rd.Id = id
		return 0, errors.New("NAME_EXSIT")
	}
	return dao.Insert("tb_sys_flow_chart", rd)
}

func Flow_Chart_Update(rd *Flow_Chart) error {
	update_sql := `update tb_sys_flow_chart set content=?,update_time=? where id=?`
	_, e := dao.GetOrm().Raw(update_sql, rd.Content, util.GetUtcTimeStr(), rd.Id).Exec()
	return e
}

func Flow_Chart_Del(id int) error {
	_, e := dao.GetOrm().Raw(`delete from tb_sys_flow_chart where id=?`, id).Exec()
	return e
}
