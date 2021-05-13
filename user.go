package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser 创建一个用户的Api
func NewUser(conn net.Conn) *User {

	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听当前用户消息的通道
	go user.ListenMessage()

	return user
}

// ListenMessage 监听当前User channel的方法，一旦有消息 就直接发送给对端客户端
func (that *User) ListenMessage() {
	for {
		msg := <-that.C

		that.conn.Write([]byte(msg + "\n"))
	}
}
