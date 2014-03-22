package service

import (
	"fmt"
	"time"
	"ultralisk/net/ulsocket"
	"ultralisk/util"
	"ultralisk/util/cmd"
)

// 服务器的开启关闭配置常数
const (
	OPEN               = 1
	CLOSE              = 0
	DEFAULT_HEART_BEAT = time.Millisecond * 50
)

// 服务器的运行状态
const (
	RUNTIME_STATE_DEFAULT = iota // 默认状态(尚未初始化完毕)
	RUNTIME_STATE_RUNNING        // 运行中
	RUNTIME_STATE_DESTROY        // 已经关闭
)

// service 接口定义,每个具体的service均应该实现此接口
type Service interface {
	SetInfo(cf ServiceInfo) error                              // 设置信息
	GetInfo() *ServiceInfo                                     // 获取信息
	Run() error                                                // 启动
	Stop() error                                               // 结束
	Reset() error                                              // 复位(停止服务,并回复到初始状态)
	DoCmd(session *cmd.CmdSession, cmdData *cmd.CmdData) error // 命令行响应
	SetRuntimeState(state int) int                             // 状态设置
	GetRuntimeState() int                                      // 状态检查
}

type ServiceCreateFunc func(sType string) (*Service, error)

// 服务器中的service管理器
type ServiceManager struct {
	Conf       *ServerConfig       // 全服的当前配置信息
	Factory    ServiceCreateFunc   // Service的工厂方法
	sMap       map[string]*Service // 实际service列表
	cmdMachine *cmd.CmdMachine
}

func (this *ServiceManager) Init(factoryFunc ServiceCreateFunc) error {
	if factoryFunc == nil {
		return util.Error("factoryFunc is nil,init ServiceManager faild")
	}

	this.sMap = make(map[string]*Service)
	this.Conf = new(ServerConfig)
	this.Factory = factoryFunc
	return nil
}

func (this *ServiceManager) LoadConfig(cfgName string) error {
	var err error
	this.Conf, err = LoadServerConfig(cfgName)
	return err
}

// 启动各项服务
func (this *ServiceManager) StartServices(factoryFunc ServiceCreateFunc) error {

	// 确定工厂方法
	facF := factoryFunc
	if facF == nil {
		facF = this.Factory
	}

	if facF == nil {
		return util.Error("Factory is nil")
	}

	// 准备遍历Service列表
	list := this.Conf.Services
	sNum := len(list)

	var s *Service
	var sInfo *ServiceInfo
	var err error
	var k string

	// 遍历服务器信息
	for i := 0; i < sNum; i++ {
		sInfo = &list[i]

		k = sInfo.Name
		s = this.Get(k)

		// 从表中查找当前的信息
		if s != nil {
			util.ShowWarnning(fmt.Sprintf("service %s was exist when start", k))
			(*s).Reset()     // 对已经存在的service进行重置
			this.Set(k, nil) // 清理原有的对象
			s = nil
		} else {
			if sInfo.Status == OPEN {
				s, err = facF(sInfo.Type)
				if s == nil || err != nil {
					util.ShowError("service factory err:", k, err)
				}
			}
		}

		if s != nil && sInfo.Status == OPEN {
			// 启动有效Service 的goroutin
			(*s).SetInfo(*sInfo)
			this.Set(k, s)
			util.ShowInfo(fmt.Sprintf("service run type=%s , name = %s", sInfo.Type, sInfo.Name))
			go (*s).Run()
		}
	}

	return nil
}

// 获取指定的Service句柄
func (this *ServiceManager) Get(key string) *Service {
	if key == "" {
		return nil
	}

	s, ok := this.sMap[key]
	if !ok || s == nil {
		return nil
	}
	return s
}

// 添加一个Service到ServerManager中
// s == nil 表示删除
func (this *ServiceManager) Set(key string, s *Service) (*Service, error) {
	if key == "" {
		return nil, util.Error("ServiceInfo.Set() key is empty")
	}

	ts, _ := this.sMap[key]
	this.sMap[key] = s
	return ts, nil
}

func (this *ServiceManager) StartCmd(ip, port, password, title string) error {
	cm := cmd.NewCmdMachine(password, title, this.DoCmd)
	this.cmdMachine = cm
	go ulsocket.StartTcp(ip, port, cm.DoTcp)
	return nil
}

func (this *ServiceManager) DoCmd(session *cmd.CmdSession, cmdData *cmd.CmdData) error {
	return nil
}
