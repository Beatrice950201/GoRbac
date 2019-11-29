package model

import (
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/models/extend"
	"regexp"
	"strings"
	"time"
)

type AdminUser struct {
	Id          int     `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Pid         int     `orm:"index;column(pid);type(int);default(0);description(上级用户)" json:"pid" form:"pid"`
	Username    string  `orm:"unique;column(username);size(32);type(char);default();description(用户名)"  json:"username" form:"username"`
	Nickname    string  `orm:"column(nickname);size(100);type(char);default();description(用户昵称)"  json:"nickname" form:"nickname"`
	Password    string  `orm:"column(password);size(96);type(char);default();description(密码)"  json:"password" form:"password"`
	Email       string  `orm:"column(email);size(64);type(char);default();description(邮箱)"  json:"email" form:"email"`
	Mobile      string  `orm:"column(mobile);size(11);type(char);default();description(手机号码)"  json:"mobile" form:"mobile"`
	Avatar      string  `orm:"column(avatar);size(200);type(char);default();description(头像)"  json:"avatar" form:"avatar"`
	Role        int     `orm:"index;column(role);type(int);default(0);description(角色组ID)"  json:"role" form:"role"`
	CreateTime  time.Time  `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time  `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
	LastLoginIp int64   `orm:"column(last_login_ip);default(0);description(最后登录IP)" json:"last_login_ip"`
	Sort        int     `orm:"column(sort);size(10);default(0);description(排序序号)" json:"sort" form:"sort"`
	Status      int8    `orm:"index;column(status);size(1);type(int);default(0);description(启用状态)" json:"status" form:"status"`
}

// 设置引擎
func (m *AdminUser) TableEngine() string {
	return "INNODB"
}

// 登陆方法
func Login(username string,password string,IP uint32)(int,string) {
	username = strings.Trim(username," ")
	password = strings.Trim(password," ")
	var user AdminUser
	o := orm.NewOrm()
	userTable := new(AdminUser)
	qs := o.QueryTable(userTable)
	if match,_ := regexp.MatchString(`/^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$/`,username);match == true{
		qs = qs.Filter("email", username)
	}else if match,_ := regexp.MatchString(`/^1\d{10}$/`,username);match == true {
		qs = qs.Filter("mobile", username)
	}else {
		qs = qs.Filter("username", username)
	}
	if err := qs.Filter("status", 1).One(&user);err != nil {
		return 0 ,"未查询到用户或者用户被禁用！"
	}
	if pwdMatch := extend.ComparePasswords(user.Password, password);pwdMatch != true{
		return 0 ,"账号或密码错误！"
	}
    return AutoLogin(user.Id,IP)
}

// 自动登陆程序
func AutoLogin(id int,IP uint32) (uid int, err string)  {
	o := orm.NewOrm()
	user := AdminUser{Id: id}
	if err := o.Read(&user);err != nil{
	   return 0,"用户不存在！"
	}
	if user.Status <= 0{
	   return 0,"用户被禁用！"
	}
	if id != 1{
		if user.Role == 0{
			return 0,"禁止访问，原因：未分配角色！"
		}
		if IsAccess(user.Role) == false{
			return 0,"该角色禁止访问！"
		}
	}
	user.LastLoginIp = int64(IP)
	if _, err := o.Update(&user); err != nil {
		return 0,"更新IP失败，登陆失败！"
	}
	return user.Id,""
}

// 获取一条用户信息
func FindOne(id int) (user AdminUser)  {
	o := orm.NewOrm()
	user = AdminUser{Id: id}
	_ = o.Read(&user)
	return
}

// 查找用户名是否存在
func FindUsernameExtends(username string) bool  {
	var user AdminUser
	o := orm.NewOrm();userTable := new(AdminUser)
	if err := o.QueryTable(userTable).Filter("username", username).One(&user);err == nil && user.Id >0 {
		return true
	}else {
		return false
	}
}
// 手机号码是否已经存在
func FindMobileExtends(mobile string) bool  {
	var user AdminUser
	o := orm.NewOrm();userTable := new(AdminUser)
	if err := o.QueryTable(userTable).Filter("mobile", mobile).One(&user);err == nil && user.Id >0 {
		return true
	}else {
		return false
	}
}
// 邮箱是否已经存在
func FindEmailExtends(email string) bool  {
	var user AdminUser
	o := orm.NewOrm()
	userTable := new(AdminUser)
	if err := o.QueryTable(userTable).Filter("email", email).One(&user);err == nil && user.Id >0 {
		return true
	}else {
		return false
	}
}

// 获取层级结构
func ToLayerUsers(roles []*AdminUser,pid int,level int) (trees []*AdminUser)  {
	for _,v := range roles {
		if pid == v.Pid{
			titlePrefix := strings.Repeat("&nbsp;", level * 8) + "┝ "
			if pid >0{
				v.Nickname = titlePrefix + v.Nickname
			}
			child := ToLayerUsers(roles,v.Id,level+1)
			trees = append(append(trees,v),child...)
		}
	}
	return
}

// 获取当前用户根用户
func RootUserData(uid int) AdminUser {
	user := FindOne(uid)
	if user.Pid == 0{
		return user
	}else {
		return RootUserData(user.Pid)
	}
}