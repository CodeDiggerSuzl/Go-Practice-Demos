package main

import (
	"fmt"
	"net"
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

// global channel, deliver msg message to all clients
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

// handle connection
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

	// 一个匿名 go 程,用,来发用户的消息
	go func() {
		// 循环读取消息
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				p("[ 📢  Detected @" + clnt.Name + " has disconnected ⤵️ ] \n")
				message <- "[ 📢 \a @" + clnt.Name + " has disconnected ⤵️ ] \n"
				return
			}
			if err != nil {
				p("err during conn.Read", err)
				return
			}
			// 获取用户指令 来判断用户意图

			msg := string(buf[:n-1])
			// get all client list
			if msg == "WHO" && len(msg) == 3 {
				conn.Write([]byte(" 👥  All online users:\n"))
				//  遍历 map
				for _, client := range onLineUserMap {
					userInfo := "🕸  IpAddress: " + client.Addr + " 🗿  Name: " + client.Name + "\n"
					_, err := conn.Write([]byte(userInfo))
					if err != nil {
						p("err during 'WHO'", err)
						return
					}
				}
				// rename : 判断开头是否为 RENAME2: && 利用切片来获取命名
				// "RENAME2:"
			} else if len(msg) >= 9 && msg[:8] == "RENAME2:" {
				newName := msg[8:]
				renameClnt(clnt, newName)
			} else {
				// send msg to all clnts
				message <- produceMsg(clnt, msg)
			}
		}
	}()
	// send "user is login" line to global channel
	// TODO
	p("[ 📢   #" + clnt.Addr + " | @" + clnt.Name + " - has login ⬆️  ]")
	message <- "[ 📢 \a  #" + clnt.Addr + " | @" + clnt.Name + " - has login ⬆️  ]"
	for {

	}
}

// rename current client
func renameClnt(clnt Client, name string) {
	clnt.Name = name
	onLineUserMap[clnt.Addr] = clnt
}

// produce message
func produceMsg(clnt Client, msg string) string {
	return "[ 📣 \a #" + clnt.Addr + " @" + clnt.Name + "] says: \n" + msg + "\n--------------------------"
}

// write message to client
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
