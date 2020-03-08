package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// 用户结构体类型
// C : each client's channel
// name : client's name
// Addr : address
type Client struct {
	id   int
	C    chan string
	Name string
	Addr string
}

// 创建全局map，存储在线用户
var onLineUserMap map[string]Client

// 创建全局 channel 传递用户消息。
var message = make(chan string)

// client id
var clntId = 0
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
		clntId++
		// 启动 go 程 处理客户端请求
		go HandlerConnect(conn)
	}
}

// handle connection
func HandlerConnect(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("🗣" +
		"You are now in a public chat room,with all user inside.\n\n" +
		"USAGE:\n\tWHO: To list all the users online\n\t" +
		"RENAME2:<name>: To change your user name\n\t" +
		"EXIT: To exit the chat room\n\t" +
		"CHAT:<UserId>: To start a chat with one user.\n" +
		"Here are the current users:\n"))
	// get user ip and port
	curClntAddr := conn.RemoteAddr().String()
	// create Client
	clnt := Client{clntId, make(chan string), curClntAddr, curClntAddr}
	// send "user is login" line to global channel
	p("[ 📢   #" + clnt.Addr + " | @" + clnt.Name + " - has login ⬆️  ]")
	message <- "[ 📢 \a  #" + clnt.Addr + " | @" + clnt.Name + " - has login ⬆️  ]"
	// put new clnt into online user map, key:curClntAddr value : Clint
	onLineUserMap[curClntAddr] = clnt
	ListClnts(conn)
	// 创建专门用来给当前用户送数据的 go 程
	go WriteMsg2Client(clnt, conn)

	// a chan watch the client is quit or not
	quitStat := make(chan bool)
	// a chan watch the client is typing or not
	hasData := make(chan bool)
	// 一个匿名 go 程,用,来发用户的消息
	go func() {
		// 循环读取消息
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)
			if n == 0 {
				quitStat <- true
				p("[ 📢  Detected @" + clnt.Name + " has disconnected ⤵️ ] \n")
				return
			}
			if err != nil {
				p("err during conn.Read", err)
				return
			}
			// 获取用户指令 来判断用户意图 去除 \n
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
				// CHAT:2
			case len(msg) >= 6 && msg[:5] == "CHAT:":
				ChatWithUser(conn, msg)
			case msg == "EXIT" && len(msg) == 4:
				quitStat <- true
				conn.Write([]byte("bye~ \n"))
				conn.Close()
				return
			default:
				// if code upto here which means has data flo
				hasData <- true
				// send msg to all clnts
				message <- produceMsg(clnt, msg)
			}
		}
	}()
	// use loop to watch the quitStat has data flow or not
	for {
		select {
		case <-quitStat:
			// 关闭 channel,终止 write to client
			close(clnt.C)
			// delete clnt in map
			delete(onLineUserMap, clnt.Addr)
			// broadcast clint
			message <- "[ 📢 \a @" + clnt.Name + " has disconnected ⤵️ ] \n"
			return
		case <-hasData:
			// doing nothing; refresh next case
		case <-time.After(2 * time.Minute):
			conn.Write([]byte("You have been kicked out,due to 120s time expire.Bye ~"))
			// ❌ the code bellow doesn't work
			// clnt.C <- "You have been kicked out,due to 120s time expire.Bye ~"
			delete(onLineUserMap, clnt.Addr)
			message <- "[ 📢 \a @" + clnt.Name + " has been kick out,due to 120s time expire ⤵️ ] \n"
			return
		}
	}
}

// 1. chat with user by id and send
// 2. after the specific user get agreed
// 3. chat with user and stop get msg from public room
// 4. after stop signal wo get into public room
func ChatWithUser(conn net.Conn, msg string) {
	// get user id
	splitWithSpace := strings.Split(msg, " ")
	if len(splitWithSpace) >= 2 {
		conn.Write([]byte("CHAT:<UserId>: To start a chat with one user."))
		return
	}
	count := 0
	s := strings.Split(splitWithSpace[0], ":")
	for _, v := range onLineUserMap {
		k, _ := strconv.Atoi(s[1])
		if v.id == k {
			count++
		}
	}
	switch {
	case count == 0:
		conn.Write([]byte("Can't find the user\n"))
		return
	case count >= 2:
		conn.Write([]byte("Find more than one user\n"))
		return
	default:

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
		id := strconv.Itoa(client.id)
		userInfo := "🥕 UserId: " + id + "🕸  IpAddress: " + client.Addr + " 🗿  Name: " + client.Name + "\n"
		_, err := conn.Write([]byte(userInfo))
		if err != nil {
			p("err during 'WHO'", err)
			return
		}
	}
}

// produce message
func produceMsg(clnt Client, msg string) (buf string) {
	return "[ 📣 \a #" + clnt.Addr + " @" + clnt.Name + "] says: \n" + msg + "\n--------------------------"
}

// write message to client
func WriteMsg2Client(clnt Client, conn net.Conn) {
	// 监听用户自带 channel 是否有消息
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
