package verify

import (
	"github.com/astaxie/beego/validation"
	"github.com/beatrice950201/GoRbac/models/model"
	"regexp"
)

func UserForm(user model.AdminUser,IsEdit bool) (Ok bool,message string) {
	valid := validation.Validation{}
	if res := valid.Match(user.Username,regexp.MustCompile(`^\w{3,}(\.\w+)*@[A-z0-9]+(\.[A-z]{2,5}){1,2}$`),"username").Message("请填写正确的用户名");!res.Ok{
		return false,error.Error(res.Error)
	}
	if e := valid.Required(user.Nickname, "nickname").Message("请填写用户昵称");e.Ok == false{
		return false,error.Error(e.Error)
	}
	if IsEdit == false{
		if e := valid.Required(user.Password, "password").Message("请填写用户密码");e.Ok == false{
			return false,error.Error(e.Error)
		}
	}
	//if e := valid.Email(user.Email, "email").Message("请填写正确的用户邮箱");e.Ok == false{
	//	return false,error.Error(e.Error)
	//}
	//if e := valid.Mobile(user.Mobile, "mobile").Message("请填写正确的用户手机号码");e.Ok == false{
	//	return false,error.Error(e.Error)
	//}
	//if e := valid.Required(user.Avatar, "avatar").Message("请上传用户头像");e.Ok == false{
	//	return false,error.Error(e.Error)
	//}
	if e := valid.Required(user.Role, "role").Message("请选择用户权限组");e.Ok == false{
		return false,error.Error(e.Error)
	}
	if IsEdit == false {
		if model.FindUsernameExtends(user.Username) {
			return false, "该用户名已经存在..."
		}
		//if model.FindMobileExtends(user.Mobile) {
		//	return false, "该手机号已经存在..."
		//}
		//if model.FindEmailExtends(user.Email) {
		//	return false, "该邮箱已经存在..."
		//}
	}
	return true,""
}
