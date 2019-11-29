package admin

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/GoRbac/controllers"
	"github.com/beatrice950201/GoRbac/models/cache"
	"github.com/beatrice950201/GoRbac/models/extend"
	"github.com/beatrice950201/GoRbac/models/model"
	"strconv"
	"strings"
)

type MainController struct {
	controllers.BaseController
	User       model.AdminUser
	Controller string
	Action     string
}

// 准备下一级构造函数
type NextPreparer interface {
	NextPrepare()
}

// 构造函数
func (c *MainController) NestPrepare()  {
    c.IsLogin()
	c.Controller,c.Action = c.CandA()
	c.checkRoleMenu()
	if app, ok := c.AppController.(NextPreparer); ok {
		app.NextPrepare()
	}
	c.Module = "admin"
	c.SetLayout()
	c.Data["headerMenus"]     = c.TopMenu()
	c.Data["MenuCurrentId"]   = model.CurrentId(c.Controller,c.Action,true)
	c.Data["sideBarMenus"]    = c.SideBarMenus(c.ThisGroupId())
	c.Data["ThisGroupId"]     = c.ThisGroupId()
	c.Data["BreadcrumbMenus"] = c.Breadcrumb()
	c.Data["users"] = c.User
}

// 检测节点是否允许访问
func (c *MainController) checkRoleMenu()  {
	menu := model.CurrentId(c.Controller,c.Action,false)
	if menu.Id <= 0{
		c.Abort500("未检测到当前节点数据！")
	}
	if menu.Status == 0{
		c.Abort500("当前节点已被禁用！")
	}
	isBool := model.CheckAuth(menu.Id,c.User.Id,c.User.Role)
	if isBool == false{
		c.Abort500("您的权限无法访问当前节点！")
	}
}
// 获取小地图导航
func (c *MainController) Breadcrumb() (menus [] model.AdminMenu) {
	 menu := model.CurrentId(c.Controller,c.Action,false)
	 menus = model.BreadcrumbMenu(menu,c.User.Id,c.User.Role)
	return
}

// 获取CookeID
func (c *MainController) ThisGroupId() (ThisGroupId int)  {
	var err error
	if ThisGroupId,err = strconv.Atoi(c.Ctx.GetCookie("ThisGroupId")); err != nil || ThisGroupId <= 0{
		ThisGroupId = 1
	}
	return
}

// 基础布局
func  (c *MainController) SetLayout()  {
	c.Layout = "admin/base.html"
	headFile,scriptFile := c.files()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHead"] = headFile
	c.LayoutSections["Scripts"] = scriptFile
}

// 获得扩展模板
func (c *MainController) files() (headFile,scriptFile string)  {
	cname,aname := c.GetControllerAndAction()
	cname = strings.Replace(cname, "Controller", "", -1)
	tpl := strings.Replace(extend.SnakeString(cname), "_", "/", -1) + "/" + strings.ToLower(aname) + ".html"
	headFile = c.Module + "/common/styles/" + tpl
	scriptFile = c.Module + "/common/script/" + tpl
	if isHeadFile,_ := extend.PathExists("./views/" + headFile);!isHeadFile{
		headFile = ""
	}
	if isScriptFile,_ := extend.PathExists("./views/" + scriptFile);!isScriptFile{
		scriptFile = ""
	}
	return headFile,scriptFile
}

// 是否登陆处理
func (c *MainController) IsLogin()  {
   uid := c.GetSession("uid")
   if uid == nil || uid.(int) <= 0 {
     if c.AutoLogin(){
     	c.Ctx.Redirect(302,beego.URLFor("IndexController.Index"))
	 }else {
		 c.Ctx.Redirect(302,beego.URLFor("SignController.Login"))
	 }
   }else {
	   c.User = model.FindOne(uid.(int))
   }
}

// 是否满足自动登陆
func (c *MainController) AutoLogin() bool {
   uid := c.Ctx.GetCookie("login_remember")
   status := false
   if intUID, err := strconv.Atoi(uid); err == nil{
	   IP := extend.Ip2long(c.Ctx.Input.IP())
	   if intUID > 0 {
		   if AutoStatus,_ := model.AutoLogin(intUID,IP); AutoStatus >0 {
			   c.SetSession("uid",AutoStatus)
			   status = true
		   }
	   }
   }
    return  status
}

// 获取顶部菜单
func (c *MainController) TopMenu() (menus []*model.AdminMenu){
	debug, _ := beego.AppConfig.Bool("debug")
	cacheTag  := "role_header_" + strconv.Itoa(c.User.Role)
	menusCache  := cache.Bm.Get(cacheTag)
	if menusCache == "" {
		if topMenus,_,err := model.HeaderBar();err == nil{
			for _,v := range topMenus {
				if model.CheckAuth(v.Id,c.User.Id,c.User.Role) == true {
					if init := model.InitUrlFor(v.Id,c.User.Id,c.User.Role); init != ""{
						v.UrlValue = beego.URLFor(init)
					}
					menus = append(menus,v)
				}
			}
			if debug == false{
				_ = extend.SetCache(cacheTag,menus)
			}
		}
	}else {
		menus = menusCache.([]*model.AdminMenu)
	}
	return
}

// 获取左侧菜单
func (c *MainController) SideBarMenus(ThisGroupId int) (menus []*model.LayerMenu) {
	debug, _ := beego.AppConfig.Bool("debug")
	cacheTag  := "role_left_" + strconv.Itoa(c.User.Role) + "_" + strconv.Itoa(ThisGroupId)
	menusCache  := cache.Bm.Get(cacheTag)
	if menusCache == "" {
		if SqlMenus,_,err := model.SideBar();err == nil{
			var m []*model.AdminMenu
			for _,v := range SqlMenus {
				if model.CheckAuth(v.Id,c.User.Id,c.User.Role) == true {
					if v.UrlType == "module_admin" && v.UrlValue != ""{
						v.UrlValue = beego.URLFor(v.UrlValue) //todo 处理URL参数问题
					}
					m = append(m,v)
				}
			}
			menus = model.ToLayer(model.MarshalMenu(m),ThisGroupId)
			if debug == false{
				_ = extend.SetCache(cacheTag,menus)
			}
		}
	}else {
		menus = menusCache.([]*model.LayerMenu)
	}
	return
}
