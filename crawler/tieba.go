package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
)

// HTTPGet a page using
func HTTPGet(url string) (result string, err error) {
    resp, err1 := http.Get(url)
    if err1 != nil {
        err = err1
        return
    }
    defer resp.Body.Close()

    buf := make([]byte, 4096)
    for {
        n, _ := resp.Body.Read(buf)
        if n == 0 {
            log.Println("Read all ... ")
            break
        }
        result += string(buf[:n])
    }
    return
}

func getToWork(start, end int) {
    log.Printf("Crawling from page %d to page %d", start, end)
    for i := start; i <= end; i++ {
        url := "https://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn=" +
            strconv.Itoa((i-1)*50)
        fmt.Printf("正在爬：%d 页，%s\n", i, url)
        result, err := HTTPGet(url)
        if err != nil && err == io.EOF {
            fmt.Println("HttpGet err:", err)
            continue
        }
        fileName := strconv.Itoa(i) + ".html"
        f, err := os.Create(fileName)
        if err != nil {
            fmt.Println("Create err:", err)
            continue
        }
        _, _ = f.WriteString(result)
        _ = f.Close()
    }
}

func main() {
    var start, end int
    fmt.Printf("Please input the start page (>=1):")
    _, _ = fmt.Scan(&start)
    fmt.Printf("Please input the end page (>=start):")
    _, _ = fmt.Scan(&end)

    getToWork(start, end)
}
