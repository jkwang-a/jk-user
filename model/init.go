package model

import (
	"fmt"
	"os/exec"
	"strings"

	//	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	tablePrefix string // 表前缀
)

func Init() {
	dsn := beego.AppConfig.String("db_url")
	if dsn == "" {
		dsn = "root:123456@tcp(127.0.0.1:3306)/dns?charset=utf8"
	}
	orm.RegisterDataBase("default", "mysql", dsn)
	tablePrefix = "tb_"
	orm.RegisterModelWithPrefix(tablePrefix)
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	orm.RunCommand()
	orm.SetMaxIdleConns("default", 100)
	orm.SetMaxOpenConns("default", 2000)
}

// 返回真实表名
func tableName(name string) string {

	return tablePrefix + name
}
func GetEncodePwd(salt, pwd string) string {
	if salt == "" {
		//	salt = GetRandomString(8)
	}
	var encode string
	en_text := fmt.Sprintf("$1$%s", salt)
	sql := fmt.Sprintf(`SELECT ENCRYPT("%s","%s")`, pwd, en_text) //用户输入的密码+$1$+数据库查出来的密码的后面数据
	orm.NewOrm().Raw(sql).QueryRow(&encode)

	if encode == "" || len(encode) < 32 {
		cmd := exec.Command("openssl", "passwd", "-1", "-salt", salt, pwd)
		out, _ := cmd.Output()
		encode = string(out) //包含\r\n
		encode = strings.Replace(encode, "\r\n", "", 10)
		encode = strings.Replace(encode, "\n", "", 10)
	}
	beego.Debug(encode)
	return encode
}
