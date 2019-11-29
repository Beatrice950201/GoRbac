package controllers

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/beatrice950201/GoRbac/models/extend"
	"strings"
)

type BaseController struct {
	beego.Controller
	Module string
}

// 子级构造
type NestPreparer interface {
	NestPrepare()
}

// JSON 返回格式
type ResultJsonValue struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Url     string      `json:"url,omitempty"`
	Data   interface{}  `json:"data,omitempty"`
}


// 构造行数
func (c *BaseController) Prepare()  {
	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
	c.SetTemplate()
}

// 模板自动定义
func (c *BaseController) SetTemplate()  {
	 controllerName,actionName := c.GetControllerAndAction()
	 controllerName = strings.Replace(controllerName, "Controller", "", -1)
	 c.TplName = c.Module + "/" + strings.Replace(extend.SnakeString(controllerName), "_", "/", -1) + "/" + strings.ToLower(actionName) + ".html"
}

// 获取控制器AND方法
func (c *BaseController) CandA() (controller ,action string)  {
	controller,action = c.GetControllerAndAction()
	return
}

// 成功返回Json
func (c *BaseController) SuccessJson(message string,data interface{}, options ...string)  {
	var url string
	if len(options) > 0 {url = options[0]}
	c.Data["json"] = &ResultJsonValue{
		Code    :   0,
		Status  :   true,
		Data    :   data,
		Message :   message,
		Url     :   url,
	}
	c.ServeJSON()
	c.StopRun()
}


// 失败返回Json
func (c *BaseController) ErrorJson(message string, options ...string)  {
	var url string
	if len(options) > 0 {url = options[0]}
	c.Data["json"] = &ResultJsonValue{
		Code    :   -1,
		Status  :   false,
		Message :   message,
		Url     :   url,
	}
	c.ServeJSON()
	c.StopRun()
}

// 检测字符串参数返回
func (c *BaseController) GetMustString(key string, message string) string {
	str := c.GetString(key, "")
	if len(str) == 0 {
		c.Abort500(message)
	}
	return str
}
// 检测INT参数返回
func (c *BaseController) GetMustInt(key string, message string) int {
	str,_:= c.GetInt(key, 0)
	if str == 0 {
		c.Abort500(message)
	}
	return str
}

// 错误检测处理
func (c *BaseController) Abort500(err string) {
	c.Data["error"] = errors.New(err)
	c.Abort("500")
}

