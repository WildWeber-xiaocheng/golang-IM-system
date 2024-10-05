package main

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
}

// User的构造器
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
	}
	//启动监听当前user channel消息的go程
	go user.ListenMessage()

	return user
}

// 监听当前User channel， 一旦有消息，就直接发送给对应的客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.Channel
		user.conn.Write([]byte(msg + "\n"))
	}
}
