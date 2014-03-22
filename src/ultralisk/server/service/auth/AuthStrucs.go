package auth

import (
	"ultralisk/util"
)

// 玩家列表类型定义
type tUserList struct {
	list    []util.User
	mapList map[int64]util.User
}

// 创建一个玩家列表
func makeUserList(capacity int64) (*tUserList, error) {
	return nil, nil
}

// 玩家列表初始化
func (this *tUserList) init() error {
	return nil

}
