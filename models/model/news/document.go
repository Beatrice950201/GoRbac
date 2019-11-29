package news

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/GoRbac/models/extend"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type TeamNewsDocument struct {
	Id          int     `orm:"pk;auto;column(id);type(int);default(0);description(主键,自增)" json:"id" form:"id"`
	Uid         int     `orm:"index;column(uid);type(int);default(0);description(创建公司)" json:"uid" form:"uid"`
	RootUid     int     `orm:"index;column(root_uid);type(int);default(0);description(所属公司)" json:"root_uid" form:"root_uid"`
	Cid         int     `orm:"index;column(cid);type(int);default(0);description(所属分类)" json:"cid" form:"cid"`
	Nickname    string  `orm:"column(nickname);size(100);type(char);default();description(发布用户名称)" json:"nickname" form:"nickname"`
	Title       string  `orm:"column(title);size(32);type(char);default();description(分类标题)" json:"title" form:"title"`
	Description string  `orm:"column(description);size(255);type(char);default();description(栏目描述)" json:"description" form:"description"`
	IsVideo     int8    `orm:"column(is_video);size(1);type(int);default(0);description(是否存有视频)" json:"is_video" form:"is_video"`
	Covers      string  `orm:"column(covers);type(text);default();description(封面图片)" json:"covers" form:"covers"`
	Content     string  `orm:"column(content);type(text);default();description(内容html)" json:"content" form:"content"`
	View        int     `orm:"column(view);default(0);description(预览数)" json:"view" form:"view"`
	Goods       int     `orm:"column(goods);default(0);description(点赞数)" json:"goods" form:"goods"`
	Quote       int     `orm:"column(quote);default(0);description(转载)" json:"quote" form:"quote"`
	CreateTime  time.Time  `orm:"auto_now_add;type(datetime);description(创建时间)" json:"create_time"`
	UpdateTime  time.Time  `orm:"auto_now;type(datetime);description(更新时间)" json:"update_time"`
	Sort        int     `orm:"column(sort);size(10);default(0);description(排序序号)" json:"sort" form:"sort"`
	Status      int8    `orm:"index;column(status);size(1);type(int);default(0);description(启用状态)" json:"status" form:"status"`
}

//分类结构体
type TeamNewsDocumentCategory struct {
	TeamNewsDocument
	CategoryName string `json:"category_name"`
}

// 设置引擎为 INNODB
func (m *TeamNewsDocument) TableEngine() string {
	return "INNODB"
}

// 结构体转换 加入 分类名称
func TeamNewsDocumentCategoryDispose(list []*TeamNewsDocument) (res []*TeamNewsDocumentCategory) {
	if j, _ := json.Marshal(list); len(j) >0 {
		_ = json.Unmarshal(j, &res)
	}
	for _,v := range res {
		cate := TeamNewsCategoryOne(v.Cid)
		v.CategoryName = cate.Title
	}
	return res
}

// 获取一条数据
func TeamNewsDocumentOne(id int) TeamNewsDocument {
	find := TeamNewsDocument{Id: id}
	_ = orm.NewOrm().Read(&find)
	return find
}

// 获取公众号文章
func CatherHttpLibString(sUrl string) map[string]string {
	res := make(map[string]string)
	rsp := httplib.Get(sUrl)
	if sHtml,_ := rsp.String();sHtml != ""{
		dom,err:=goquery.NewDocumentFromReader(strings.NewReader(sHtml))
		if err == nil{
			dom.Find("#activity-name").Each(func(i int, selection *goquery.Selection) {
				res["title"] = strings.ReplaceAll(strings.ReplaceAll(selection.Text(), " ", ""), "\n", "")
			})
			dom.Find("#js_content").Each(func(i int, selection *goquery.Selection) {
				// 处理图片
				selection.Find("img[data-src]").Each(func(i int, img *goquery.Selection) {
					src,_ := img.Attr("data-src")
					if "https://mmbiz.qpic.cn" == GetQueueUrlDomain(src){
						src = "http://img01.store.sogou.com/net/a/04/link?appid=100520029&url=" + src
					}
					if lopath,e := DownloadResourceLocal(src);e == nil{
						if i > 0 && i < selection.Find("img[data-src]").Length()-1 && res["covers"] == ""  {
							if extend.GetFileSize(lopath) > 1024 * 10{
								res["covers"] = lopath
							}
						}
						img.SetAttr("src","/"+lopath)
					}else {
						img.SetAttr("src",src)
					}
				})
				fmt.Println(res["covers"])
				// 处理视频
				selection.Find("iframe[data-src]").Each(func(i int, iframe *goquery.Selection) {
					src,_ := iframe.Attr("data-src")
					src = strings.Replace(src,"https://v.qq.com/iframe/preview.html","http://v.qq.com/txp/iframe/player.html",-1)
					iframe.SetAttr("src",src)
					//iframe.Remove()
				})
				html,_ := selection.Html()
				// 背景图片
				reg := regexp.MustCompile(`background-image:.*?url\(&#34;([\s\S]*?)&#34;\);`)
				result := reg.FindAllStringSubmatch(html, -1)
				for _,v := range result{
					paths := v[len(v)-1]
					src := paths
					if "https://mmbiz.qpic.cn" == GetQueueUrlDomain(src){
						src = "http://img01.store.sogou.com/net/a/04/link?appid=100520029&url=" + src
						if lopath,e := DownloadResourceLocal(src);e == nil{
							html = strings.Replace(html,"url(&#34;"+paths,"url(&#34;"+"/"+lopath,-1)
						}else {
							html = strings.Replace(html,"url(&#34;"+paths,"url(&#34;"+src,-1)
						}
					}
				}
				res["content"] = RemoveLineHtml(html)
			})
			// 描述
			dom.Find("#activity-name").Each(func(i int, selection *goquery.Selection) {
				res["description"] = strings.ReplaceAll(strings.ReplaceAll(selection.Text(), " ", ""), "\n", "")
			})
			// 预览量
		}
	}
	return res
}

// 清理连续换行
func RemoveLineHtml(sHtml string) string  {
	reg, _ := regexp.Compile("\n")
	return reg.ReplaceAllString(sHtml, "")
}

// 获得连接URL更域名
func GetQueueUrlDomain(URL string) string {
	reg := regexp.MustCompile(`(https|http):\/\/([^\/]+)`)
	if result := reg.FindAllStringSubmatch(URL, -1);len(result) >0 && result[0][0] != ""{
		return result[0][0]
	}else {
		return ""
	}
}

// 下载远程资到本地
func DownloadResourceLocal(cloudPath string) (localPath string,err error) {
	if  res, e := http.Get(cloudPath); e == nil {
		fileName := path.Base(cloudPath)
		filePathSave := extend.CreateDateDir("static/upload/")
		fileNameSave := time.Now().Format("20060102150405") + extend.GetRandomString(5) + path.Ext(fileName)
		defer res.Body.Close()
		reader := bufio.NewReaderSize(res.Body, 32 * 1024)
		if file, e := os.Create(filePathSave + "/" + fileNameSave);e != nil{
			err = e
		}else {
			writer := bufio.NewWriter(file)
			_, _ = io.Copy(writer, reader)
			localPath = strings.Replace(filePathSave + "\\" + fileNameSave, "\\", "/", 3)
		}
	}else {
		err = e
	}
	if extend.GetFileSize(localPath) == 0{
		err = errors.New("下载失败！！！")
	}
	return
}
