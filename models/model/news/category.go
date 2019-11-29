package news

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// 设置表结构
type TeamNewsCategory struct {
	Id          int     `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Uid         int     `orm:"index;column(uid);type(int);default(0);description(创建公司)" json:"uid" form:"uid"`
	Pid         int     `orm:"index;column(pid);type(int);default(0);description(上级ID)" json:"pid" form:"pid"`
	Nickname    string  `orm:"column(nickname);size(100);type(char);default();description(发布用户名称)" json:"nickname" form:"nickname"`
	PushUid     int     `orm:"index;column(push_uid);type(int);default(0);description(发布用户ID)" json:"push_uid" form:"push_uid"`
	Title       string  `orm:"column(title);size(32);type(char);default();description(分类标题)" json:"title" form:"title"`
	Description string  `orm:"column(description);size(255);type(char);default();description(栏目描述)" json:"description" form:"description"`
	CreateTime  time.Time  `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time  `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
	Sort        int     `orm:"column(sort);size(10);default(0);description(排序序号)" json:"sort" form:"sort"`
	Status      int8    `orm:"index;column(status);size(1);type(int);default(0);description(启用状态)" json:"status" form:"status"`
}

// 设置引擎为 INNODB
func (m *TeamNewsCategory) TableEngine() string {
	return "INNODB"
}

// 获取一条数据
func TeamNewsCategoryOne(id int) TeamNewsCategory {
	cate := TeamNewsCategory{Id: id}
	_ = orm.NewOrm().Read(&cate)
	return cate
}

// 获取符合状态列表
func TeamNewsCategoryStatusList(status int) []*TeamNewsCategory {
	var list []*TeamNewsCategory
	_, _ = orm.NewOrm().QueryTable(new(TeamNewsCategory)).Filter("status", status).All(&list)
	return list
}

// 获取公司所属分类列表
func TeamNewsCategoryStatusListTeam(status int,root int) []*TeamNewsCategory {
	var list []*TeamNewsCategory
	_, _ = orm.NewOrm().QueryTable(new(TeamNewsCategory)).Filter("status", status).Filter("uid", root).All(&list)
	return list
}



