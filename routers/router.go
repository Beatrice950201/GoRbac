package routers

import (
	"github.com/astaxie/beego"
	"github.com/beatrice950201/GoRbac/controllers"
	"github.com/beatrice950201/GoRbac/controllers/admin"
	"github.com/beatrice950201/GoRbac/controllers/admin/news"
	"github.com/beatrice950201/GoRbac/controllers/index"
)

func init() {

	beego.ErrorController(&controllers.ErrorController{})

	// 前台
	beego.Include(&index.HomeController{})

	// 后台
	beego.AddNamespace(
		beego.NewNamespace("/admin",beego.NSInclude(
			&admin.IndexController{},
			&admin.SignController{},
			&admin.MenuController{},
			&admin.RoleController{},
			&admin.UserController{},
			&news.NewsDocument{},
			&news.NewsCategory{},
		)))
}
