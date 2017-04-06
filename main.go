package main

import (
	_ "ums/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"ums/models"
)

func main() {

	orm.Debug = true
	orm.RegisterModel(new(models.User))
	//注册驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//注册默认数据库
	orm.RegisterDataBase("default", "mysql", "root:abcd1234@tcp(127.0.0.1:3306)/beego?charset=utf8") //密码为空格式

	o := orm.NewOrm()
	o.Using("blog")
	qs := o.QueryTable(models.User{})

	var counts, _ = qs.Count()

	fmt.Println(counts)

	beego.Run()
}
