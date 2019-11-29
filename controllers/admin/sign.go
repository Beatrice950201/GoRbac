package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/beatrice950201/GoRbac/controllers"
	"github.com/beatrice950201/GoRbac/models/extend"
	"github.com/beatrice950201/GoRbac/models/model"
	"path"
	"strconv"
	"strings"
	"time"
)

type SignController struct {
	controllers.BaseController
}

 // 构造函数
func (c *SignController) NestPrepare()  {
	c.Module = "admin"
}

// @router /sign/login [get,post]
func (c *SignController) Login()  {
	if c.IsAjax(){
		username := c.GetMustString("login-username", "请正确填写用户名...")
		password := c.GetMustString("login-password", "请正确填写密码...")
		rememberString := c.GetString("login-remember", "off")
		remember := false
		if rememberString == "on"{
			remember = true
		}
		IP := extend.Ip2long(c.Ctx.Input.IP())
		if uid,err := model.Login(username,password,IP);uid <=0{
			c.ErrorJson(err)
		}else {
			c.SetSession("uid",uid)
			if remember { // 写入cooke
				c.Ctx.SetCookie("login_remember",strconv.Itoa(uid),86400*7, "/")
			}
			c.Ctx.SetCookie("ThisGroupId",strconv.Itoa(1),86400, "/")
			c.SuccessJson("登陆系统成功！正在跳转...",nil,beego.URLFor("IndexController.Index"))
		}
	}
}

// @router /sign/login_quit [get]
func (c *SignController) LoginQuit()  {
	 c.DelSession("uid")
	 c.Ctx.SetCookie("login_remember","0",100, "/")
	 c.Ctx.Redirect(302,beego.URLFor("IndexController.Index"))
}

// @router /sign/uploader [post]
func (c *SignController) Uploader()  {
	f, h, err := c.GetFile("file")
	if err != nil {
		c.ErrorJson("文件上传失败：" + error.Error(err))
	}
	defer f.Close()
	filePathSave := extend.CreateDateDir("static/upload/")
	fileNameSave := time.Now().Format("20060102150405") + path.Ext(h.Filename)
	if err = c.SaveToFile("file", filePathSave + "/" + fileNameSave);err != nil{
		c.ErrorJson("文件上传失败：" + error.Error(err))
	}else {
		res := strings.Replace(filePathSave + "\\" + fileNameSave, "\\", "/", 3)
		c.SuccessJson("文件上传成功！",res)
	}
}

// @router /sign/uploader_edit [post]
func (c *SignController) UploaderEdit(){
	f, h, err := c.GetFile("upload")
	if err != nil {
		c.ErrorJson("文件上传失败：" + error.Error(err))
	}
	defer f.Close()
	filePathSave := extend.CreateDateDir("static/upload/")
	fileNameSave := time.Now().Format("20060102150405") + path.Ext(h.Filename)
	if err = c.SaveToFile("upload", filePathSave + "/" + fileNameSave);err != nil{
		logs.Error(err)
	}else {
		json := make(map[string]interface{})
		json["uploaded"] = true
		json["url"] = "/" + strings.Replace(filePathSave + "\\" + fileNameSave, "\\", "/", 3)
		c.Data["json"] = json
		c.ServeJSON()
	}
}