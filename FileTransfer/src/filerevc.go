package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Listen to create socket
	listen, err := net.Listen("tcp", ":8848")
	defer listen.Close()
	if err != nil {
		fmt.Println("err during net.Listen", err)
		return
	}

	// accept
	conn, err := listen.Accept()
	defer conn.Close()
	if err != nil {
		fmt.Println("err during listen.Accept()", err)
		return
	}

	// get file name and return "ok"
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("err during conn.Read", err)
		return
	}
	fileName := string(buf[:n])
	_, err = conn.Write([]byte("ok"))
	if err != nil {
		fmt.Println("err during conn.Write", err)
		return
	}

	// save file
	saveFile(conn, fileName)
}

func saveFile(conn net.Conn, name string) {
	// create file
	file, err := os.Create(name)
	if err != nil {
		fmt.Println("err during os.Create", err)
		return
	}
	defer file.Close()
	buf := make([]byte, 4096)
	for {
		n, _ := conn.Read(buf)
		if n == 0 {
			fmt.Println("file recv all done")
			return
		}
		file.Write(buf[:n])
	}
}

// func main() {
// 	fmt.Println("✅ staring to create socket of tcp server")
// 	listener, err := net.Listen("tcp", ":8848")
// 	if err != nil {
// 		fmt.Println("❌ during net.Listen", err)
// 		return
// 	}
// 	defer listener.Close()
// 	conn, err := listener.Accept()
// 	if err != nil {
// 		fmt.Println("❌ err during listener.Accept()", err)
// 		return
// 	}
// 	defer conn.Close()
// 	// get file name and return "ok",file name should be less than 1024 bytes
// 	buf := make([]byte, 1024)
// 	n, err := conn.Read(buf)
// 	if err != nil {
// 		fmt.Println("❌ err during conn.Read", err)
// 		return
// 	}
// 	fileName := string(buf[:n])
// 	_, err = conn.Write([]byte("ok"))
// 	if err != nil {
// 		fmt.Println("❌ err during conn.Write", err)
// 	}
//
// 	// file save
// 	recvFile(conn, fileName)
// }
//
// func recvFile(conn net.Conn, name string) {
// 	// create file
// 	file, err := os.Create(name)
// 	defer file.Close()
// 	if err != nil {
// 		fmt.Println("❌ err during create file")
// 		return
// 	}
// 	buf := make([]byte, 4096)
// 	for {
// 		n, _ := conn.Read(buf)
// 		if n == 0 {
// 			fmt.Println("✅ file receive completed")
// 			return
// 		}
// 		file.Write(buf[:n])
// 	}
// }
