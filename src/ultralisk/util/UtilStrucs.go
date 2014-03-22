package util

/**
公共数据定义区
用于定义项目内各个模块均要用到的公共数据结构
*/
import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	UTF8_BOM     = "\xEF\xBB\xBF"
	UTF8_BOM_LEN = len(UTF8_BOM)
)

func CheckUTF8_BOM(data []byte) bool {
	if data == nil {
		return false
	}
	length := len(data)
	if length < UTF8_BOM_LEN {
		return false
	}

	if string(data[0:UTF8_BOM_LEN]) == UTF8_BOM {
		return true
	}
	return false
}

// session的数据结构
type Session struct {
	SessionId string      // session id
	Data      interface{} // session data
}

// 用户账号的数据结构
type User struct {
	SessionId     string // 用户的 session id
	SafeCode      string // 账号安全码（游戏中的二级安全验证码）
	SafeLevel     int32  // 用户的安全等级（类似于GM等级）
	Id            int64  // 用户的账号id
	Name          string // 名称
	RealName      string // 用户的真实姓名
	UpId          int64  // 推广ID
	Role          string // 权限值
	LoginCount    int32  // 登陆次数
	LastLoginTime int64  // 最后登陆时间
	LastLoginIp   string // 最后登陆IP
	BuildTime     int64  // 发生时间
	CreateTime    int64  // 记录时间
	CreateIp      string // 提交IP
	Sex           int64  // 性别
	Age           int16  // 年龄
	Area          int64  // 区域ID
	Mobile        string // 手机号
	Email         string // email
	Qq            string // qq
	MicroChannel  string // 微信
	Money0        int64  // 货币1 货币名称从config表里取
	Money1        int64  // 货币2 货币名称从config表里取
	Money2        int64  // 货币3 货币名称从config表里取
	Money3        int64  // 货币3 货币名称从config表里取
	Money4        int64  // 货币3 货币名称从config表里取
}

// 命令行处理相关数据和接口
type CmdInfo struct {
	SourceCmd string
	list      []string
}

func (this *CmdInfo) ParseCmd(srcStr string) error {
	return nil
}

type CommonDataPool struct {
	datas map[string]*CommonData
}

func (this *CommonDataPool) Init() {
	this.datas = make(map[string]*CommonData)
}

func (this *CommonDataPool) Put(key string, data *CommonData) error {
	if key == "" || data == nil {
		return Error("put common data err key=%s,data=%v", key, data)
	}
	this.datas[key] = data
	return nil
}

func (this *CommonDataPool) Get(key string) *CommonData {
	if key == "" {
		return nil
	}
	v, ok := this.datas[key]
	if !ok {
		return nil
	}
	return v
}

// 统一数据格式, 对应data.json
type CommonData struct {
	FileName string     // 原始文件的名字
	JsonStr  string     // 原始的json字符串
	Fields   []string   `fields` // 数据名称列表
	Types    []string   `types`  // 数据类型列表
	Values   [][]string `values` // 数据表
	ValueBuf [][]byte   // 每个value对应的json字符串(相当于预先打包好的单条数据的jsonStr)
}

// 从json格式中读取统一数据
func (this *CommonData) DecodeJsonFile(filename string) error {
	var e error

	f, e := os.Open(filename)
	if e != nil {
		return e
	}
	defer f.Close()

	var buf, jsonBuffer []byte

	bufLen := 1024
	buf = make([]byte, bufLen) // read buf

	var n, count int
	count = 0
	for {
		n, e = f.Read(buf)
		if e != nil || n <= 0 {
			break
		}

		if jsonBuffer == nil { // 首次读取
			if n < bufLen {
				jsonBuffer = buf[0:n]
				break
			}
			jsonBuffer = make([]byte, bufLen<<2)
		}

		jsonBuffer = append(jsonBuffer[0:count], buf[:n]...)
		count += n
	}

	// check uft8 bom
	if CheckUTF8_BOM(jsonBuffer) {
		jsonBuffer = jsonBuffer[UTF8_BOM_LEN:]
	}

	e = json.Unmarshal(jsonBuffer, this)
	if e != nil {
		return e
	}
	this.FileName = filename
	this.JsonStr = string(jsonBuffer)
	return nil
}

func (this *CommonData) Print() {
	ShowDebug(fmt.Sprintf("\n\n-----start print file \"%s\"------------------------", this.FileName))
	if this.Fields != nil {
		ShowDebug("fields <", len(this.Fields), "> data")
		for k, v := range this.Fields {
			ShowDebug("\t", k, v, "\t")
		}
	}

	if this.Types != nil {
		ShowDebug("types <", len(this.Types), "> data")
		for k, v := range this.Types {
			ShowDebug("\t", k, v, "\t")
		}
	}

	if this.Values != nil {
		ShowDebug("values <", len(this.Values), "> data")
		for k, v := range this.Values {
			ShowDebug("\tvalue", k, " <", len(v), "> data")
			for x, y := range v {
				ShowDebug("\t\tvalues", k, x, y, "\t")
			}
		}
	}
	ShowDebug(fmt.Sprintf("\n-----end print file \"%s\"------------------------\n\n", this.FileName))

}

// 通用的消息数据结构
type Message struct {
	Version int32
	Cmd     int64
	Data    interface{}
}

func (this *Message) ParseJsonString(data []byte) error {
	return json.Unmarshal(data, this)
}

func (this *Message) PackJsonString() ([]byte, error) {
	return json.Marshal(this)
}

func (this *Message) PrintInfo() {
	ShowDebug("msg cmd:", this.Cmd)
	ShowDebug("msg ver:", this.Version)
	ShowDebug("msg dat:", this.Data)
}
