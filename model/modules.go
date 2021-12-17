package model

import (
	"errors"
	"fmt"
	"jk-common/dao"
	"sort"
	"strings"
)

//组织机构

type Modules struct {
	Id          int
	Name        string
	Url         string
	Levelcode   string
	Parent_id   int
	Seq_id      int //排序码
	Key_code    string
	Item_type   int //0 menu 1,permisson
	Status      int
	Create_time string
	Update_time string
}
type Modules_Ext struct {
	Id          int
	Name        string
	Url         string
	Levelcode   string
	Parent_id   int
	Seq_id      int //排序码
	Key_code    string
	Item_type   int //0 menu 1,permisson
	Status      int
	Flag        int
	Create_time string
	Update_time string
}

type Menu struct {
	Id        int    `json:"id"`
	Parent_id int    `json:"parent_id"`
	Name      string `json:"name"`
	Name_en   string `json:"name_en"`
	Icon      string `json:"icon"`
	Remark    string `json:"remark"`
	Module    string `json:"module"`
	Type      string `json:"type"`
	Url       string `json:"url"`
	Params    string `json:"params"`
	Is_navi   int    `json:"is_navi"`
	Sort      int    `json:"sort"`
	Target    string `json:"target"`
	Status    int    `json:"status"`
}

type Menu_Sort struct {
	Id []int
}

var list []Modules_Ext

func Modules_List(keys string) interface{} {
	sql := `select * from tb_sys_modules order by  levelcode `
	dao.GetOrm().Raw(sql).QueryRows(&list)
	return Modules_GetTree(list, 0)
}

type SortTreeItem []Tree

func (s SortTreeItem) Len() int { return len(s) }

func (s SortTreeItem) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SortTreeItem) Less(i, j int) bool { return s[i].Seq_id < s[j].Seq_id }

type Tree struct {
	Modules
	Children []Tree
}

func Copy_Module_Item(t *Tree, v *Modules_Ext) {
	t.Id = v.Id
	t.Name = v.Name
	t.Url = v.Url
	t.Levelcode = v.Levelcode
	t.Key_code = v.Key_code
	t.Status = v.Status
	t.Item_type = v.Item_type
	t.Seq_id = v.Seq_id
	t.Parent_id = v.Parent_id
	t.Create_time = v.Create_time
	t.Update_time = v.Update_time
}

func Modules_GetTree(list []Modules_Ext, parent_id int) []Tree {
	var tree_list []Tree
	for k, v := range list {
		if list[k].Flag != 0 {
			continue
		}
		if v.Parent_id == parent_id {
			list[k].Flag = 1
			t := Tree{}
			Copy_Module_Item(&t, &v)
			t.Children = Modules_GetTree(list, v.Id)

			sort.Sort(SortTreeItem(t.Children))
			tree_list = append(tree_list, t)
			//fmt.Println(t.Id)
		}
	}
	sort.Sort(SortTreeItem(tree_list))
	return tree_list
}

//增加菜单
func Modules_Add_Menu(rd *Modules) (int64, error) {

	id, err := dao.Insert("tb_sys_modules", rd)
	if err != nil {
		return 0, err
	}
	rd.Item_type = 0
	sql := `select levelcode from tb_sys_modules where id=?`
	parent_code := ""
	dao.GetOrm().Raw(sql, rd.Parent_id).QueryRow(&parent_code)
	code := fmt.Sprintf("%s%.03d", parent_code, id)
	sql = `update tb_sys_modules set levelcode =? where id=?`
	dao.GetOrm().Raw(sql, code, id).Exec()
	return id, nil
}

//增加菜单的权限
func Modules_Add_Permission(rd *Modules) (int64, error) {

	if rd.Parent_id == 0 || rd.Key_code == "" {
		return 0, errors.New("PARAM_ERROR")
	}
	sql := `select id from tb_sys_modules where name=?`
	r_id := 0
	dao.GetOrm().Raw(sql, rd.Name).QueryRow(&r_id)
	if r_id > 0 {
		rd.Id = r_id
		return 0, errors.New("NAME_EXSIT")
	}
	rd.Item_type = 1
	id, err := dao.Insert("tb_sys_modules", rd)
	if err != nil {
		return 0, err
	}
	sql = `select levelcode from tb_sys_modules where id=?`
	parent_code := ""
	dao.GetOrm().Raw(sql, rd.Parent_id).QueryRow(&parent_code)
	code := fmt.Sprintf("%s%.03d", parent_code, id)
	sql = `update tb_sys_modules set levelcode =? where id=?`
	dao.GetOrm().Raw(sql, code, id).Exec()
	return id, nil
}

func Modules_Update(rd *Modules) error {
	search := dao.SearchMap{}
	search.Put("id", rd.Id)
	_, e := dao.Update("tb_sys_modules", search, rd, "levelcode")
	return e
}

func Modules_Del(id int) error {
	//是否有子节点
	rd := Modules{}
	dao.GetOrm().Raw(`select * from tb_sys_modules where id=?`, id).QueryRow(&rd)
	if rd.Id == 0 {
		return errors.New("ITEM_NOT_FOUND")
	}
	count := 0
	dao.GetOrm().Raw(`select count(*) from tb_sys_modules where parent_id=?`, rd.Id).QueryRow(&count)
	if count > 0 {
		return errors.New("CHILD_EXSIT")
	}
	_, e := dao.GetOrm().Raw(`delete from tb_sys_modules where id=?`, rd.Id).Exec()
	return e
}
func Modules_GetFullPath(id int) string {
	sql := `select levelcode from tb_sys_modules where id=?`
	code := ""
	dao.GetOrm().Raw(sql, id).QueryRow(&code)
	fmt.Println(code)
	if code == "" {
		return ""
	}
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
	sql = fmt.Sprintf("select id,name from tb_sys_modules where id in  (%s) order by levelcode ", ids)
	list := []Modules{}
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
func Modules_move(node_id int, to_parent_id int) error {
	orm := dao.GetOrm()
	sql := `select levelcode from tb_sys_modules where id=?`
	src_code := ""
	orm.Raw(sql, node_id).QueryRow(&src_code)
	fmt.Println(src_code)
	sql = `select levelcode from tb_sys_modules where id=?`
	parent_code := ""
	orm.Raw(sql, to_parent_id).QueryRow(&parent_code)
	fmt.Println(parent_code)
	l_src := len(src_code)

	src_new_code := fmt.Sprintf("%s%.03d", parent_code, node_id)
	fmt.Println(src_new_code)
	sql = `update tb_sys_modules set levelcode=?,parent_id= ? where id=?`
	orm.Raw(sql, src_new_code, to_parent_id, node_id).Exec()
	//取出来源的所有子节点
	child_list := []Modules{}
	orm.Raw(`select *  from tb_sys_modules where parent_id=?`, node_id).QueryRows(&child_list)
	for _, v := range child_list {
		sub_code := src_new_code + v.Levelcode[l_src:len(v.Levelcode)]
		fmt.Println(sub_code)
		orm.Raw(`update tb_sys_modules set levelcode=? where id=? `, sub_code, v.Id).Exec()
	}
	return nil
}

func All_Modules_List(keys string) interface{} {
	list = []Modules_Ext{}
	sql := `select * from tb_sys_modules where status=1 order by  levelcode `
	dao.GetOrm().Raw(sql).QueryRows(&list)
	return Modules_GetTree(list, 0)

}

func GetUserMenu(user_id int) interface{} {
	list = []Modules_Ext{}
	modules := ""
	role_name := ""
	sql_module := `select b.module_list,b.comment from tb_sys_users a
	left join tb_sys_role b on b.id=a.role_id
	where a.id=? `
	dao.GetOrm().Raw(sql_module, user_id).QueryRow(&modules, &role_name)
	sql := ""
	if modules != "" && role_name != "admin" {
		sql = "select * from tb_sys_modules where id in (" + strings.Trim(modules, "[^\"]+") + ") order by  levelcode"

	}
	if role_name == "admin" {
		sql = "select * from tb_sys_modules order by  levelcode"
	}
	dao.GetOrm().Raw(sql).QueryRows(&list)
	return list

}

func GetUserAction(user_id int) interface{} {
	list = []Modules_Ext{}
	modules := ""
	role_name := ""
	sql_module := `select b.module_list,b.comment from tb_sys_users a
	left join tb_sys_role b on b.id=a.role_id
	where a.id=? `
	dao.GetOrm().Raw(sql_module, user_id).QueryRow(&modules, &role_name)

	if modules != "" && role_name != "admin" {

		sql := "select * from tb_sys_modules where id in (" + strings.Trim(modules, "[^\"]+") + ") and key_code!='' order by  levelcode"
		dao.GetOrm().Raw(sql).QueryRows(&list)
	}
	if role_name == "admin" {
		sql := "select * from tb_sys_modules where  key_code!='' order by  levelcode"
		dao.GetOrm().Raw(sql).QueryRows(&list)
	}
	return list

}

func Get_Menu(user_id int) interface{} {
	menu := []Menu{}
	sql := "select * from tb_sys_menu order by sort"
	dao.GetOrm().Raw(sql).QueryRows(&menu)
	return menu

}

//增加菜单
func Add_Menu(rd *Menu) (int64, error) {

	id, err := dao.Insert("tb_sys_menu", rd)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//删除菜单
func Del_Menu(id int) error {
	//是否有子节点
	rd := Menu{}
	dao.GetOrm().Raw(`select * from tb_sys_menu where id=?`, id).QueryRow(&rd)
	if rd.Id == 0 {
		return errors.New("ITEM_NOT_FOUND")
	}
	count := 0
	dao.GetOrm().Raw(`select count(*) from tb_sys_menu where parent_id=?`, rd.Id).QueryRow(&count)
	if count > 0 {
		return errors.New("CHILD_EXSIT")
	}
	_, e := dao.GetOrm().Raw(`delete from tb_sys_menu where id=?`, rd.Id).Exec()
	return e
}

//修改菜单
func Edit_Menu(rd *Menu) error {
	search := dao.SearchMap{}
	search.Put("id", rd.Id)
	_, e := dao.Update("tb_sys_menu", search, rd, "")
	return e
}

//菜单启停
func Start_Stop_Menu(rd *Menu) error {
	_, e := dao.GetOrm().Raw(`update tb_sys_menu set status=? where id=?`, rd.Status, rd.Id).Exec()
	return e
}

//根据权限获取菜单
func Get_Menu_With_Auth(user_id int) interface{} {
	menu_auth := ""
	sql_menu := `select a.menu_auth from tb_sys_auth_rule a 
				left join tb_sys_user_group b on a.group_id = b.id 
				left join tb_sys_users c on b.id = c.group_id 
				where c.id=?`
	dao.GetOrm().Raw(sql_menu, user_id).QueryRow(&menu_auth)
	menu := []Menu{}
	sql := "select a.* from tb_sys_menu a where a.id in (" + menu_auth + ") order by a.sort"
	dao.GetOrm().Raw(sql).QueryRows(&menu)
	return menu

}

//菜单排序
func Sort_Menu(rd *Menu_Sort) (e error) {
	sort := 0
	for _, v := range rd.Id {
		sort = sort + 1
		_, e = dao.GetOrm().Raw(`update tb_sys_menu set sort=? where id=?`, sort, v).Exec()
	}
	return e
}
