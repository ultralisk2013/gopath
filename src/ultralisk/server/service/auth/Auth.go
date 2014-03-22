package auth

/*
	本模块功能为用户游戏账号的认证
	暂时只实现http一种认证方式
*/

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"ultralisk/net/ulhttp"
	"ultralisk/net/ulsocket"
	"ultralisk/util"
	"ultralisk/util/service"
)

type ServiceAuth struct {
	serviceInfo   service.ServiceInfo
	serviceConfig util.CommonConfig
	titleOnShow   string
	runtimeState  int

	// 本Service用到的配置段
	serverId         int64  // 服务器唯一id
	serverKey        string // 服务器key __Q:可能不需要这个东西
	debugMod         bool   // 调试开关
	httpIp, httpPort string // http监听的地址和端口
	httpReadBuff     []byte // http的读取缓冲
	tcpIp, tcpPort   string // tcp监听的地址和端口
	tcpReadBuff      []byte // tcp的读取缓冲
	// 其他数据

}

////////////////////////////////////////////////////////////////////////////////
// Service interface 的实现
////////////////////////////////////////////////////////////////////////////////
// 设置信息
func (this *ServiceAuth) SetInfo(cf service.ServiceInfo) error {
	this.serviceInfo = cf
	this.titleOnShow = cf.Name
	return nil
}

// 获取信息
func (this *ServiceAuth) GetInfo() *service.ServiceInfo {
	return &this.serviceInfo
}

// 启动
func (this *ServiceAuth) Run() error {

	this.titleOnShow = this.serviceInfo.Name

	if !this.onInit() {
		return util.Error("service <%s> start failed", this.serviceInfo.Name)
	}
	this.showInfo("starting")

	return nil
}

// 结束
func (this *ServiceAuth) Stop() error {
	return nil
}

// 复位
func (this *ServiceAuth) Reset() error {
	return nil
}

// 命令行响应
func (this *ServiceAuth) DoCmd(cmd *util.CmdInfo) (string, error) {
	return "", nil
}

// 状态设置
func (this *ServiceAuth) SetRuntimeState(state int) int {
	state, this.runtimeState = this.runtimeState, state
	return state
}

// 状态检查
func (this *ServiceAuth) GetRuntimeState() int {
	return this.runtimeState
}

////////////////////////////////////////////////////////////////////////////////
// 自身逻辑的实现
////////////////////////////////////////////////////////////////////////////////
// service初始化
func (this *ServiceAuth) onInit() bool {
	var err error

	// 各种初始化

	// 配置文件载入和解析
	fmt.Println()
	cfgFileName := this.serviceInfo.ConfigFile
	err = this.loadConfigs(cfgFileName)
	if err != nil {
		this.showError("load config err", cfgFileName, err)
		return false
	}

	err = this.parseConfigs()
	if err != nil {
		this.showError("parse config err", cfgFileName, err)
		return false
	}
	fmt.Println()

	// 启动网络监听
	err = this.startHttpService(this.httpIp, this.httpPort)
	if err != nil {
		this.showError("auth.onInit startHttpService:", err.Error())
		return false
	}

	err = this.startSocketService(this.tcpIp, this.tcpPort)
	if err != nil {
		this.showError("auth.onInit startSocketService:", err.Error())
		return false
	}
	//err = this.startSocketService("", "")
	//if err != nil {
	//	this.showError("auth.onInit startSocketService:", err.Error())
	//	return false
	//}

	return true
}

// service结束
func (this *ServiceAuth) onFinal() error {
	return nil
}

// 装载配置文件
func (this *ServiceAuth) loadConfigs(cfgFileName string) error {

	var err error

	err = this.serviceConfig.Load(cfgFileName)

	if err != nil {
		this.showError("parseConfig err", cfgFileName, err)
	} else {
		this.serviceConfig.PrintAllConfigs()
	}
	return err
}

// 解析配置文件
func (this *ServiceAuth) parseConfigs() error {

	// 将参数依次取出并使用
	var err error
	var key, str string
	var i64 int64
	cnf := &this.serviceConfig

	// 服务器统一id
	key = "serverId"
	i64, err = cnf.GetConfStringToInt64(key)
	if err != nil {
		this.showError(key, err)
	} else {
		this.serverId = i64
	}

	// http监听配置
	key = "httpIp"
	str, err = cnf.GetConfString(key)
	if err != nil {
		this.showError(key, err)
	} else {
		this.httpIp = str
	}

	key = "httpPort"
	str, err = cnf.GetConfString(key)
	if err != nil {
		this.showError(key, err)
	} else {
		this.httpPort = str
	}

	// tcp监听配置
	key = "tcpIp"
	str, err = cnf.GetConfString(key)
	if err != nil {
		this.showError(key, err)
	} else {
		this.tcpIp = str
	}

	key = "tcpPort"
	str, err = cnf.GetConfString(key)
	if err != nil {
		this.showError(key, err)
	} else {
		this.tcpPort = str
	}

	return err
}

// 启动http监听服务
func (this *ServiceAuth) startHttpService(ip, port string) error {

	if ip == "" || port == "" {
		return util.Error("startHttpService addr err: ip=%s,port=%s", ip, port)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", this.httpOnRequest)

	addr := fmt.Sprintf("%s:%s", ip, port)
	server := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	httpHandler, err := ulhttp.CreateHttpHandler(server)
	if err != nil {
		return err
	}

	this.showDebug("service auth http starting :", addr)
	go httpHandler.RunHttp()

	return nil
}

// 启动socket监听服务
func (this *ServiceAuth) startSocketService(ip, port string) error {
	ulsocket.StartTcp(ip, port, this.doTcp)
	return nil
}

// tcp逻辑处理

func (this *ServiceAuth) doTcp(conn net.Conn) error {

	var err error
	//var length int
	var handler *ulsocket.MessageHandler

	// 创建消息处理句柄并设置参数
	rbSize, wbSize := 0, 0
	handler, err = ulsocket.CreateMessageHandler(conn)
	if err != nil {
		this.showError("CreateMessageHandler err,", rbSize, wbSize, conn, this.tcpDoMessageCallback)
		return nil
	}

	// 设置回调
	handler.SetCallback(ulsocket.CB_ON_RECIVE, this.tcpOnReceiveCallback)
	handler.SetCallback(ulsocket.CB_ON_DO_MESSAGE, this.tcpDoMessageCallback)
	handler.SetCallback(ulsocket.CB_ON_WRITE, this.tcpOnWriteCallback)

	// 启动
	err = handler.RunTcp(time.Millisecond * 50)
	if err != nil {
		this.showError("do tcp err :", err.Error())
	}
	return nil
}

func (this *ServiceAuth) httpOnRequest(w http.ResponseWriter, req *http.Request) {

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil || buf == nil {
		return
	}

	length := len(buf)

	str := string(buf[0:length])
	this.showInfo("this is ServiceAuth.httpOnRequest", str)
	w.Write([]byte(str))
	return
}

// MessageHandler的三个回调
func (this *ServiceAuth) tcpOnReceiveCallback(buf []byte, bufSize int) error {
	return nil
}

func (this *ServiceAuth) tcpDoMessageCallback(handler *ulsocket.MessageHandler) error {
	return nil
}

func (this *ServiceAuth) tcpOnWriteCallback(handler *ulsocket.MessageHandler) error {
	return nil
}
