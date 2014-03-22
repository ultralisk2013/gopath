package service

import (
	"encoding/xml"
	"runtime"
	"ultralisk/util"
)

type ServerConfig struct {
	ServerName string
	Version    string
	DebugMode  bool
	IP         string
	Port       string
	Services   []ServiceInfo `xml:"Service"`
	GOMAXPROCS int           // 允许使用的最大cpu数量
}

type ServiceInfo struct {
	XMLName     xml.Name `xml:"Service"`
	Type        string
	Name        string
	ConfigFile  string
	Params      string
	DebugMode   bool
	Description string
	Status      int32
}

// 获取指定Service的配置参数
func (this *ServerConfig) GetServiceInfo(serviceName string) *ServiceInfo {

	for _, v := range this.Services {
		if v.Name == serviceName {
			return &v
		}
	}
	return nil
}

// 获取指定Service的配置参数
func (this *ServerConfig) GetServiceStatus(serviceName string) int32 {

	for _, v := range this.Services {
		if v.Name == serviceName {
			return v.Status
		}
	}
	return CLOSE
}

// 打印配置信息
func (this *ServerConfig) PrintConfigInfo() {
	var st string
	var debugDesc string

	// 服务器名称以及版本号
	util.ShowInfo(""+this.ServerName, "(", this.Version, ")")

	// 服务器监听情况
	util.ShowInfo("Listen Addr\t", this.IP+":"+this.Port)

	// 调试模式开关
	if this.DebugMode {
		st = "Open"
	} else {
		st = "Close"
	}
	util.ShowInfo("DebugMode\t", st)

	// CPU 相关信息
	util.ShowInfo("LOC_NUMCPU\t", runtime.NumCPU())
	util.ShowInfo("GOMAXPROCS\t", this.GOMAXPROCS)

	// 各项服务器信息
	util.ShowInfo("Service\t", "Index\t", "Key\t", "Status\t", "ConfigFile\t", "Params\t", "Description\t")

	var count int32 = 0
	for k, v := range this.Services {

		if v.Status == 1 {
			st = "Open"
			count++
		} else {
			st = "Close"
		}

		if v.DebugMode {
			debugDesc = "(debug)"
		} else {
			debugDesc = ""
		}

		util.ShowInfo("Service\t", k, "\t", v.Name+debugDesc, "\t", st+"\t", v.ConfigFile+"\t", v.Params+"\t", v.Description)
	}
	util.ShowInfo("Service Count:\t", count, "\t\t\t")

}

// 装载服务器配置
func LoadServerConfig(fileName string) (*ServerConfig, error) {
	if fileName == "" {
		util.ShowInfo("Config file name is nil")
		return nil, util.Error("config file name is nil")
	}
	cf := new(ServerConfig)
	err := util.LoadXmlFile(fileName, cf)

	if err != nil {
		return nil, err
	} else {
		return cf, nil
	}
}
