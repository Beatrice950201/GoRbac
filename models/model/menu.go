package model

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/models/extend"
	"strconv"
	"time"
)

// 设置表结构
type AdminMenu struct {
	Id          int     `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Pid         int     `orm:"index;column(pid);type(int);default(0);description(上级ID)"  json:"pid" form:"pid"`
	Title       string  `orm:"column(title);size(32);type(char);default();description(节点标题)" json:"title" form:"title"`
	Icon        string  `orm:"column(icon);size(64);type(char);default();description(图标)" json:"icon" form:"icon"`
	UrlType     string  `orm:"column(url_type);size(16);type(char);default();description(连接类型)" json:"url_type" form:"url_type"`
	UrlValue    string  `orm:"column(url_value);size(255);type(char);default();description(连接值)" json:"url_value" form:"url_value"`
	IsMenu      int8    `orm:"column(is_menu);size(1);type(int);default(1);description(是否左侧菜单)" json:"is_menu" form:"is_menu"`
	OnlineHide  int8    `orm:"index;column(online_hide);type(tinyint);default(0);description(生产环境是否隐藏)" json:"online_hide" form:"online_hide"`
	CreateTime  time.Time  `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time  `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
	Sort        int     `orm:"column(sort);size(10);default(0);description(排序序号)" json:"sort" form:"sort"`
	Status      int8    `orm:"index;column(status);size(1);type(int);default(0);description(启用状态)" json:"status" form:"status"`
	Params      string  `orm:"column(params);size(255);type(char);default();description(附加参数)" json:"params" form:"params"`
}

type LayerMenu struct {
	AdminMenu
	Child []*LayerMenu `json:"child" form:"children"`
}

type SaveMenu struct {
	Id int `json:"id"`
	Children []*SaveMenu `json:"children"`
}

// 设置引擎为 INNODB
func (m *AdminMenu) TableEngine() string {
	return "INNODB"
}

// 获取顶部菜单
func HeaderBar() (menus []*AdminMenu,lens int64,e error)  {
	debug, _ := beego.AppConfig.Bool("debug")
	o := orm.NewOrm()
	userTable := new(AdminMenu)
	qs := o.QueryTable(userTable).Filter("status", 1).Filter("pid", 0)
	if debug == false{
		//qs = qs.Filter("online_hide", 1)
	}
	lens,e = qs.All(&menus)
	return
}

//获取初始化跳转地址
func InitUrlFor(id ,uid ,role int) string {
	var menus []*AdminMenu
	o := orm.NewOrm()
	userTable := new(AdminMenu)
	_, _ = o.QueryTable(userTable).Filter("status", 1).Filter("pid", id).OrderBy("-sort").All(&menus)
	for _, v := range menus {
		if CheckAuth(v.Id,uid,role) == true{
			if v.UrlValue != "" {
				return v.UrlValue
			}else {
				var menu []*AdminMenu
				_, _ = o.QueryTable(userTable).Filter("status", 1).Filter("url_value__isnull", false).Filter("pid", v.Id).OrderBy("-sort").All(&menu)
				for _, vs := range menu {
					if CheckAuth(vs.Id,uid,role) == true{
						return vs.UrlValue
					}
				}
			}
		}
	}
	return ""
}

// 获取左边栏目
func SideBar() (menus []*AdminMenu,lens int64,e error) {
	debug, _ := beego.AppConfig.Bool("debug")
	o := orm.NewOrm()
	userTable := new(AdminMenu)
	qs := o.QueryTable(userTable).Filter("status", 1).Filter("is_menu", 1)
	if debug == false{
		//qs = qs.Filter("online_hide", 1)
	}
	lens,e = qs.OrderBy("-sort").All(&menus)
	return
}

// 节点是否授权
func CheckAuth(id ,uid ,role int) bool {
	status := false
	if uid == 1 { // 最高管理员
		status = true
	}else {
		roleIds := Roles(role)
		if extend.InArray(id,roleIds){
			status =  true
		}
	}
	return status
}

//// 生产树结构
func ToLayer(list []*LayerMenu,pid int) (tree []* LayerMenu) {
	for _, v := range list {
		if v.Pid == pid{
			if child := ToLayer(list,v.Id); len(child) > 0{
				v.Child = child
			}
			 tree = append(tree,v)
		 }
	}
	return tree
}

// 转换结构体
func MarshalMenu(menus []*AdminMenu)  (tree []* LayerMenu) {
	if jsonData, _ := json.Marshal(menus); len(jsonData) >0 {
		_ = json.Unmarshal(jsonData, &tree)
	}
	return
}

// 获取当前节点ID；用于定位左侧菜单高亮
func CurrentId(controller,action string,IsSideBar bool) (menu AdminMenu) {
	where := controller + "." + action
	o := orm.NewOrm()
	userTable := new(AdminMenu)
	qs := o.QueryTable(userTable)
	_ = qs.Filter("url_value", where).Filter("pid__gt", 0).One(&menu)
	if IsSideBar == true && menu.IsMenu == 0{
		_ = qs.Filter("id", menu.Pid).One(&menu)
	}
	return
}

// 获取一条数据
func MenuOne(id int) (menu AdminMenu)  {
	o := orm.NewOrm()
	menu = AdminMenu{Id: id}
	_ = o.Read(&menu)
  return
}

// 获取小地图导航
func BreadcrumbMenu(menu AdminMenu,uid ,role int) []AdminMenu {
   var menus []AdminMenu
   if menu.Pid >0{
   	 intOne := MenuOne(menu.Pid)
   	 child := BreadcrumbMenu(intOne,uid ,role)
   	 if len(child) > 0  {
		 for _,v := range child {
		 	if  v.UrlValue == "" {
				if init := InitUrlFor(v.Id,uid ,role); init != ""{
					v.UrlValue = beego.URLFor(init)
				}
			}
			menus = append(menus, v)
		 }
	 }
   }
   if menu.UrlValue != ""{
	   menu.UrlValue = beego.URLFor(menu.UrlValue)
   }
	menus = append(menus, menu)
   return menus
}

// 获取所有顶部模块
func GroupMenu() (list []*AdminMenu) {
	o := orm.NewOrm()
	userTable := new(AdminMenu)
	_, _ = o.QueryTable(userTable).Filter("pid", 0).OrderBy("-sort").All(&list)
	return
}

// 获取所有节点
func GroupAllMenu() (list []*AdminMenu) {
	o := orm.NewOrm()
	userTable := new(AdminMenu)
	_, _ = o.QueryTable(userTable).OrderBy("-sort").All(&list)
	return
}

//是否存在子节点
func IsExtendChildMenu(id int) bool  {
	o := orm.NewOrm()
	userTable := new(AdminMenu)
	if num, err := o.QueryTable(userTable).Filter("pid", id).OrderBy("-sort").Count();err == nil && num >0{
		return  true
	}else {
		return false
	}
}

// 解析成可以写入数据库的格式
func ParseMenuSave(menus []*SaveMenu,pid int) []AdminMenu  {
	var result []AdminMenu
	sort := 1
	for _, v := range menus {
		ResultOne := MenuOne(v.Id)
		ResultOne.Id = v.Id
		ResultOne.Pid = pid
		ResultOne.Sort = sort
		result = append(result,ResultOne)
		if len(v.Children) > 0{
			child := ParseMenuSave(v.Children,v.Id)
			result = append(result,child...)
		}
		sort ++
	}
	return result
}

// 获取树结构HTML渲染列表页
func MenuIndexTreeHtml(list []*AdminMenu,pid int,root int) string {
    var html string
	for _,v := range list {
       if pid == v.Pid{
		   disable := "";if v.Status == 0{disable = "dd-disable"}
		   html += "<li class='dd-item dd3-item "+ disable +"' data-id='"+ strconv.Itoa(v.Id) +"'>"
		   html += "<div class='dd-handle dd3-handle'>拖拽</div><div class='dd3-content'><i class='fa fa-"+ v.Icon +"'></i> " + v.Title
		   if v.UrlValue != ""{
			   html += "<span class='link'><i class='fa fa-link'></i> "+ beego.URLFor(v.UrlValue) +"</span>"
		   }
		   html += "<div class='action'>"
		   html += "<a href='"+beego.URLFor("MenuController.Create","id",v.Id,"root",root)+"' data-toggle='tooltip' data-original-title='新增子节点'><i class='list-icon fa fa-plus fa-fw'></i></a>"
		   html += "<a href='"+beego.URLFor("MenuController.Edit","id",v.Id,"root",root)+"' data-toggle='tooltip' data-original-title='编辑'><i class='list-icon fa fa-pencil-alt fa-fw'></i></a>"
		   if v.Status <= 0{
			   // 启用
			   html += "<a href='javascript:void(0);' data-action='"+beego.URLFor("MenuController.Status","root",root)+"' data-ids='"+ strconv.Itoa(v.Id) +"' class='enable' data-toggle='tooltip' data-original-title='启用'><i class='list-icon fa fa-check-circle fa-fw'></i></a>"
		   }else {
			   // 禁用
			   html += "<a href='javascript:void(0);' data-action='"+beego.URLFor("MenuController.Status","root",root)+"' data-ids='"+ strconv.Itoa(v.Id) +"' class='disable' data-toggle='tooltip' data-original-title='禁用'><i class='list-icon fa fa-ban fa-fw'></i></a>"
		   }
		   html += "<a href='javascript:void(0);' data-action='"+beego.URLFor("MenuController.Delete","root",root)+"' data-ids='"+ strconv.Itoa(v.Id) +"' data-toggle='tooltip' data-original-title='删除' class='delete'><i class='list-icon fa fa-times fa-fw'></i></a></div>"
		   html += "</div>"
		   // 下级节点
		   if ChlidHtml := MenuIndexTreeHtml(list,v.Id,root); ChlidHtml != ""{
			   html += "<ol class='dd-list'>"+ChlidHtml+"</ol>"
		   }
		   html += "</li>"
	   }
	}
	return html
}