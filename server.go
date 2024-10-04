package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

// Server的构造器
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
	}
	return server
}

// 处理业务
func (server *Server) Handler(conn net.Conn) {
	fmt.Println("New connection from ", conn.RemoteAddr())
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
