package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int

	//在线用户列表 key是用户名 value是用户实例
	OnlineMap map[string]*User
	//读写锁，用于互斥访问OnlineMap
	mapLock sync.RWMutex

	//消息广播的channel
	Message chan string
}

// Server的构造器
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel的go程，一旦有消息就发送给全部的在线User
// 在server启动时就开启
func (server *Server) ListenMessager() {
	for {
		msg := <-server.Message
		//将msg发送给全部的在线User
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
			user.Channel <- msg
		}
		server.mapLock.Unlock()

	}
}

// 将消息传入到server channel中
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

// 处理业务
func (server *Server) Handler(conn net.Conn) {
	//fmt.Println("New connection from ", conn.RemoteAddr())
	user := NewUser(conn)
	//用户上线，将用户加入到OnlineMap
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	//广播当前用户上线消息
	server.BroadCast(user, "已上线")

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			//n表示读出的数据长度
			n, err := conn.Read(buf) //读数据
			if n == 0 {
				server.BroadCast(user, "下线")
				//退出当前协程
				return
			}
			if err != nil && err != io.EOF { //读数据出错
				fmt.Println("conn Read err:", err)
				return
			}

			//提取用户的消息（去除'\n'）
			msg := string(buf[:n-1])
			//将得到的消息进行广播
			server.BroadCast(user, msg)
		}
	}()

	//当前handler阻塞，如果不阻塞，则执行完上条语句后，该go程直接结束了
	select {}
}

// 启动服务器
func (server *Server) Start() {
	//socket listen 使用tcp协议，利用Sprintf函数将ip和port拼接成地址
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听Message的go程
	go server.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		//处理业务
		go server.Handler(conn)
	}

}
