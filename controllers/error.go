package controllers

import (
	"fmt"
)

type ErrorController struct {
	BaseController
}

// 构造行数
func (c *ErrorController) NestPrepare()  {
	c.Module = "admin"
}

// 404 错误处理
func (c *ErrorController) Error404() {
	if c.IsAjax() {
		c.Ctx.Output.Status = 200
		c.ErrorJson("非法访问")
	} else {
		c.Data["content"] = "非法访问"
	}
}

// 500 处理功能
func (c *ErrorController) Error500() {
	err := c.Data["error"].(error)
	if c.IsAjax() {
		c.Ctx.Output.Status = 200
		c.ErrorJson(error.Error(err))
	} else {
		c.Data["content"] = fmt.Sprintf("错误：%s", error.Error(err))
	}
}

// 数据库检测功能
func (c *ErrorController) ErrorDb() {
	if c.IsAjax() {
		c.Ctx.Output.Status = 200
		c.ErrorJson("database is now down")
	} else {
		c.Data["content"] = fmt.Sprintf("错误：%s", "database is now down")
	}
}
