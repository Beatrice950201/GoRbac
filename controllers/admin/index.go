package admin

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"runtime"
)

type IndexController struct {
	MainController
}

// @router /index/index [get]
func (c *IndexController) Index()  {
	c.Data["GoVersion"] = runtime.Version()
	c.Data["CpuInfo"] = c.GetCpuInfo()
	c.Data["Memory"] = c.GetMemory()
	c.Data["MemoryUsedPercent"] = c.UsedPercentMemory()
	c.Data["SystemVersion"] = runtime.GOOS
}

// 获取内存使用率
func (c *IndexController) UsedPercentMemory() float64 {
	VirtualMemory,_:= mem.VirtualMemory()
	return 100 - VirtualMemory.UsedPercent
}

//CPU信息
func (c *IndexController) GetCpuInfo() string {
	info,_ := cpu.Info()
	return info[0].ModelName
}

// 内存信息
func (c *IndexController) GetMemory() uint64 {
	VirtualMemory,_:= mem.VirtualMemory()
	return VirtualMemory.Total / 1024 / 1024 / 1024
}

