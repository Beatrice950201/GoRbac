package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/models/model"
	"github.com/beatrice950201/GoRbac/models/model/news"
	_ "github.com/go-sql-driver/mysql"
)

func init()  {
	orm.Debug, _ = beego.AppConfig.Bool("debug")
	username := beego.AppConfig.String("db_username")
	password := beego.AppConfig.String("db_password")
	database := beego.AppConfig.String("db_database")
	port := beego.AppConfig.String("db_port")
	host := beego.AppConfig.String("db_host")
	charset := beego.AppConfig.String("db_charset")
	prefix := beego.AppConfig.String("db_prefix")
	dnsOpt := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s`, username, password,host,port, database,charset,`Asia%2FShanghai`)
	_ = orm.RegisterDataBase(
		"default",
		"mysql",
		dnsOpt,
		30)
	orm.RegisterModelWithPrefix(
		prefix,
		new(model.AdminMenu),
		new(model.AdminUser),
		new(model.AdminRole),
		new(model.AdminConfig),
		new(news.TeamNewsCategory),
		new(news.TeamNewsDocument),
	)
	_ = orm.RunSyncdb("default", false, true)
}
