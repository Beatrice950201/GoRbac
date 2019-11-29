package cache

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"time"
)

func init()  {
	folderName := time.Now().Format("20060102")
	logsPath := beego.AppConfig.String("logs_path") + folderName + ".log"
	option := fmt.Sprintf(`{"filename":"%s","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`,logsPath)
	_ = logs.SetLogger(logs.AdapterFile,option)
}
