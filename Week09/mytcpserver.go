package week09

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

// MyTCPServer 某某TCPServer
type MyTCPServer struct {
	l      *net.TCPListener
	ctx    context.Context
	cancel context.CancelFunc
}

// NewMyTCPServer ...
func NewMyTCPServer(addr string) (*MyTCPServer, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &MyTCPServer{
		l:      l,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Run 运行
func (m *MyTCPServer) Run() {
	fmt.Println("server run")
	for {
		select {
		case <-m.ctx.Done():
			m.l.Close()
			return
		default:
			conn, err := m.l.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			ch := make(chan []byte)
			go m.write(conn, ch)
			go m.read(conn, ch)
		}
	}
}

// Stop 退出
func (m *MyTCPServer) Stop() {
	m.cancel()
	fmt.Println("server stop")
}

// read
func (m *MyTCPServer) read(conn net.Conn, ch chan []byte) {
	ctx, _ := context.WithCancel(m.ctx)
	r := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return
		default:
			l, err := r.ReadBytes('\n')
			if err == io.EOF {
				time.Sleep(1 * time.Millisecond)
				continue
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(l) > 0 {
				ch <- l
			}
		}
	}
}

// write
func (m *MyTCPServer) write(conn net.Conn, ch chan []byte) {
	ctx, _ := context.WithCancel(m.ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-ch:
			ret := fmt.Sprintf("got: %s\n", string(data))
			_, err := conn.Write([]byte(ret))
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
