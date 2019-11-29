package verify

import (
	"github.com/astaxie/beego/validation"
	"github.com/beatrice950201/GoRbac/models/model"
)

func RoleForm(menu model.AdminRole) (Ok bool,message string) {
	valid := validation.Validation{}
	if e := valid.Required(menu.Name, "name").Message("角色名称不能为空");e.Ok == false{
		return false,error.Error(e.Error)
	}
	if e := valid.Required(menu.Description, "description").Message("角色描述不能为空");e.Ok == false{
		return false,error.Error(e.Error)
	}
	if e := valid.Required(menu.MenuAuth, "menu_auth").Message("请至少选择一个授权节点！");e.Ok == false{
		return false,error.Error(e.Error)
	}
	return true,""
}
