package admin

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/models/model"
	"github.com/beatrice950201/GoRbac/models/verify"
	"strconv"
)

type RoleController struct {
	MainController
}

// 获取授权节点
func (c *RoleController) buildJsTree(roleId int) (authTree []*model.RolesJsTree)  {
	var Tree []*model.RolesJsTree
	Group := model.GroupMenu()
	if jsonData, _ := json.Marshal(Group); len(jsonData) >0 {
		_ = json.Unmarshal(jsonData, &Tree)
	}
	for _, v := range Tree {
		if model.CheckAuth(v.Id,c.User.Id,c.User.Role) == true{
			LayerChild := model.ToLayer(model.MarshalMenu(model.GroupAllMenu()),v.Id)
			v.HtmlTree = model.BuildJsTree(LayerChild, model.Roles(roleId),c.User.Id,c.User.Role)
			authTree = append(authTree,v)
		}
	}
  return
}

// @router /role/index [get]
func (c *RoleController) Index()  {
	var roles []*model.AdminRole
	var lists []*model.RoleExtend
	o := orm.NewOrm()
	userTable := new(model.AdminRole)
	_, _ = o.QueryTable(userTable).All(&roles)
	roles = model.ToLayerRole(roles,0,0)
	if jsonData, _ := json.Marshal(roles); len(jsonData) >0 {
		_ = json.Unmarshal(jsonData, &lists)
	}
	for _, v := range lists {
		if name := model.RoleName(v.Pid);name == ""{
			v.ParentName = "顶级角色"
		}else {
			v.ParentName = name
		}
	}
	c.Data["roles"] = lists
}

// @router /role/status [post]
func (c *RoleController) Status()  {
	 ids := c.GetStrings("ids")
	 if len(ids) <= 0{c.ErrorJson("请选择您要操作的数据！")}
	 status,_ := c.GetInt("status",0)
	var(message string;statusRes bool);o := orm.NewOrm();_ = o.Begin()// 启动事务
     for _, v := range ids {
     	 id,_:= strconv.Atoi(v)
     	 role := model.RoleOne(id)
     	 role.Status = int8(status)
		 if _, err := o.Update(&role);err == nil {
			 message = "更新状态成功";statusRes = true
		 }else {
			 message = "更新状态失败，请稍后再试！！";statusRes = false;break
		 }
	 }
	if statusRes == true{
		_ = o.Commit()
		c.SuccessJson(message,nil,beego.URLFor("RoleController.Index"))
	}else {
		_ = o.Rollback()
		c.ErrorJson(message)
	}

}

// @router /role/access [post]
func (c *RoleController) Access()  {
	ids := c.GetStrings("ids")
	if len(ids) <= 0{c.ErrorJson("请选择您要操作的数据！")}
	status,_ := c.GetInt("status",0)
	o := orm.NewOrm()
	for _, v := range ids {
		id,_:= strconv.Atoi(v)
		role := model.RoleOne(id)
		role.Access = int8(status)
		if _, err := o.Update(&role);err == nil {
			c.SuccessJson("更新状态成功",nil)
		}else {
			c.ErrorJson("更新状态失败，请稍后再试！")
		}
	}
}

// @router /role/delete [post]
func (c *RoleController) Delete()  {
	ids := c.GetStrings("ids")
	if len(ids) <= 0{c.ErrorJson("请选择您要操作的数据！")}
	var(message string;statusRes bool);o := orm.NewOrm();_ = o.Begin()// 启动事务
	for _, v := range ids {
		id,_:= strconv.Atoi(v)
		if id > 1{
			if false == model.IsExtendChildRole(id) {
				if num, err := o.Delete(&model.AdminRole{Id: id}); err == nil && num > 0 {
					message = "删除分组成功";statusRes = true
				}else {
					message = "删除分组失败，请稍后再试！！";statusRes = false;break
				}
			}else {
				message = "请先删除所有子角色组！";statusRes = false;break
			}
		}else {
			message = "此角色组不允许该操作！";statusRes = false;break
		}
	}
	if statusRes == true{
		_ = o.Commit()
		c.SuccessJson(message,nil,beego.URLFor("RoleController.Index"))
	}else {
		_ = o.Rollback()
		c.ErrorJson(message)
	}
}

// @router /role/create [get,post]
func (c *RoleController) Create()  {
	if c.IsAjax(){
		role  := model.AdminRole{}
		if err := c.ParseForm(&role); err != nil {
			c.ErrorJson("操作错误:"+error.Error(err))
		}
		Ok,Message := verify.RoleForm(role)
		if Ok == false{
			c.ErrorJson(Message)
		}else {
			o := orm.NewOrm()
			_, err := o.Insert(&role)
			if err == nil {
				c.SuccessJson("插入分组成功",nil,beego.URLFor("RoleController.Index"))
			}else {
				c.ErrorJson("插入分组失败，请稍后再试！")
			}
		}
	}else {
		c.Data["pids"] = model.ToLayerRole(model.RoleAll(),0,0)
		c.Data["group"] = c.buildJsTree(0)
	}
}

// @router /role/edit [get,post]
func (c *RoleController) Edit()  {
	id := c.GetMustInt("id","非法跳转....")
	if c.IsAjax(){
		role  := model.RoleOne(id)
		if err := c.ParseForm(&role); err != nil {
			c.ErrorJson("操作错误:"+error.Error(err))
		}
		Ok,Message := verify.RoleForm(role)
		if Ok == false{
			c.ErrorJson(Message)
		}else {
			o := orm.NewOrm()
			_, err := o.Update(&role)
			if err == nil {
				c.SuccessJson("更新角色成功",nil,beego.URLFor("RoleController.Index"))
			}else {
				c.ErrorJson("更新角色失败，请稍后再试！")
			}
		}
	}else {
		c.Data["pids"]  = model.ToLayerRole(model.RoleAll(),0,0)
		c.Data["group"] = c.buildJsTree(id)
		c.Data["info"]  = model.RoleOne(id)
	}
}
