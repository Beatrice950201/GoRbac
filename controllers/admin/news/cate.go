package news

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/controllers/admin"
	"github.com/beatrice950201/GoRbac/models/extend"
	"github.com/beatrice950201/GoRbac/models/model"
	"github.com/beatrice950201/GoRbac/models/model/news"
	"strconv"
)

type NewsCategory struct {
	admin.MainController
}

// @router /news/cate/index [get]
func (c *NewsCategory) Index()  {
	page, _ := c.GetInt("page",1)
	rootUsers := model.RootUserData(c.User.Id)
	pageSize, _ := beego.AppConfig.Int("web_page_size")
	var cate []*news.TeamNewsCategory
	_, _ = orm.NewOrm().QueryTable(new(news.TeamNewsCategory)).Filter("uid",rootUsers.Id).Limit(pageSize, page*pageSize-pageSize).All(&cate)
	count,_ := orm.NewOrm().QueryTable(new(news.TeamNewsCategory)).Filter("uid",rootUsers.Id).Count()
	c.Data["Page"] = extend.PageUtil(int(count),page,pageSize,cate)
}

// @router /news/cate/create [get,post]
func (c *NewsCategory) Create()  {
	if c.IsAjax(){
		cate  := news.TeamNewsCategory{}
		if err := c.ParseForm(&cate); err != nil {
			c.ErrorJson("操作错误:"+error.Error(err))
		}
		if cate.Title == ""{
			c.ErrorJson("请填写分类标题...")
		}
		rootGroup := model.RootUserData(c.User.Id)
		cate.PushUid = c.User.Id
		cate.Nickname = c.User.Nickname
		cate.Uid = rootGroup.Id  // 所属公司
		_, err := orm.NewOrm().Insert(&cate)
		if err == nil {
			c.SuccessJson("添加栏目成功",nil,beego.URLFor("NewsCategory.Index"))
		}else {
			c.ErrorJson("添加栏目失败~")
		}
	}
}

// @router /news/cate/edit [get,post]
func (c *NewsCategory) Edit() {
   id := c.GetMustInt("id","非法操作！")
   if c.IsAjax(){
	   cate  := news.TeamNewsCategoryOne(id)
	   if err := c.ParseForm(&cate); err != nil {
		   c.ErrorJson("操作错误:"+error.Error(err))
	   }
	   if cate.Title == ""{
	   	   c.ErrorJson("请填写分类标题！")
	   }
	   if _, err := orm.NewOrm().Update(&cate);err == nil {
		   c.SuccessJson("更新分类成功！",nil,beego.URLFor("NewsCategory.Index"))
	   }else {
		   c.ErrorJson("更新分类失败！")
	   }
   } else {
	   c.Data["info"] = news.TeamNewsCategoryOne(id)
   }
}

// @router /news/cate/delete [post]
func (c *NewsCategory) Delete()  {
	ids := c.GetStrings("ids")
	if len(ids) <= 0{c.ErrorJson("请选择您要操作的数据！")}
	var(
		message string
		statusRes bool
	)
	o := orm.NewOrm();_ = o.Begin()// 启动事务
	for _, v := range ids {
		id,_:= strconv.Atoi(v)
		if num, err := o.Delete(&news.TeamNewsCategory{Id: id}); err == nil && num > 0 {
			message = "删除分类成功";statusRes = true
		}else {
			message = "删除分类失败，请稍后再试！";statusRes = false;break
		}
	}
	if statusRes == true{
		_ = o.Commit()
		c.SuccessJson(message,nil,beego.URLFor("NewsCategory.Index"))
	}else {
		_ = o.Rollback()
		c.ErrorJson(message)
	}
}

// @router /news/cate/status [post]
func (c *NewsCategory) Status()  {
	id := c.GetMustInt("ids","非法操作！！")
	status,_ := c.GetInt("status",0)
	info := news.TeamNewsCategoryOne(id)
	info.Status = int8(status)
	if _, err := orm.NewOrm().Update(&info); err == nil {
		c.SuccessJson("更新状态成功;",nil,beego.URLFor("NewsCategory.Index"))
	} else {
		c.ErrorJson("更新状态失败！")
	}
}