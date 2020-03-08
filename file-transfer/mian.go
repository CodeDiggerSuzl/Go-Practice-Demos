package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("USAGE:\n go run main.go [arguments]")
    // get command line args ,split by space 
    args := os.Args
    if len(args) != 2{
        fmt.Println("❌ format  go run **.go fileName")
        return
    }
    path := args[1]
    // get file status
    fileInfo, err := os.Stat(path)
    if err != nil {
        fmt.Println("❌ err during os.Stat")
        return
    }
    name := fileInfo.Name()
    fmt.Println("file name is :",name)
}
