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

type NewsDocument struct {
	admin.MainController
}

// @router /news/document/index [get]
func (c *NewsDocument) Index()  {
	rootUsers := model.RootUserData(c.User.Id)
	page, _ := c.GetInt("page",1)
	cid, _ := c.GetInt("cid",0)
	pageSize, _ := beego.AppConfig.Int("web_page_size")
	var docs []*news.TeamNewsDocument
	qm := orm.NewOrm().QueryTable(new(news.TeamNewsDocument)).Filter("root_uid",rootUsers.Id)
	if cid > 0{
		qm = qm.Filter("cid",cid)
	}
	_, _ = qm.Limit(pageSize, page*pageSize-pageSize).All(&docs)
	count,_ := qm.Count()
	c.Data["Page"] = extend.PageUtil(int(count),page,pageSize,news.TeamNewsDocumentCategoryDispose(docs))
	c.Data["cate"] = news.TeamNewsCategoryStatusListTeam(1,rootUsers.Id)
	c.Data["cid"] = cid
}

// @router /news/document/create [get,post]
func (c *NewsDocument) Create()  {
	rootUsers := model.RootUserData(c.User.Id)
	if c.IsAjax(){
		doc  := news.TeamNewsDocument{}
		if err := c.ParseForm(&doc); err != nil {
			c.ErrorJson("操作错误:"+error.Error(err))
		}
		if doc.Title == ""{ c.ErrorJson("请填写内容标题...") }
		if doc.Description == ""{ c.ErrorJson("请填写内容描述...") }
		if doc.Covers == ""{ c.ErrorJson("请上传封面图...") }
		if doc.Content == ""{ c.ErrorJson("请填写内容...") }
		doc.Uid = c.User.Id
		doc.Nickname = c.User.Nickname
		doc.RootUid = rootUsers.Id
		_, err := orm.NewOrm().Insert(&doc)
		if err == nil {
			c.SuccessJson("添加文章成功",nil,beego.URLFor("NewsDocument.Index"))
		}else {
			c.ErrorJson("添加文章失败~")
		}
	}
	c.Data["cate"] = news.TeamNewsCategoryStatusListTeam(1,rootUsers.Id)
}

// @router /news/document/edit [get,post]
func (c *NewsDocument) Edit()  {
	rootUsers := model.RootUserData(c.User.Id)
	id := c.GetMustInt("id","参数非法！！")
	if c.IsAjax(){
		docs  := news.TeamNewsDocumentOne(id)
		if err := c.ParseForm(&docs); err != nil {
			c.ErrorJson("操作错误:"+error.Error(err))
		}
		if docs.Title == ""{ c.ErrorJson("请填写内容标题...") }
		if docs.Description == ""{ c.ErrorJson("请填写内容描述...") }
		if docs.Covers == ""{ c.ErrorJson("请上传封面图...") }
		if docs.Content == ""{ c.ErrorJson("请填写内容...") }
		if _, err := orm.NewOrm().Update(&docs);err == nil {
			c.SuccessJson("更新文章成功！",nil,beego.URLFor("NewsDocument.Index"))
		}else {
			c.ErrorJson("更新文章失败！")
		}
	}
	c.Data["info"] = news.TeamNewsDocumentOne(id)
	c.Data["cate"] = news.TeamNewsCategoryStatusListTeam(1,rootUsers.Id)
}

// @router /news/document/delete [post]
func (c *NewsDocument) Delete() {
	ids := c.GetStrings("ids")
	if len(ids) <= 0{c.ErrorJson("请选择您要操作的数据！")}
	var(
		message string
		statusRes bool
	)
	o := orm.NewOrm();_ = o.Begin()// 启动事务
	for _, v := range ids {
		id,_:= strconv.Atoi(v)
		if num, err := o.Delete(&news.TeamNewsDocument{Id: id}); err == nil && num > 0 {
			message = "删除文章成功";statusRes = true
		}else {
			message = "删除文章失败，请稍后再试！";statusRes = false;break
		}
	}
	if statusRes == true{
		_ = o.Commit()
		c.SuccessJson(message,nil,beego.URLFor("NewsDocument.Index"))
	}else {
		_ = o.Rollback()
		c.ErrorJson(message)
	}
}

// @router /news/document/status [post]
func (c *NewsDocument) Status() {
	id := c.GetMustInt("ids","非法操作！！")
	status,_ := c.GetInt("status",0)
	info := news.TeamNewsDocumentOne(id)
	info.Status = int8(status)
	if _, err := orm.NewOrm().Update(&info); err == nil {
		c.SuccessJson("更新状态成功;",nil,beego.URLFor("NewsDocument.Index"))
	} else {
		c.ErrorJson("更新状态失败！")
	}
}

// @router /news/document/gather [post]
func (c *NewsDocument) Gather() {
	urls := c.GetMustString("url","请填写链接地址....")
    str := news.CatherHttpLibString(urls)
    if len(str) > 1{
		c.SuccessJson("采集成功！",str)
	}else {
		c.ErrorJson("数据采集不完整，采集失败...")
	}
}

