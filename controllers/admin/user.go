package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/models/extend"
	"github.com/beatrice950201/GoRbac/models/model"
	"github.com/beatrice950201/GoRbac/models/verify"
	"os"
	"strconv"
)

type UserController struct {
	MainController
}

// @router /user/index [get]
func (c *UserController) Index()  {
	var users []*model.AdminUser
	o := orm.NewOrm()
	userTable := new(model.AdminUser)
	_, _ = o.QueryTable(userTable).All(&users)
	c.Data["lists"] = model.ToLayerUsers(users,0,0)
}

// @router /user/create [get,post]
func (c *UserController) Create() {
	if c.IsAjax() {
		user := model.AdminUser{}
		if err := c.ParseForm(&user); err != nil {
			c.ErrorJson("操作错误:" + error.Error(err))
		}
		Ok, Message := verify.UserForm(user, false)
		if Ok == true {
			user.Password = extend.HashAndSalt(user.Password)
			o := orm.NewOrm()
			_, err := o.Insert(&user)
			if err == nil {
				c.SuccessJson("添加用户成功", nil, beego.URLFor("UserController.Index"))
			} else {
				c.ErrorJson("添加用户失败！")
			}
		} else {
			c.ErrorJson(Message)
		}
	} else {
		c.Data["group"] = model.ToLayerRole(model.RoleAll(), 0, 0)
		c.Data["pid"], _ = c.GetInt("pid", 0)
	}
}

// @router /user/edit [get,post]
func (c *UserController) Edit()  {
    id := c.GetMustInt("id","非法操作！")
    if id == 1{c.Abort500("该用户不允许此操作！")}
   if c.IsAjax(){
	   user  := model.FindOne(id)
	   if err := c.ParseForm(&user); err != nil {
		   c.ErrorJson("操作错误:"+error.Error(err))
	   }
	   Ok,Message := verify.UserForm(user,true)
	   if Ok == true{
	   	   if c.GetString("password","") != "" {
			  user.Password = extend.HashAndSalt(user.Password)
		   }
		   o := orm.NewOrm()
		   _, err := o.Update(&user)
		   if err == nil {
			   c.SuccessJson("更新用户成功",nil,beego.URLFor("UserController.Index"))
		   }else {
			   c.ErrorJson("更新个人资料失败")
		   }
	   }else {
		   c.ErrorJson(Message)
	   }
   }else {
	   c.Data["group"] = model.ToLayerRole(model.RoleAll(),0,0)
	   c.Data["info"] = model.FindOne(id)
   }
}

// @router /user/status [post]
func (c *UserController) Status()  {
	ids := c.GetStrings("ids")
	if len(ids) <= 0{c.ErrorJson("请选择您要操作的数据！")}
	status,_ := c.GetInt("status",0)
	var(
		message string
		statusRes bool
	)
	o := orm.NewOrm();_ = o.Begin()// 启动事务
	for _, v := range ids {
		id,_:= strconv.Atoi(v)
		if id > 1 {
			role := model.FindOne(id)
			role.Status = int8(status)
			if _, err := o.Update(&role); err == nil {
				message = "更新状态成功";statusRes = true
			} else {
				message = "更新个人资料失败！";statusRes = false;break
			}
		}else {
			message = "此用户不允许该操作！";statusRes = false;break
		}
	}
	if statusRes == true{
		_ = o.Commit()
		c.SuccessJson(message,nil,beego.URLFor("UserController.Index"))
	}else {
		_ = o.Rollback()
		c.ErrorJson(message)
	}
}

// @router /user/delete [post]
func (c *UserController) Delete()  {
	ids := c.GetStrings("ids")
	if len(ids) <= 0{c.ErrorJson("请选择您要操作的数据！")}
	var(
		message string
		statusRes bool
	)
	o := orm.NewOrm();_ = o.Begin()// 启动事务
	for _, v := range ids {
		id,_:= strconv.Atoi(v)
		if id > 1{
			if num, err := o.Delete(&model.AdminUser{Id: id}); err == nil && num > 0 {
				if avatar := model.FindOne(id).Avatar; avatar != ""{
					_= os.Remove(avatar)
				}
				message = "删除用户成功";statusRes = true
			}else {
				message = "删除用户失败，请稍后再试！";statusRes = false;break
			}
		}else {
			message = "此用户不允许该操作！";statusRes = false;break
		}
	}
	if statusRes == true{
		_ = o.Commit()
		c.SuccessJson(message,nil,beego.URLFor("UserController.Index"))
	}else {
		_ = o.Rollback()
		c.ErrorJson(message)
	}
}

// @router /user/profile [post,get]
func (c *UserController) Profile()  {
	if c.IsAjax(){
		user  := model.FindOne(c.User.Id)
		if err := c.ParseForm(&user); err != nil {
			c.ErrorJson("操作错误:"+error.Error(err))
		}
		if c.GetString("password","") != "" {
			user.Password = extend.HashAndSalt(user.Password)
		}
		o := orm.NewOrm()
		if _, err := o.Update(&user); err == nil {
			c.Ctx.SetCookie("ThisGroupId",strconv.Itoa(1),86400, "/")
			c.SuccessJson("更新个人资料成功",nil,beego.URLFor("UserController.Profile"))
		}else {
			c.ErrorJson("更新个人资料失败!")
		}
	}else {
		c.Data["info"] = model.FindOne(c.User.Id)
	}
}
