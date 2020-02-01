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

// 创建全局map，存储在线用户
var onLineUserMap map[string]Client

// 创建全局 channel 传递用户消息。
var message = make(chan string)

var p = fmt.Println

// handle connection
func HandlerConnect(conn net.Conn) {
	defer conn.Close()

	// get user ip and port
	curClntAddr := conn.RemoteAddr().String()
	// create Client
	clnt := Client{make(chan string), curClntAddr, curClntAddr}

	// put new clnt into online user map, key:curClntAddr value : Clint
	onLineUserMap[curClntAddr] = clnt

	// 创建专门用来给当前用户送数据的 go 程
	go WriteMsg2Client(clnt, conn)
	// send "user is login" line to global channel
	p("[ 📢   #" + clnt.Addr + " | @" + clnt.Name + " - has login ⬆️  ]")
	message <- "[ 📢 \a  #" + clnt.Addr + " | @" + clnt.Name + " - has login ⬆️  ]"

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

			switch {
			// list clients
			case msg == "WHO" && len(msg) == 3:
				ListClnts(conn)

			// rename
			case len(msg) >= 9 && msg[:8] == "RENAME2:":

				clnt.Name = msg[8:] // 修改结构体成员name
				onLineUserMap[clnt.Addr] = clnt
				// rename(clnt, msg) // 更新 onLineUserMap
				conn.Write([]byte("✅  Rename successfully\n"))

			default:
				// send msg to all clnts
				message <- produceMsg(clnt, msg)
			}
		}
	}()
	for {

	}
}

func rename(clnt Client, msg string) {
	clnt.Name = msg[8:] // 修改结构体成员name
	onLineUserMap[clnt.Addr] = clnt
}

// list all clients
func ListClnts(conn net.Conn) {
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
}

// produce message
func produceMsg(clnt Client, msg string) (buf string) {
	buf = "[ 📣 \a #" + clnt.Addr + " @" + clnt.Name + "] says: \n" + msg + "\n--------------------------"
	return
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

// for global message channel and onLineUserMap
func Manager() {
	// init map
	onLineUserMap = make(map[string]Client)
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
