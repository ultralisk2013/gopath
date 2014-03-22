package ulsocket

import (
	"fmt"
	"net"
	"ultralisk/util"
)

type DoTcpFunc func(conn net.Conn) error

func StartTcp(ip, port string, doTcpFunc DoTcpFunc) error {
	var err error

	if ip == "" || port == "" {
		return util.Error("startCmd addr err: ip=%s,port=%s", ip, port)
	}

	addr := fmt.Sprintf("%s:%s", ip, port)

	l, err := net.Listen("tcp", addr) // 启动监听
	if err != nil {
		return err
	}

	util.ShowDebug("StartCmd tcp starting :", addr)

	// 连接监听
	var c net.Conn
	for {
		c, err = l.Accept()
		if err != nil {
			if c != nil {
				err = c.Close()
				if err != nil {
					util.ShowError("cmd l.Accept err", err)
				}
			}
		} else {
			util.ShowInfo("cmd connection ok:", c.LocalAddr(), c.RemoteAddr())
			// 有新连接建立
			go doTcpFunc(c)
		}
	}
	return nil
}
