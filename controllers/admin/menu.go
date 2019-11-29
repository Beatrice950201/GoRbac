package admin

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/models/model"
	"github.com/beatrice950201/GoRbac/models/verify"
	"strconv"
)

type MenuController struct {
	MainController
}

// @router /menu/index [get]
func (c *MenuController) Index()  {
	var groupMenuCurrentId int
	var err error
	if groupMenuCurrentId ,err = c.GetInt("id",0); err != nil || groupMenuCurrentId <= 0 {
		groupMenuCurrentId = 1
	}
	root := groupMenuCurrentId
	c.Data["groupMenu"] = model.GroupMenu()
	c.Data["groupMenuCurrentId"] = groupMenuCurrentId
	c.Data["menus"] = model.MenuIndexTreeHtml(model.GroupAllMenu(),groupMenuCurrentId,root)
}

// @router /menu/cookies [post]
func (c *MenuController) Cookies()  {
	id,err := c.GetInt("id", 0)
	if err != nil || id <= 0 {
		c.ErrorJson("非法操作，请稍后再试")
	}
	if init := model.InitUrlFor(id,c.User.Id,c.User.Role); init != ""{
		c.Ctx.SetCookie("ThisGroupId",strconv.Itoa(id),3600 * 12, "/")
		c.SuccessJson("写入Cookie成功",init)
	}else {
		c.ErrorJson("当前分组无子节点")
	}
}

// @router /menu/create [post,get]
func (c *MenuController) Create()  {
	if c.IsAjax(){
		menu  := model.AdminMenu{}
		if err := c.ParseForm(&menu); err != nil {
			c.ErrorJson("操作错误:"+error.Error(err))
		}
        Ok,Message := verify.MenuForm(menu)
        if Ok == false{
			c.ErrorJson(Message)
		}else {
			o := orm.NewOrm()
			_, err := o.Insert(&menu)
			if err == nil {
				c.SuccessJson("插入节点成功",nil,beego.URLFor("MenuController.Index","id",c.ThisGroupId()))
			}else {
				c.ErrorJson("插入节点失败，请稍后再试！")
			}
		}
	}else {
		c.Data["root"] ,_= c.GetInt("root",1)
		c.Data["pid"] ,_= c.GetInt("id",0)
		c.Data["groupMenu"] = model.GroupMenu()
	}
}

// @router /menu/edit [post,get]
func (c *MenuController) Edit()  {
	if c.IsAjax(){
		id := c.GetMustInt("id","非法操作，请稍后再试...")
		menu := model.MenuOne(id)
		if err := c.ParseForm(&menu); err != nil {
			c.ErrorJson("操作错误:"+error.Error(err))
		}
		Ok,Message := verify.MenuForm(menu)
		if Ok == false{
			c.ErrorJson(Message)
		}else {
			o := orm.NewOrm()
			_, err := o.Update(&menu)
			if err == nil {
				c.SuccessJson("插入更新成功",nil,beego.URLFor("MenuController.Index","id",c.ThisGroupId()))
			}else {
				c.ErrorJson("更新节点失败，请稍后再试！")
			}
		}
	}else {
		if id,err := c.GetInt("id",0); id > 0{
			c.Data["info"]   = model.MenuOne(id)
			c.Data["root"] ,_= c.GetInt("root",1)
			c.Data["groupMenu"] = model.GroupMenu()
		}else {
			panic(err) //todo 跳转错误页面
		}
	}
}

// @router /menu/status [post]
func (c *MenuController) Status()  {
	root,_:= c.GetInt("root",1)
	id := c.GetMustInt("id","非法操作，请稍后再试...")
	menu := model.MenuOne(id)
	menu.Status,_ = c.GetInt8("status",0)
	o := orm.NewOrm()
	_, err := o.Update(&menu)
	if err == nil {
		c.SuccessJson("状态更新成功",nil,beego.URLFor("MenuController.Index","id",root))
	}else {
		c.ErrorJson("状态更新失败，请稍后再试！")
	}
}

// @router /menu/delete [post]
func (c *MenuController) Delete()  {
	root,_:= c.GetInt("root",1)
	id := c.GetMustInt("id","非法操作，请稍后再试...")
	if false == model.IsExtendChildMenu(id) {
		o := orm.NewOrm()
		if num, err := o.Delete(&model.AdminMenu{Id: id}); err == nil && num > 0 {
			c.SuccessJson("删除节点成功",nil,beego.URLFor("MenuController.Index","id",root))
		}else {
			c.ErrorJson("删除节点失败，请稍后再试！")
		}
	}else {
		c.ErrorJson("请先删除所有子节点！")
	}
}

// @router /menu/save [post]
func (c *MenuController) Save()  {
	 var tree []*model.SaveMenu
     menus := c.GetString("menu")
     root,_  := c.GetInt("root",0)
	 _ = json.Unmarshal([]byte(menus), &tree)
	 result := model.ParseMenuSave(tree,root)
     for _, v := range result {
		 o := orm.NewOrm();_, _ = o.Update(&v)
	 }
	c.SuccessJson("节点更新成功",nil,beego.URLFor("MenuController.Index","id",root))
}

// @router /menu/child_tree [post]
func (c *MenuController) ChildTree()  {
	id := c.GetMustInt("id","非法请求,请稍后再试...")
	AllList := model.MarshalMenu(model.GroupAllMenu())
	list := model.ToLayer(AllList,id)
    c.SuccessJson("",list)
}