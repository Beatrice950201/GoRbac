package cache

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
)

var Bm cache.Cache

func init()  {
	cachepath := beego.AppConfig.String("cachepath")
	filesuffix := beego.AppConfig.String("filesuffix")
	option := fmt.Sprintf(`{"CachePath":"%s","FileSuffix":"%s","DirectoryLevel":"2","EmbedExpiry":"120"}`,cachepath,filesuffix)
	Bm, _ = cache.NewCache("file", option)
}
