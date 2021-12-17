package model

import (
	"jk-common/dao"
)

type Rule struct {
	Id          int    `json:"id"`
	Module      string `json:"module"`
	Group_id    int    `json:"group_id"`
	Name        string `json:"name"`
	Menu_auth   string `json:"menu_auth"`
	Action_auth string `json:"action_auth"`
	Log_auth    string `json:"log_auth"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	Create_time string `json:"create_time"`
	Update_time string `json:"update_time"`
}

type BatchDel struct {
	Id     []int
	Status int
}

func Rule_List() interface{} {
	list := []Rule{}
	sql := `select * from tb_sys_auth_rule `
	dao.GetOrm().Raw(sql).QueryRows(&list)
	return list
}

//增加规则
func Rule_Add(rd *Rule) (int64, error) {

	id, err := dao.Insert("tb_sys_auth_rule", rd)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//删除规则
func Rule_Del(id int) error {
	_, e := dao.GetOrm().Raw(`delete from tb_sys_auth_rule where id=?`, id).Exec()
	return e
}

//修改规则
func Rule_Update(rd *Rule) error {
	search := dao.SearchMap{}
	search.Put("id", rd.Id)
	_, e := dao.Update("tb_sys_auth_rule", search, rd, "")
	return e
}

//批量删除规则
func Batch_Del(rd *BatchDel) (e error) {
	for _, v := range rd.Id {
		_, e = dao.GetOrm().Raw(`delete from tb_sys_auth_rule where id=?`, v).Exec()
	}
	return e
}

//修改规则状态
func Edit_Rule_Status(rd *BatchDel) (e error) {
	for _, v := range rd.Id {
		_, e = dao.GetOrm().Raw(`update tb_sys_auth_rule set status=? where id=?`, rd.Status, v).Exec()
	}
	return e
}
