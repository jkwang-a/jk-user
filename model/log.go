package model

import (
	"fmt"
	"jk-common/dao"
	//	"jk-common/util"
	//	"jk-common/util/thread"
)

type Log struct {
	Id          int
	Action      string
	Content     string
	User_id     int
	Create_time string
	Client_ip   string
}
type Log_info struct {
	Id          int
	Action      string
	Content     string
	User_id     int
	User_name   string
	Create_time string
	Client_ip   string
}

type Page_Req struct {
	Page_no   int
	Page_size int
}

//保存日志.

func Log_Save(log Log) {
	/*log := Log{}
	log.User_id = user_id
	log.Action = action
	log.Content = content
	client := thread.Get_Value("client_ip")
	client_ip := "127.0.0.1"
	if client != nil {
		client_ip = client.(string)
	}
	log.Client_ip = client_ip*/
	dao.Insert("tb_sys_log", &log)

}

func Log_Query(page_index int, page_size int, order_field, order string,
	start_time string, end_time string, key string) (int, interface{}) {
	sql := `select a.*,b.name as user_name from tb_sys_log a
	left join tb_sys_users b on b.id=a.user_id where 1=1 `
	if order_field == "" {
		order_field = "create_time"
	}
	if key != "" {
		sql += fmt.Sprintf(" and (a.content like '%%" + key + "%%' or a.action like '%%" + key + "%%')")
	}
	if start_time != "" {
		sql += fmt.Sprintf(" and a.create_time >='%s'", start_time)
	}
	if end_time != "" {
		sql += fmt.Sprintf(" and a.create_time <='%s'", end_time)

	}
	if order == "" {
		order = "desc"
	}
	list := []Log_info{}
	_, total, _ := dao.QueryPage(sql, page_index, page_size, order_field, order, &list, Log_info{})
	return total, list
}
