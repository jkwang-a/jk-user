package model

import (
	"encoding/json"
	"jk-common/dao"
	"jk-common/util"
	"strconv"

	"github.com/astaxie/beego"

	//	model "jk-user/model"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

//用户表
type Users struct {
	Id          int
	Name        string
	Password    string
	Email       string
	Parent_id   int
	Role_id     int
	Client_id   int
	Status      int
	Api_key     string
	Org_id      int
	Create_time string
	Update_time string
}

type Users_Ext struct {
	Id           int
	Name         string
	Password     string `json:"-"`
	Email        string
	Role_id      int
	Parent_id    int
	Org_id       int
	Role_name    string
	Org_name     string
	Org_fullname string
	Menu_list    string //菜单权限列表
	Op_list      string //操作权限
	Api_key      string
	Create_time  string
	Update_time  string
	Status       int
}

type Access_list struct {
	Access_count int
	Date_time    string
	User_name    string
	User_count   int
}

func User_List(key string, pageindex, pagesize int, order_field, desc string, user_id int) ([]Users_Ext, int, int) {
	org_list := ""
	role_name := ""
	sql_user := `select b.org_list,b.comment from tb_sys_users a 
				left join tb_sys_role b on a.role_id = b.id where a.id=?`
	dao.GetOrm().Raw(sql_user, user_id).QueryRow(&org_list, &role_name)
	org_list = strings.Trim(org_list, "[^\"]+")
	sql := `select a.*,b.name as org_name,c.name as role_name from tb_sys_users a
	   left join tb_sys_org b on b.id=a.org_id
	   left join tb_sys_role c on c.id=a.role_id where 1=1`
	if org_list != "" && role_name != "admin" && role_name != "" {
		sql += fmt.Sprintf(" and a.org_id in (%s) ", org_list)
	}
	if key != "" {
		sql += fmt.Sprintf(" and (a.name like '%s%%' or a.email like '%s%%')", key, key)
	}
	sql += "  order by id desc"
	list := []Users_Ext{}
	dao.GetOrm().Raw(sql).QueryRows(&list)
	total_page, total_rd, _ := dao.QueryPage(sql, pageindex, pagesize, order_field, desc, &list, Users_Ext{})
	return list, total_page, total_rd
}

func User_Add(u *Users) (int64, error) {
	sql := `select id from tb_sys_users where email=?`
	id := 0
	dao.GetOrm().Raw(sql, u.Email).QueryRow(&id)
	if id > 0 {
		u.Id = id
		return 0, errors.New("USER_EXSIT")
	}
	u.Api_key = util.GetUUID()
	u.Status = 0
	rand := util.GetRandomString(8)
	hash := rand + u.Password
	md5 := util.GetMd5Base64([]byte(hash))
	u.Password = fmt.Sprintf("$1$%s$%s", rand, md5)

	//content := fmt.Sprintf("USER-[%s]", u.Name)
	//
	return dao.Insert("tb_sys_users", u)
}
func User_Update(u *Users) error {
	search := dao.SearchMap{}
	search.Put("id", u.Id)
	_, e := dao.Update("tb_sys_users", search, u, "api_key,password")
	return e
}

func User_Del(id int) error {
	_, e := dao.GetOrm().Raw(`delete from tb_sys_users where id=?`, id).Exec()
	return e
}

func User_Get(Id int) (*Users, error) {
	sql := `select * from tb_sys_users where id=?`
	user := Users{}
	dao.GetOrm().Raw(sql, Id).QueryRow(&user)
	if user.Id == 0 {
		return nil, errors.New("USER_NOT_EXIST")
	}
	return &user, nil
}

func Login(username, passwd string) (*Users_Ext, error) {
	fmt.Println("====== 来自网关的调用 ", username, passwd)

	sql := `select a.*,b.name as org_name,c.name as role_name from tb_sys_users a
	   left join tb_sys_org b on b.id=a.org_id
	   left join tb_sys_role c on c.id=a.role_id
	   where a.email=? and a.status=1`
	var user Users_Ext
	//取对应用户
	err := dao.GetOrm().Raw(sql, username).QueryRow(&user)
	if err != nil {
		return &user, errors.New("ERR_LOG")
	}
	u_pass := user.Password //
	fmt.Println("------ ", u_pass)
	inf := strings.Split(u_pass, "$")
	if len(inf) < 2 {
		return nil, errors.New("ERR_PWD")
	}
	//$ 1 $ jKGeTfew $ ztNROU0rxEzw7tpiFJPlsA==
	salt := inf[2]
	hash := salt + passwd
	//e_pass := GetEncodePwd(salt, passwd)
	md5 := util.GetMd5Base64([]byte(hash)) //这个结果就是第三段
	u_pass = inf[3]
	fmt.Println("===", md5, u_pass) //数据库存储密码的第2段和用户输入密码的 解析结果，数据库存储密码的第3段
	if md5 != u_pass {
		fmt.Println("====== 不想等===")
		return nil, errors.New("ERR_PWD")
	}
	//	content := fmt.Sprintf("USER-[%s]", username)
	menu_list, op_list := User_Get_Right(user.Id)
	user.Menu_list = menu_list
	user.Op_list = op_list
	//	model.Log_Save(user.Id, "LOGIN", content)
	fmt.Println(util.ObjToString(user))
	return &user, nil
}

func User_GetByApiKey(ApiKey string) (*Users, error) {
	sql := "select id from tb_sys_users where api_key=? and status=1"
	var user Users
	//取对应用户
	err := dao.GetOrm().Raw(sql, ApiKey).QueryRow(&user)
	if err != nil {
		return nil, errors.New("ERR_API_KEY")
	}
	return &user, nil
}
func User_Login_AZure(token string) (*Users_Ext, error) {
	list := strings.Split(token, ".")
	if len(list) < 2 {
		return nil, errors.New("ERR_AZURE_TOKEN")
	}
	str := list[1]
	decodeBytes, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.New("ERR_AZURE_TOKEN")
	}
	type AZure_User struct {
		Name        string
		Unique_name string
		Upn         string
	}
	z_user := AZure_User{}
	code := string(decodeBytes)
	util.StringToObj(code, &z_user)
	user := Users_Ext{}
	if z_user.Name != "" && z_user.Unique_name != "" {
		sql := `select a.*,b.name as org_name,c.name as role_name from tb_sys_users a
	   left join tb_sys_org b on b.id=a.org_id
	   left join tb_sys_role c on c.id=a.role_id
	   where a.email=? and a.status=1`
		orm := dao.GetOrm()
		orm.Raw(sql, z_user.Unique_name).QueryRow(&user)
		if user.Id == 0 {
			//新增用户
			info := Users{}
			info.Name = z_user.Name
			info.Email = z_user.Unique_name
			info.Status = 0
			info.Role_id = 0 //默认增加
			info.Org_id = 1
			User_Add(&info)

		}
		orm.Raw(sql, z_user.Unique_name).QueryRow(&user)
		return &user, nil
	}
	return nil, errors.New("ERR_AZURE_TOKEN")
}

//取用户的权限列表
func User_Get_Right(id int) (string, string) {
	role_id := 0
	org_list := ""
	modules_list := ""
	sql := `select a.role_id ,b.org_list,b.module_list from tb_sys_users a
	 left join tb_sys_role b on b.id=a.role_id
	where a.id=?`
	dao.GetOrm().Raw(sql, id).QueryRow(&role_id, &org_list, &modules_list)
	org_users := []int{}
	if org_list != "" {
		sql = `select id from tb_sys_users where org_id in (?)`
		dao.GetOrm().Raw(sql, org_list).QueryRows(&org_users)
	}
	//取用户菜单权限
	menu_list := []int{}
	sql = `select * from tb_sys_modules where item_type=0 and id in (?)`
	dao.GetOrm().Raw(sql, modules_list).QueryRows(&menu_list)
	menu_ids := ""
	for k, v := range menu_list {
		menu_ids += fmt.Sprintf("%d", v)
		if k < len(menu_list)-1 {
			menu_ids += ","
		}
	}
	//取用户操作权限
	op_list := []int{}
	op_keys := ""
	sql = `select key_code from tb_sys_modules where item_type=1 and id in (?)`
	dao.GetOrm().Raw(sql, modules_list).QueryRows(&op_list)
	for k, v := range op_list {
		op_keys += fmt.Sprintf("%d", v)
		if k < len(op_list)-1 {
			op_keys += ","
		}
	}
	fmt.Println(menu_ids, op_keys)
	return menu_ids, op_keys
}

//取指定用户名下管理的用户
func User_Get_Mgr_userList(id int) interface{} {
	type Info struct {
		Id   int
		Name string
	}
	role_id := 0
	org_list := ""
	modules_list := ""
	sql := `select a.role_id ,b.org_list,b.module_list from tb_sys_users a
	 left join tb_sys_role b on b.id=a.role_id
	where a.id=?`
	dao.GetOrm().Raw(sql, id).QueryRow(&role_id, &org_list, &modules_list)
	org_users := []Info{}
	org_list += "," + fmt.Sprintf("%d", id)
	if org_list != "" {
		sql = `select id,name from tb_sys_users where org_id in (?)`
		dao.GetOrm().Raw(sql, org_list).QueryRows(&org_users)
	}
	/*
		user_ids := ""
		for k, v := range org_users {
			user_ids += fmt.Sprintf("%s", v)
			if k < len(org_users)-1 {
				user_ids += ","
			}
		}*/
	return org_users
}

func User_GetById(Id string) (*Users_Ext, error) {
	sql := `select a.*,b.comment as role_name from tb_sys_users a 
			left join tb_sys_role b on a.role_id = b.id where a.id=?`
	var user Users_Ext
	//取对应用户
	err := dao.GetOrm().Raw(sql, Id).QueryRow(&user)
	if err != nil {
		return nil, errors.New("ERR_ID")
	}
	return &user, nil
}

func AccessStatistics() []Access_list {
	access_list := []Access_list{}
	sql := `select count(*) as access_count,DATE_FORMAT(create_time,'%Y-%m-%d') as date_time from tb_sys_log 
	where action like '%LOGIN%' and TO_DAYS(NOW())-TO_DAYS(create_time)<=15 GROUP BY DATE_FORMAT(create_time,'%Y-%m-%d') `
	dao.GetOrm().Raw(sql).QueryRows(&access_list)
	for k, v := range access_list {
		userName := ""
		userCount := 0
		date_str := v.Date_time
		sql_user := fmt.Sprintf("select count(*) as user_count, b.name as user_name from tb_sys_log a left join tb_sys_users b on b.id = a.user_id where a.action like '%s%%' and a.create_time like '%s%%' GROUP BY a.user_id ORDER BY user_count desc", "LOGIN", date_str)
		dao.GetOrm().Raw(sql_user).QueryRow(&userCount, &userName)
		access_list[k].User_name = userName
		access_list[k].User_count = userCount
	}
	return access_list
}

type All_User struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Api_key         string `json:"api_key"`
	Group_id        int    `json:"group_id"`
	Parent_id       int
	Group_name      string `json:"group_name"`
	Org_id          int    `json:"org_id"`
	Org_name        string `json:"org_name"`
	Parent_org_name string `json:"parent_org_name"`
	Last_login      string `json:"last_login"`
	Last_ip         string `json:"last_ip"`
	Is_delete       int    `json:"is_delete"`
	Status          int    `json:"status"`
	Create_time     string `json:"create_time"`
	Update_time     string `json:"update_time"`
}

type User_Edit struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Api_key     string `json:"api_key"`
	Group_id    int    `json:"group_id"`
	Parent_id   int
	Org_id      int    `json:"org_id"`
	Last_login  string `json:"last_login"`
	Last_ip     string `json:"last_ip"`
	Is_delete   int    `json:"is_delete"`
	Status      int    `json:"status"`
	Create_time string `json:"create_time"`
	Update_time string `json:"update_time"`
}

type Item struct {
	Items []All_User
}

type User_Del_Edit struct {
	Id     []int
	Status int
}

func Get_User_List(page_index int, page_size int, key string, group_id string, status string, org_id string, members string, query_user_id string) (int, interface{}) {

	sql := `select a.*,b.name as group_name,c.name as org_name,d.name as parent_org_name from tb_sys_users a 
			left join tb_sys_user_group b on b.id=a.group_id 
			left join tb_sys_org c on c.id=a.org_id
			left join tb_sys_org d on d.id = c.parent_id 
	 		where 1=1 `
	if key != "" {
		sql += fmt.Sprintf(" and (a.name like '%%" + key + "%%' or a.email like '%%" + key + "%%')")
	}
	if group_id != "" {
		id, _ := strconv.Atoi(group_id)
		sql += fmt.Sprintf(" and a.group_id=%d", id)
	}
	if status != "" {
		st, _ := strconv.Atoi(status)
		sql += fmt.Sprintf(" and a.status=%d", st)
	}
	if org_id != "" {
		orgId, _ := strconv.Atoi(org_id)
		sql += fmt.Sprintf(" and a.org_id=%d", orgId)
	}
	if members != "" {
		sql += " and a.id in (" + strings.Trim(members, "[^\"]+") + ")"
	}
	if query_user_id != "" {
		userId, _ := strconv.Atoi(query_user_id)
		sql += fmt.Sprintf(" and a.id=%d", userId)
	}
	sql += ` and a.is_delete != 1  order by a.name asc`

	list := []All_User{}
	_, total, _ := dao.QueryPage(sql, page_index, page_size, "", "", &list, All_User{})
	item := Item{}
	item.Items = list
	return total, item
}

func Add_User(u *User_Edit) (int64, error) {
	sql := `select id from tb_sys_users where email=?`
	id := 0
	dao.GetOrm().Raw(sql, u.Email).QueryRow(&id)
	if id > 0 {
		u.Id = id
		return 0, errors.New("USER_EXSIT")
	}
	u.Api_key = util.GetUUID()
	u.Status = 0
	rand := util.GetRandomString(8) //8个随机字符

	hash := rand + u.Password
	md5 := util.GetMd5Base64([]byte(hash)) //8个随机字符+，再md5 ,再 $1$8个随机字符$md5结果=用户密码
	u.Password = fmt.Sprintf("$1$%s$%s", rand, md5)

	fmt.Println("----- 创建用户生成密码 -----,", hash, md5, u.Password)

	return dao.Insert("tb_sys_users", u)
}

func Update_User(u *User_Edit) error {
	search := dao.SearchMap{}
	search.Put("id", u.Id)
	if u.Last_login == "" {
		u.Last_login = util.GetUtcTimeStr()
	}
	_, e := dao.Update("tb_sys_users", search, u, "api_key,password")
	return e
}

func Del_User(rd *User_Del_Edit) (e error) {
	for _, v := range rd.Id {
		_, e = dao.GetOrm().Raw(`delete from tb_sys_users where id=?`, v).Exec()
	}
	return e
}

func Edit_User_Status(rd *User_Del_Edit) (e error) {
	for _, v := range rd.Id {
		_, e = dao.GetOrm().Raw(`update tb_sys_users set status=? where id=?`, rd.Status, v).Exec()
	}
	return e
}

func User_Login(username, passwd string, client_ip string) (*User_Edit, error) {
	sql := `select a.* from tb_sys_users a where a.email=? `
	var user User_Edit
	//取对应用户
	err := dao.GetOrm().Raw(sql, username).QueryRow(&user)
	if err != nil {
		return &user, errors.New("ERR_LOG")
	}
	if user.Id != 0 && user.Status == 0 {
		return nil, errors.New("PERMISSION DENIED")
	}
	u_pass := user.Password //
	inf := strings.Split(u_pass, "$")
	if len(inf) < 2 {
		return nil, errors.New("ERR_PWD")
	}
	salt := inf[2]
	//hash := salt + passwd
	e_pass := GetEncodePwd(salt, passwd) //util.GetMd5Base64([]byte(hash)) //
	u_pass = inf[3]
	fmt.Println("=====", e_pass, u_pass)
	if e_pass != user.Password {
		return nil, errors.New("ERR_PWD")
	}
	user.Last_ip = client_ip
	user.Last_login = util.GetUtcTimeStr()
	Update_User(&user)
	user.Password = ""
	fmt.Println(util.ObjToString(user))
	return &user, nil
}

func User_Login_With_AZure(token string, client_ip string) (*User_Edit, error) {
	list := strings.Split(token, ".")
	if len(list) < 2 {
		return nil, errors.New("ERR_AZURE_TOKEN")
	}
	str := list[1]
	decodeBytes, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.New("ERR_AZURE_TOKEN")
	}
	type AZure_User struct {
		Name        string
		Unique_name string
		Upn         string
	}
	z_user := AZure_User{}
	code := string(decodeBytes)
	util.StringToObj(code, &z_user)
	user := User_Edit{}
	if z_user.Name != "" && z_user.Unique_name != "" {
		sql := `select a.* from tb_sys_users a
	   where a.email=? `
		orm := dao.GetOrm()
		orm.Raw(sql, z_user.Unique_name).QueryRow(&user)
		if user.Id != 0 && user.Status == 0 {
			return nil, errors.New("PERMISSION DENIED")
		}
		if user.Id == 0 {
			//新增用户
			//info := User_Edit{}
			user.Name = z_user.Name
			user.Email = z_user.Unique_name
			user.Status = 0
			user.Group_id = 2
			user.Last_ip = client_ip
			user.Create_time = util.GetUtcTimeStr()
			user.Update_time = util.GetUtcTimeStr()
			Add_User(&user)

		} else {
			user.Last_ip = client_ip
			user.Last_login = util.GetUtcTimeStr()
			Update_User(&user)
		}
		orm.Raw(sql, z_user.Unique_name).QueryRow(&user)
		return &user, nil
	}
	return &user, errors.New("ERR_AZURE_TOKEN")
}

func Get_User_Detail_By_Key(ApiKey string) (*User_Edit, error) {
	sql := "select id from tb_sys_users where api_key=? and status=1"
	var user User_Edit
	//取对应用户
	err := dao.GetOrm().Raw(sql, ApiKey).QueryRow(&user)
	if err != nil {
		return nil, errors.New("ERR_API_KEY")
	}
	return &user, nil
}

func Get_UserId_With_ParentId(user_id string) (string, error) {
	var id_str = ""
	sql := "select GROUP_CONCAT(id) from tb_sys_users where parent_id=? or id=?"
	userId, _ := strconv.Atoi(user_id)
	err := dao.GetOrm().Raw(sql, userId, userId).QueryRow(&id_str)
	if err != nil {
		return "", errors.New("ERR")
	}
	return id_str, nil
}

func Update_New_User(email, api_key string) error {
	sql := "update tb_sys_users set api_key=? where email=?"
	_, err := dao.GetOrm().Raw(sql, api_key, email).Exec()
	if err != nil {
		return errors.New("ERR")
	}
	return nil
}

func Update_All_User(data_str string) error {
	var data_map map[string]interface{}
	err := json.Unmarshal([]byte(data_str), &data_map)
	if err != nil {
		beego.Debug(err)
		return err
	}
	var user_map []map[string]interface{}
	user_byte, _ := json.Marshal(data_map["data"])
	err = json.Unmarshal([]byte(string(user_byte)), &user_map)
	if err != nil {
		beego.Debug(err)
		return err
	}
	sql := `select a.id,a.api_key from tb_sys_users a where a.email=?`
	sql_update := `update tb_sys_users set api_key=? where id=?`
	for _, v := range user_map {
		email := v["Email"]
		api_key := v["Api_key"]
		user_id := 0
		old_api_key := ""
		err = dao.GetOrm().Raw(sql, email).QueryRow(&user_id, &old_api_key)
		if user_id > 0 && api_key != old_api_key {
			_, err = dao.GetOrm().Raw(sql_update, api_key, user_id).Exec()
			if err != nil {
				return errors.New("ERR")
			}
		}
	}
	return nil
}
