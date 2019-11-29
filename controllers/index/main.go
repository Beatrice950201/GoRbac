package index

import (
	"github.com/beatrice950201/eibk/controllers"
)

type MainController struct {
	controllers.BaseController
}

// 准备下一级构造函数
type NextPreparer interface {
	NextPrepare()
}

// 构造函数
func (c *MainController) NestPrepare()  {
	if app, ok := c.AppController.(NextPreparer); ok {
		app.NextPrepare()
	}
	c.Module = "index"
}