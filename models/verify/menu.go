package verify

import (
	"github.com/astaxie/beego/validation"
	"github.com/beatrice950201/GoRbac/models/model"
)

func MenuForm(menu model.AdminMenu) (Ok bool,message string) {
	valid := validation.Validation{}
	if e := valid.Required(menu.Title, "title").Message("标题不能为空");e.Ok == false{
		return false,error.Error(e.Error)
	}
	//if e := valid.Required(menu.Pid, "pid").Message("请选择上级节点");e.Ok == false{
	//	return false,error.Error(e.Error)
	//}
	if e := valid.Required(menu.UrlType, "url_type").Message("请选择节点类型");e.Ok == false{
		return false,error.Error(e.Error)
	}
	//if e := valid.Required(menu.UrlValue, "url_value").Message("请填写节点地址");e.Ok == false{
	//	return false,error.Error(e.Error)
	//}
	if e := valid.Required(menu.Icon, "icon").Message("请填写节点图标");e.Ok == false{
		return false,error.Error(e.Error)
	}
	return true,""
}