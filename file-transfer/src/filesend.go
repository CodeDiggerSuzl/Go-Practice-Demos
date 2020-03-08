package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

// file transfer using tcp
func main() {
	// get os args
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Please type the file path")
		return
	}
	filePath := args[1]

	// get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("err during os.Stat", err)
		return
	}
	name := fileInfo.Name()

	// send file name and receive "ok"
	conn, err := net.Dial("tcp", ":8848")
	defer conn.Close()
	if err != nil {
		fmt.Println("err during net.Dial", err)
		return
	}
	fmt.Println("connect to server was created correctly")
	_, err = conn.Write([]byte(name))
	if err != nil {
		fmt.Println("err during conn.Write", err)
		return
	}
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("err during conn.Read")
		return
	}
	if string(buf[:n]) == "ok" {
		// read and send file
		sendFile(conn, filePath)
	}

}

// send file
func sendFile(conn net.Conn, path string) {
	fmt.Println("start opening file")

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println("err during os.Open", err)
		return
	}

	buf := make([]byte, 4096)

	// while read while write
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("file read all done")
				return
			}
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println("err during conn.Write")
			return
		}
	}
}

//
// // main func of file send
// func main() {
// 	args := os.Args
// 	if len(args) != 2 {
// 		fmt.Println("❌ plz type the file name")
// 		return
// 	}
// 	// get the file name
// 	fileInfo, err := os.Stat(args[1])
// 	if err != nil {
// 		fmt.Println("❌ err during os.Stat", err)
// 		return
// 	}
// 	fileName := fileInfo.Name()
//
// 	// client send connection request
// 	conn, err := net.Dial("tcp", ":8848")
// 	if err != nil {
// 		fmt.Println("❌ err during net.Dial", err)
// 		return
// 	}
// 	fmt.Println("✅ conn created：", conn)
// 	// !! don't forget to write
// 	defer conn.Close()
// 	_, err = conn.Write([]byte(fileName))
// 	if err != nil {
// 		fmt.Println("❌ err during conn.Write(fileName)", err)
// 		return
// 	}
// 	// get info send form server
// 	buf := make([]byte, 16)
// 	n, err := conn.Read(buf)
// 	if err != nil {
// 		fmt.Println("❌ err during conn.Read", err)
// 		return
// 	}
// 	if string(buf[:n]) != "ok" {
// 		fmt.Println("❌ get from server is not string 'ok' ")
// 		return
// 	} else {
// 		sendFile(conn, args[1])
// 	}
// }
//
// // send file function
// func sendFile(conn net.Conn, filePath string) {
// 	// read file with ReadOnly
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("❌ err during os.Open", err)
// 		return
// 	}
// 	defer file.Close()
// 	// read from local and send to net
// 	buf := make([]byte, 4096)
// 	for {
// 		n, err := file.Read(buf)
// 		if err != nil {
// 			if err == io.EOF {
// 				fmt.Println("✅ File Transfer Completed")
// 			} else {
// 				fmt.Println("❌ err during file.Read", err)
// 			}
// 			return
// 		}
// 		_, err = conn.Write(buf[:n])
// 		if err != nil {
// 			fmt.Println("❌ err during conn.Write", err)
// 			return
// 		}
// 	}
// }
