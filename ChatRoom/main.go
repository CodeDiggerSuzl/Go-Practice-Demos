package main

import (
	"fmt"
	"net"
	"time"
)

// 用户结构体类型
// C : each client's channel
// name : client's name
// Addr : address
type Client struct {
	C    chan string
	Name string
	Addr string
}

// 全局map 存储线上用户
var onLineUserMap = make(map[string]Client)

// global channel, deliver user message
var message = make(chan string)

var p = fmt.Println

func main() {
	// 创建监听套接字
	listener, err := net.Listen("tcp", ":8848")
	defer listener.Close()
	if err != nil {
		p(" ❌ err during net.Listen", err)
		return
	}
	// 创建 manager，for global map and message
	go Manager()
	// 循环监听客户端连接请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			p("❌ err during listen.Close", err)
			return
		}
		// 启动 go 程 处理客户端请求
		go HandlerConnect(conn)
	}
}

// for global message channel and onLineUserMap
func Manager() {
	// 循环从 message 中读取数据
	for {
		// watch the global message if has data
		msg := <-message

		// send msg to each of online user map
		for _, clnt := range onLineUserMap {
			clnt.C <- msg
		}
	}
}

// 处理链接
func HandlerConnect(conn net.Conn) {
	defer conn.Close()
	// get user ip and port
	addr := conn.RemoteAddr().String()

	// create Client
	clnt := Client{make(chan string), addr, addr}

	// put new clnt into online user map, key:addr value : Clint
	onLineUserMap[addr] = clnt

	// 创建专门用来给当前用户送数据的 go 程
	go WriteMsg2Client(clnt, conn)

	// send "user is login" line to global channel
	// TODO
	message <- "[⬆️ " + string(time.Now().Second()) + " " + clnt.Name + " has login]"
	for {

	}
}

// this func is for write message to client
func WriteMsg2Client(clnt Client, conn net.Conn) {
	// 监听用户自带 channel 是否有消息
	// TODO
	for msg := range clnt.C {
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			p("❌ err during conn.Write", err)
			return
		}
	}
}
