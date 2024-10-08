package main

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn

	//当前用户属于哪个server，这样就可以利用user来访问server
	server *Server
}

// User的构造器
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
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

// 用户上线
func (user *User) Online() {
	//用户上线，将用户加入到OnlineMap
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	//广播当前用户上线消息
	user.server.BroadCast(user, "已上线")
}

// 用户下线
func (user *User) Offline() {
	//用户下线，将用户从OnlineMap删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	//广播当前用户下线消息
	user.server.BroadCast(user, "下线")
}

// 用户处理消息
func (user *User) DoMessage(msg string) {
	//目前只有广播
	user.server.BroadCast(user, msg)
}
