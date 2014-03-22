package util

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	COMMEM_PARA_SEP = ":"
)

// ------------------------------------------------------------------
// ==> xml配置相关 Begin
// ------------------------------------------------------------------

// 装载Xml配置文件
// 参数格式由struct决定
func LoadXmlFile(confFile string, v interface{}) error {
	// 打开配置文件
	file, err := os.Open(confFile) // For read access.
	if err != nil {
		ShowError("LoadXmlFile read file err", confFile, err)
		return err
	}
	defer file.Close()

	// 读取数据
	data, err := ioutil.ReadAll(file)
	if err != nil {
		ShowError("LoadXmlFile ReadAll", err)
		return err
	}

	// 解析xml
	err = xml.Unmarshal(data, v)
	if err != nil {
		ShowError("LoadXmlFile Unmarshal err", err)
		return err
	}

	return nil
}

// ------------------------------------------------------------------
// ==> xml配置相关 End
// ------------------------------------------------------------------

// ------------------------------------------------------------------
// ==> txt配置相关 Begin
// ------------------------------------------------------------------

// 通用配置struct: 用于存储配置项
type CommonConfig struct {
	XMLName  xml.Name `xml:"Config"`
	FileName string
	Parames  []string `xml:"p"`
	strs     map[string]string
	pCount   int
}

func parseOnePara(paraStr string) (string, string, error) {
	if paraStr == "" {
		return "", "", Error("paraStr is empty")
	}

	ss := strings.SplitN(paraStr, COMMEM_PARA_SEP, 2)

	sNum := len(ss)
	var k, v string

	if sNum < 1 {
		return "", "", Error("paraStr format err")
	} else if sNum == 1 {
		k = strings.TrimSpace(ss[0])
		v = ""
	} else {
		k = strings.TrimSpace(ss[0])
		v = strings.TrimSpace(ss[1])
	}
	return k, v, nil
}

func (this *CommonConfig) Load(fileName string) error {
	err := LoadXmlFile(fileName, this)
	if err != nil {
		return err
	}

	this.FileName = fileName
	this.pCount = 0

	var k, v string

	pNum := len(this.Parames)

	for i := 0; i < pNum; i++ {
		k, v, err = parseOnePara(this.Parames[i])
		if err != nil {
			ShowError("CommonConfig.Load parse err", this.FileName, i, err)
			continue
		}

		err = this.SetConfString(k, v)
		if err != nil {
			ShowError("CommonConfig.Load set err", this.FileName, i, k, v, err)
			continue
		}
	}

	return nil
}

// 设置相关函数
func (this *CommonConfig) SetConfString(key string, value string) error {

	if this.strs == nil {
		this.pCount = 0
		this.strs = make(map[string]string, 20)
	}

	v, ok := this.strs[key]
	if ok {
		ShowWarnning("String Conf Was Redefined: ", key, " ", v, " => ", value)
	} else {
		this.pCount++
	}
	this.strs[key] = value
	return nil
}

// 获取相关函数
func (this *CommonConfig) GetConfString(key string) (string, error) {
	if this.strs == nil {
		return "", Error("GetConfString() config list is empty")
	}

	v, ok := this.strs[key]
	if ok {
		return v, nil
	}
	return "", Error("GetConfString() failed, key=(%s)", key)
}

func (this *CommonConfig) GetConfStringToInt64(key string) (int64, error) {
	v, err := this.GetConfString(key)
	if err != nil {
		return 0, err
	}

	if v != "" {
		n, err := strconv.ParseInt(v, 0, 64)
		if err != nil {
			return 0, Error("GetConfStringToInt64() failed key=(%s) v=(%s)", key, v)
		}
		return n, nil
	}

	return 0, Error("GetConfStringToInt64() failed key=(%s)", key)
}

func (this *CommonConfig) PrintAllConfigs() {
	ShowDebug("Config file name:", this.FileName, "paraCount", this.pCount)
	i := 0
	for k, v := range this.strs {
		ShowDebug(i, "\t", k, "\t", v)
		i++
	}
}

// 装载txt格式的配置文件
// 一行为一个配置项
// 采用 key:value 形式的参数
// value: 表示此参数为int64形式或者bool形式
// "value" 表示这个值是string形式
func LoadTxtFile(confFile string, cf *CommonConfig) error {
	// 打开配置文件
	file, err := os.Open(confFile) // For read access.
	if err != nil {
		ShowError("LoadTxtFile read file err", confFile, err)
		return err
	}
	defer file.Close()

	// 读取数据
	data, err := ioutil.ReadAll(file)
	if err != nil {
		ShowError("LoadTxtFile ReadAll", err)
		return err
	}

	// TODO:如何处理数据?

	if data != nil {
		return nil
	}
	return nil
}

// ------------------------------------------------------------------
// ==> txt配置相关 End
// ------------------------------------------------------------------
