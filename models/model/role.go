package model

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/models/extend"
	"strconv"
	"strings"
	"time"
)

type AdminRole struct {
	Id          int     `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Pid         int     `orm:"index;column(pid);type(int);default(0);description(上级ID)"  json:"pid" form:"pid"`
	Name        string  `orm:"column(name);size(32);type(char);default();description(分组名称)"  json:"name" form:"name"`
	Description string  `orm:"column(description);size(200);type(char);default();description(分组描述)"  json:"description" form:"description"`
	MenuAuth    string  `orm:"column(menu_auth);type(text);default();description(节点字符串)"  json:"menu_auth" form:"menu_auth"`
	CreateTime  time.Time  `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time  `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
	Sort        int     `orm:"column(sort);size(10);default(0);description(排序序号)" json:"sort" form:"sort"`
	Status      int8    `orm:"index;column(status);size(1);type(int);default(0);description(启用状态)" json:"status" form:"status"`
	Access      int8    `orm:"index;column(access);type(tinyint);default(0);description(是否可以登陆后台)" json:"access" form:"access"`
}

type RoleExtend struct {
	AdminRole
	ParentName string `json:"parent_name"`
}

type RolesJsTree struct {
	AdminMenu
   HtmlTree string `json:"html_tree"`
}

// 设置引擎为 INNODB
func (m *AdminRole) TableEngine() string {
	return "INNODB"
}

// 获取是否可以登陆后台
func IsAccess(id int)bool  {
	o := orm.NewOrm()
	role := AdminRole{Id: id}
	if err := o.Read(&role);err != nil{
		return false
	}
	if role.Status == 0 || role.Access == 0 {
		return false
	}
	return true
}

// 获取节点组
func Roles(id int) (ids []int) {
	o := orm.NewOrm()
	role := AdminRole{Id: id}
	if err := o.Read(&role);err == nil{
		_ = json.Unmarshal([]byte(role.MenuAuth), &ids)
	}
	return
}

//获取节点名称
func RoleName(id int) (name string)  {
	o := orm.NewOrm()
	role := AdminRole{Id: id}
	if err := o.Read(&role);err == nil{
		name = role.Name
	}
	return
}

// 获取一条数据
func RoleOne(id int) (role AdminRole)  {
	o := orm.NewOrm()
	role = AdminRole{Id: id}
	_ = o.Read(&role)
	return
}

// 获取所有角色组
func RoleAll() (roles []*AdminRole) {
	o := orm.NewOrm()
	userTable := new(AdminRole)
	_, _ = o.QueryTable(userTable).All(&roles)
	return
}

// 检测是否存在子节点
func IsExtendChildRole(id int) bool {
	o := orm.NewOrm()
	userTable := new(AdminRole)
	if num, err := o.QueryTable(userTable).Filter("pid", id).Count();err == nil && num >0{
		return  true
	}else {
		return false
	}
}

// 获取层级结构
func ToLayerRole(roles []*AdminRole,pid int,level int) (trees []*AdminRole)  {
	for _,v := range roles {
		if pid == v.Pid{
			titlePrefix := strings.Repeat("&nbsp;", level * 8) + "┝ "
			if pid >0{
				v.Name = titlePrefix + v.Name
			}
           child := ToLayerRole(roles,v.Id,level+1)
           trees = append(append(trees,v),child...)
		}
	}
	return
}



// 组装节点Html
func BuildJsTree(menus []*LayerMenu,menuAuth []int,uid ,role int) (HtmlTree string) {
     if len(menus) >0{
		 MapJson := make(map[string]interface{})
		 MapJson["opened"] = true
		 MapJson["selected"] = false
		 MapJson["icon"] = ""
		 for _,v := range menus {
			 if CheckAuth(v.Id,uid,role) == true {
				 MapJson["icon"] = "fa fa-fw fa-" + v.Icon
				 MapJson["selected"] = extend.InArray(v.Id, menuAuth)
				 if optionJson, _ := json.Marshal(MapJson); len(optionJson) > 0 {
					 optionString := string(optionJson)
					 optionUrlValue := v.UrlValue
					 if optionUrlValue != "" {
						 optionUrlValue = " ( " + beego.URLFor(optionUrlValue) + " ) "
					 }
					 titleCore := v.Title + optionUrlValue
					 if len(v.Child) > 0 {
						 titleCore += BuildJsTree(v.Child, menuAuth, uid, role)
					 }
					 HtmlTree += "<li id='" + strconv.Itoa(v.Id) + "' data-jstree='" + optionString + "'>" + titleCore + "</li>"
				 }
			 }
		 }
	 }
	return "<ul>" + HtmlTree + "</ul>"
}