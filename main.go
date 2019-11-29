package main

import (
	"github.com/astaxie/beego"
	_ "github.com/beatrice950201/GoRbac/models"
	_ "github.com/beatrice950201/GoRbac/models/cache"
	_ "github.com/beatrice950201/GoRbac/models/logs"
	_ "github.com/beatrice950201/GoRbac/routers"
)

func main() {
	beego.Run()
}

