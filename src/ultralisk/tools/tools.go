package tools

import (
	"crypto/md5"
	"fmt"
	"io"
)

// md5生成工具
// 结果为32小写md5字符串
func Md5(str string) (string, error) {
	h := md5.New()

	_, err := io.WriteString(h, str)
	if err != nil {
		return "", err
	}

	strMd5 := fmt.Sprintf("%x", h.Sum(nil))

	return strMd5, nil
}
