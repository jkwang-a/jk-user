package model

import (
	"fmt"
	"jk-common/dao"
	"strconv"
)

type UserGroup struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	System      int    `json:"system"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
	Create_time string `json:"create_time"`
	Update_time string `json:"update_time"`
}

type UpdateStatus struct {
	Id     []int
	Status int
}

func User_Group_List(status string) interface{} {
	list := []UserGroup{}
	sql := `select a.* from tb_sys_user_group a where 1=1`
	if status != "" {
		st, _ := strconv.Atoi(status)
		sql += fmt.Sprintf(" and a.status=%d", st)
	}
	dao.GetOrm().Raw(sql).QueryRows(&list)
	return list
}

//增加用户组
func User_Group_Add(rd *UserGroup) (int64, error) {

	id, err := dao.Insert("tb_sys_user_group", rd)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//删除用户组
func User_Group_Del(id int) error {
	_, e := dao.GetOrm().Raw(`delete from tb_sys_user_group where id=?`, id).Exec()
	return e
}

//修改用户组
func User_Group_Update(rd *UserGroup) error {
	search := dao.SearchMap{}
	search.Put("id", rd.Id)
	_, e := dao.Update("tb_sys_user_group", search, rd, "")
	return e
}

//修改用户组状态
func Edit_User_Group_Status(rd *UpdateStatus) (e error) {

	for _, v := range rd.Id {
		_, e = dao.GetOrm().Raw(`update tb_sys_user_group set status=? where id=?`, rd.Status, v).Exec()
	}
	return e
}
