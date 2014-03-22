/*
	socket消息处理句柄,目前尝试的处理流程为:
	1. 网络连接获取数据
	2. 将数据强行规定为 length(int32) +  data 的形式
	3. 将消息交给doMsg函数处理
	4. 返回消息
*/

package ulsocket

import (
	"net"
	"time"
	"ultralisk/util"
)

const (
	DEFAULT_RW_BUFF_SIZE = 1024 * 1024 // 默认的消息读写buff长度  1M
	TIMEOUT_ERROR_INFO   = "i/o timeout"
)

// 回调函数相关设置
const (
	CB_ON_RECIVE     = "onReciveCallback"    // 网络数据到达后的回调
	CB_ON_DO_MESSAGE = "onDoMessageCallback" // 处理函数
	CB_ON_WRITE      = "onWriteCallback"     // 数据写入回调
)

type OnReciveFunc func(handler *MessageHandler, buf []byte, bufSize int) error
type DoMessageFunc func(handler *MessageHandler) error
type OnWriteFunc func(handler *MessageHandler) error

func CheckTimeout(err error) bool {
	if err != nil {
		e, ok := err.(net.Error)
		if ok && e.Timeout() {
			return true
		}
	}
	return false
}

// 读写buffer数据结构体
type MessageHandler struct {
	handlerType string // 消息句柄的连接类型(http/tcp/udp)
	tcpAddr     string // 连接地址
	httpAddr    string
	needClose   bool

	// 逻辑处理函数句柄
	Conn      net.Conn               // 连接接口
	callbacks map[string]interface{} // 回调列表

	// 缓存的数据句柄
	SessionData interface{}
}

// 消息处理句柄初始化
func (this *MessageHandler) Init(conn net.Conn) error {

	if conn == nil {
		return util.Error("Init MessageHandler err conn=%v", conn)
	}

	// 句柄设置
	this.Conn = conn
	return nil
}

func (this *MessageHandler) SetCallback(key string, callback interface{}) {

	if key == "" || callback == nil {
		util.ShowError("MessageHandler.SetCallback err: key=", key, "callback=", callback)
		return
	}

	var ok bool
	switch key {
	case CB_ON_RECIVE:
		_, ok = callback.(OnReciveFunc)
		break
	case CB_ON_DO_MESSAGE:
		_, ok = callback.(DoMessageFunc)
		break
	case CB_ON_WRITE:
		_, ok = callback.(OnWriteFunc)
		break
	}

	if !ok {
		util.ShowError("MessageHandler.SetCallback faild: key=", key, "callback=", callback)
		return
	}

	this.callbacks[key] = callback
}

//// 消息数据到达
//func (this *MessageHandler) OnRecived(buf []byte, dataSize int) error {
//	if this.onReciveCallback != nil {
//		this.onReciveCallback(buf, dataSize)
//	}
//	return nil
//}

//// 消息处理逻辑
//func (this *MessageHandler) DoMessage() error {

//	doMsgFunc := this.doMsgCallback

//	if doMsgFunc == nil {
//		return errors.New("doMsgFunc is nil")
//	}

//	err := doMsgFunc(this, this.conn)
//	return err
//}
func (this *MessageHandler) SetSessionData(sessionData interface{}) {
	this.SessionData = sessionData
}
func (this *MessageHandler) RunTcp(heartBeat time.Duration) error {

	// 设置网络连接参数
	var timeout time.Time
	var err error
	var length int

	var OnRecive OnReciveFunc
	var DoMessage DoMessageFunc
	var OnWrite OnWriteFunc
	var cb interface{}

	var ok bool

	conn := this.Conn
	defer conn.Close()

	// 网络读取缓冲创建
	buffSize := 16 * 1024
	buffRead := make([]byte, buffSize)

	cb, ok = this.callbacks[CB_ON_RECIVE]
	if ok {
		OnRecive, ok = cb.(OnReciveFunc)
	}

	cb, ok = this.callbacks[CB_ON_DO_MESSAGE]
	if ok {
		DoMessage, ok = cb.(DoMessageFunc)
	}

	cb, ok = this.callbacks[CB_ON_DO_MESSAGE]
	if ok {
		OnWrite, ok = cb.(OnWriteFunc)
	}

	// 从网络中读取数据
	for {
		// 设置超时
		timeout = time.Now()
		conn.SetReadDeadline(timeout)
		conn.SetWriteDeadline(timeout)

		// 尝试从网络中读取数据
		length, err = conn.Read(buffRead)
		if err != nil {

			if CheckTimeout(err) {
				length = 0
			} else {
				// 网络错误, 强行断开连接
				conn.Close()
				util.ShowDebug("connection was closed", err.Error())
				break
			}
		}

		if length > 0 {
			// 读到数据的话,灌注到handler中
			buffRead[length] = 0 // 尾数置零
			err = OnRecive(this, buffRead, length)
			if err != nil {
				util.ShowError("onReciveCallback retrun err: ", err.Error())
				break
			}
		} else if length < 0 {
			// 读取长度异常
			conn.Close()
			util.ShowDebug("connection was closed")
			break
		}

		// 消息处理
		err = DoMessage(this)
		if err != nil {
			util.ShowError("doMsgCallback retrun err: ", err.Error())
			break
		}

		// 回写数据
		err = OnWrite(this)
		if err != nil {
			util.ShowError("onWriteCallback retrun err: ", err.Error())
			break
		}

		if this.needClose {
			conn.Close()
		} else {
			time.Sleep(heartBeat)
		}
	}

	conn.Close()
	return nil
}

func (this *MessageHandler) Close() {
	this.needClose = true
}

// 工厂函数
func CreateMessageHandler(conn net.Conn) (*MessageHandler, error) {
	ret := new(MessageHandler)
	err := ret.Init(conn)

	if err != nil {
		return nil, err
	}
	return ret, nil
}
