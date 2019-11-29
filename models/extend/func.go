package extend

import (
	"encoding/binary"
	"github.com/beatrice950201/GoRbac/models/cache"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// 判断是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 字符串加密
func HashAndSalt(password string) string {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// 验证加密串
func ComparePasswords(hashedPwd string, plainPassword string) bool {
	byteHash := []byte(hashedPwd)
	plainPwd := []byte(plainPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}

// Ip2long 将 IPv4 字符串形式转为 uint32
func Ip2long(ipString string) uint32 {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

// 写入缓存
func SetCache(key string, val interface{}) error {
	return  cache.Bm.Put(key, val, 86400*time.Second)
}

// 判断interface是否为空
func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

// 判断是否存在某个值
func InArray(need int, needArr []int) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

// 根据当前日期来创建文件夹
func CreateDateDir(Path string) string {
	folderName := time.Now().Format("20060102")
	folderPath := filepath.Join(Path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		_ = os.Mkdir(folderPath, 0777) // 必须分成两步：先创建文件夹、再修改权限
		_ = os.Chmod(folderPath, 0777)
	}
	return folderPath
}

//生成随机字符串
func GetRandomString(length int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

/*
  回退网络模式
*/
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// 获取文件大小
func GetFileSize(filename string) int64 {
	var result int64
	_ = filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}