package week09

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestMyTCPServer(t *testing.T) {
	// 创建server
	addr := "127.0.0.1:23330"
	s, err := NewMyTCPServer(addr)
	if err != nil {
		t.Error(err)
		return
	}
	go s.Run()

	// 模拟一个客户端
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = conn.Write([]byte("test\n"))
	if err != nil {
		t.Error(err)
		return
	}
	rd := bufio.NewReader(conn)
	l, err := rd.ReadBytes('\n')
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(l))
	conn.Close()

	// 关闭服务端
	time.Sleep(1 * time.Second)
	s.Stop()
}
