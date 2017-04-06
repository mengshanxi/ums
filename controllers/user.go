package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"github.com/astaxie/beego/context"
	"time"
	"ums/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) DoLogin() {
	var username = this.GetString("username")
	var password = this.GetString("password")
	fmt.Println(username)
	fmt.Println(password)
	if (strings.EqualFold("user01", username) && strings.EqualFold("123456", password)) {
		this.Data["username"] = username
		this.TplName = "index.html"

	} else {
		this.Ctx.WriteString("用户名或密码错误")
	}

}
func (c *UserController) Login() {
	c.TplName = "login.html"
}

func (c *UserController) Logout() {
	c.TplName = "login.html"
}

func (this *UserController) DoAdd() {
	var username = this.GetString("username")
	var realname = this.GetString("realname")
	u := models.User{Username: username, Realname: realname}
	o := orm.NewOrm()
	o.Insert(&u)
	this.Ctx.WriteString("success")
}

func (this *UserController) Add() {
	this.TplName = "user-add.html"
}
func (this *UserController) Read() {

	aColumns := []string{
		"Id",
		"Username",
		"Realname",
		"Status",
		"Gid",
	}

	user := new(models.User)
	user.Gid = "hello"
	user.Id = 22
	user.Realname = "zhangsan"

	var where = make(map[string]interface{})

	fmt.Println(where)

	maps, count, counts := Datatables(aColumns, user, this.Ctx.Input, where)

	var output = make([][]interface{}, len(maps))
	for i, m := range maps {
		for _, v := range aColumns {
			if v == "Lasttime" {
				output[i] = append(output[i], m[v].(time.Time).Format("2006-01-02 15:04:05"))
			} else {
				output[i] = append(output[i], m[v])
			}
		}
	}

	data := make(map[string]interface{}, count)
	data["sEcho"] = this.GetString("sEcho")
	data["iTotalRecords"] = counts
	data["iTotalDisplayRecords"] = count
	data["aaData"] = output

	this.Data["json"] = data
	this.ServeJSON()

}
func (this *UserController) GetUsers() {

	this.TplName = "user-list.html"
}

/*
 * aColumns []string `SQL Columns to display`
 * thismodel interface{} `SQL model to use`
 * ctx *context.Context `Beego ctx which contains httpcontext`
 * maps []orm.Params `return result in a interface map as []orm.Params`
 * count int64 `return iTotalDisplayRecords`
 * counts int64 `return iTotalRecords`
 *
 */
func Datatables(aColumns []string, thismodel interface{}, Input *context.BeegoInput, where interface{}) (maps []orm.Params, count int64, counts int64) {
	/*
	 * Paging
	 * 分页请求
	 */
	iDisplayStart, _ := strconv.Atoi(Input.Query("iDisplayStart"))
	iDisplayLength, _ := strconv.Atoi(Input.Query("iDisplayLength"))
	fmt.Println(iDisplayStart)
	fmt.Println(iDisplayLength)
	/*
	 * Ordering
	 * 排序请求
	 */
	querysOrder := []string{}
	if iSortCol_0, _ := strconv.Atoi(Input.Query("iSortCol_0")); iSortCol_0 > -1 {
		ranges, _ := strconv.Atoi(Input.Query("iSortingCols"))
		for i := 0; i < ranges; i++ {
			istring := strconv.Itoa(i)
			if iSortcol := Input.Query("bSortable_" + Input.Query("iSortCol_"+istring)); iSortcol == "true" {
				sordir := Input.Query("sSortDir_" + istring)
				thisSortCol, _ := strconv.Atoi(Input.Query("iSortCol_" + istring))
				if sordir == "asc" {
					querysOrder = append(querysOrder, aColumns[thisSortCol])
				} else {
					querysOrder = append(querysOrder, "-"+aColumns[thisSortCol])
				}
			}
		}
	}
	/*
	 * Filtering
	 * 快速过滤器
	 */
	//querysFilter := []string{}
	cond := orm.NewCondition()
	if len(Input.Query("sSearch")) > 0 {
		for i := 0; i < len(aColumns); i++ {
			cond = cond.Or(aColumns[i]+"__icontains", Input.Query("sSearch"))
		}

	}
	/* Individual column filtering */
	for i := 0; i < len(aColumns); i++ {
		if Input.Query("bSearchable_"+strconv.Itoa(i)) == "true" && len(Input.Query("sSearch_"+strconv.Itoa(i))) > 0 {
			cond = cond.And(aColumns[i]+"__icontains", Input.Query("sSearch"))
		}
	}

	//where条件
	wheres, ok := where.(map[string]interface{})
	if ok {
		for k, v := range wheres {
			fmt.Println(k, v)
			cond = cond.And(k, v)
		}
	}
	fmt.Println(where)
	//用户管理GID
	gid := Input.Query("gid")
	if gid != "" {
		gid2, _ := strconv.Atoi(gid)
		cond = cond.And("gid", int64(gid2))
	}

	//客服管理
	accountid := Input.Query("aid")
	fmt.Println(accountid)
	if accountid != "" {
		aid, _ := strconv.Atoi(accountid)
		cond = cond.And("accountid", int64(aid))
	}
	/*
	 * GetData
	 * 数据请求
	 */
	o := orm.NewOrm()
	qs := o.QueryTable(thismodel)
	fmt.Println(qs)
	counts, _ = qs.Count()
	fmt.Println(counts)
	qs = qs.Limit(iDisplayLength, iDisplayStart)
	qs = qs.SetCond(cond)
	for _, v := range querysOrder {
		qs = qs.OrderBy(v)
	}
	qs.Values(&maps)
	count, _ = qs.Count()
	fmt.Println(maps)
	return maps, count, counts
}
