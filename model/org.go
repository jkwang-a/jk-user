package model

import (
	"errors"
	"fmt"
	"jk-common/dao"
	"sort"
)

//组织机构

type Org struct {
	Id          int
	Name        string
	Levelcode   string
	Parent_id   int
	Seq_id      int //排序码
	Status      int
	Create_time string
	Update_time string
}

type Org_Ext struct {
	Org
	Flag int
}

type Sort struct {
	Id []int
}

func Org_List1(keys string) []Org {
	sql := `select * from tb_sys_org order by levelcode desc ,seq_id desc`
	list := []Org{}
	dao.GetOrm().Raw(sql).QueryRows(&list)
	return list
}

var org_list []Org_Ext

func Org_List(keys string) interface{} {
	sql := `select * from tb_sys_org order by  levelcode `
	dao.GetOrm().Raw(sql).QueryRows(&org_list)
	return Org_GetTree(org_list, 0)
}

type Tree_Org struct {
	Org
	Children []Tree_Org
}

func Copy_Org_Item(t *Tree_Org, v *Org_Ext) {
	t.Id = v.Id
	t.Name = v.Name
	t.Levelcode = v.Levelcode
	t.Status = v.Status
	t.Seq_id = v.Seq_id
	t.Parent_id = v.Parent_id
	t.Create_time = v.Create_time
	t.Update_time = v.Update_time
}

type SortTreeItem_ORG []Tree_Org

func (s SortTreeItem_ORG) Len() int { return len(s) }

func (s SortTreeItem_ORG) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SortTreeItem_ORG) Less(i, j int) bool { return s[i].Seq_id < s[j].Seq_id }

func Org_GetTree(list []Org_Ext, parent_id int) []Tree_Org {
	var tree_list []Tree_Org
	for k, v := range list {
		if list[k].Flag != 0 {
			continue
		}
		if v.Parent_id == parent_id {
			list[k].Flag = 1
			t := Tree_Org{}
			Copy_Org_Item(&t, &v)
			t.Children = Org_GetTree(list, v.Id)

			sort.Sort(SortTreeItem_ORG(t.Children))
			tree_list = append(tree_list, t)
			//fmt.Println(t.Id)
		}
	}
	sort.Sort(SortTreeItem_ORG(tree_list))
	return tree_list
}

func Org_Add(rd *Org) (int64, error) {

	sql := `select id from tb_sys_org where name=?`
	r_id := 0
	dao.GetOrm().Raw(sql, rd.Name).QueryRow(&r_id)
	if r_id > 0 {
		//rd.Id = r_id
		//	return 0, errors.New("NAME_EXSIT")
	}

	id, err := dao.Insert("tb_sys_org", rd)
	if err != nil {
		return 0, err
	}
	sql = `select levelcode from tb_sys_org where id=?`
	parent_code := ""
	dao.GetOrm().Raw(sql, rd.Parent_id).QueryRow(&parent_code)
	code := fmt.Sprintf("%s%.03d", parent_code, id)
	sql = `update tb_sys_org set levelcode =? where id=?`
	dao.GetOrm().Raw(sql, code, id).Exec()
	return id, nil
}

func Org_Update(rd *Org) error {
	search := dao.SearchMap{}
	search.Put("id", rd.Id)
	_, e := dao.Update("tb_sys_org", search, rd, "levelcode")
	return e
}

func Org_Del(id int) error {
	//是否有子节点
	rd := Org{}
	dao.GetOrm().Raw(`select * from tb_sys_org where id=?`, id).QueryRow(&rd)
	if rd.Id == 0 {
		return errors.New("ITEM_NOT_FOUND")
	}
	count := 0
	dao.GetOrm().Raw(`select count(*) from tb_sys_org where parent_id=?`, rd.Id).QueryRow(&count)
	if count > 0 {
		return errors.New("CHILD_EXSIT")
	}
	_, e := dao.GetOrm().Raw(`delete from tb_sys_org where id=?`, rd.Id).Exec()
	return e
}
func Org_GetFullPath(id int) string {
	sql := `select levelcode from tb_sys_org where id=?`
	code := ""
	dao.GetOrm().Raw(sql, id).QueryRow(&code)
	fmt.Println(code)

	c := len(code) / 3
	ids := ""
	for i := 0; i < c; i++ {
		p_code := code[i*3 : (i+1)*3]
		v := ""
		begin := 0
		for k := 0; k < len(p_code); k++ {

			if p_code[k:k+1] != "0" || begin == 1 {
				v += p_code[k : k+1]
				begin = 1
			}
		}
		//	fmt.Println(v)
		ids += v
		ids += ","
	}
	if len(ids) > 0 {
		ids += "0"
	}
	sql = fmt.Sprintf("select id,name from tb_sys_org where id in  (%s) order by  levelcode ", ids)
	list := []Org{}
	dao.GetOrm().Raw(sql).QueryRows(&list)
	path := ""
	//	path_key := ""
	for k, v := range list {
		path += v.Name
		//	path_key += v.Key_code
		if k < len(list)-1 {

			path += "/"
			//		path_key += "/"
		}
	}

	return path
}

//移动节点到新
//子节点跟着移动到新的节点主要是要改变code
func Org_move(node_id int, to_parent_id int) error {
	orm := dao.GetOrm()
	sql := `select levelcode from tb_sys_org where id=?`
	src_code := ""
	orm.Raw(sql, node_id).QueryRow(&src_code)
	fmt.Println(src_code)
	sql = `select levelcode from tb_sys_org where id=?`
	parent_code := ""
	orm.Raw(sql, to_parent_id).QueryRow(&parent_code)
	fmt.Println(parent_code)
	l_src := len(src_code)

	src_new_code := fmt.Sprintf("%s%.03d", parent_code, node_id)
	fmt.Println(src_new_code)
	sql = `update tb_sys_org set levelcode=?,parent_id= ? where id=?`
	orm.Raw(sql, src_new_code, to_parent_id, node_id).Exec()
	//取出来源的所有子节点
	child_list := []Org{}
	orm.Raw(`select *  from tb_sys_org where parent_id=?`, node_id).QueryRows(&child_list)
	for _, v := range child_list {
		sub_code := src_new_code + v.Levelcode[l_src:len(v.Levelcode)]
		fmt.Println(sub_code)
		orm.Raw(`update tb_sys_org set levelcode=? where id=? `, sub_code, v.Id).Exec()
	}
	return nil
}

func Get_Org_List() interface{} {
	list := []Org{}
	sql := `select * from tb_sys_org order by  seq_id `
	dao.GetOrm().Raw(sql).QueryRows(&list)
	return list
}

//批量删除组织机构
func Batch_Del_Org(rd *BatchDel) (e error) {
	for _, v := range rd.Id {
		_, e = dao.GetOrm().Raw(`delete from tb_sys_org where id=?`, v).Exec()
	}
	return e
}

//修改组织机构状态
func Edit_Org_Status(rd *BatchDel) (e error) {
	for _, v := range rd.Id {
		_, e = dao.GetOrm().Raw(`update tb_sys_org set status=? where id=?`, rd.Status, v).Exec()
	}
	return e
}

//组织机构排序
func Sort_Org(rd *Sort) (e error) {
	sort := 0
	for _, v := range rd.Id {
		sort = sort + 1
		_, e = dao.GetOrm().Raw(`update tb_sys_org set seq_id=? where id=?`, sort, v).Exec()
	}
	return e
}
