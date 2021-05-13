package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播通道
	Message chan string
}

// NewServer 创建一个服务
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// BroadCast 广播消息
func (that *Server) BroadCast(user *User, msg string) {

	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	that.Message <- sendMsg
}

// Handler 业务处理
func (that *Server) Handler(conn net.Conn) {
	// 当前链接的任务
	fmt.Println("链接建立成功")

	// 用户上线，将用户加入到OnlineMap中
	user := NewUser(conn)

	that.mapLock.Lock()
	that.OnlineMap[user.Name] = user
	that.mapLock.Unlock()

	that.BroadCast(user, "已上线")

	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if n == 0 {
				that.BroadCast(user, "下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Err:", err)
				return
			}

			msg := string(buf[:n-1])

			// 将得到的消息进行广播

			that.BroadCast(user, msg)
		}
	}()

	// 当前handler阻塞
	select {}
}

// ListenerMessage 监听Message广播消息的进程，一旦有消息就发送给全部在线的User
func (that *Server) ListenerMessage() {

	for {

		msg := <-that.Message

		// 将msg发送给所有在线的用户
		that.mapLock.Lock()

		for _, cli := range that.OnlineMap {
			cli.C <- msg
		}

		that.mapLock.Unlock()

	}
}

// Start 启动服务器
func (that *Server) Start() {

	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", that.Ip, that.Port))

	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// close listen socket
	defer listener.Close()

	// 启动消息监听进程
	go that.ListenerMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}

		// do handler
		go that.Handler(conn)
	}

}
