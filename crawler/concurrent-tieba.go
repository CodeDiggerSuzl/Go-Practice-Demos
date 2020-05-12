package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func GetPage(url string) (result string, err error) {
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
			fmt.Println("Read the page done...")
			break
		}
		result += string(buf[:n])
	}
	return
}

// 定义一个名为 page 的通道，将通道引用和循环因子 i 一起传递到了 SpiderPage 方法中。
// 同时在 SpiderPage 函数内爬取网页数据完成后，将 i 值（代表爬取的第几页）写入 page。
// 主协程循环创建 N 个 goroutine 之后，要依次读取每一个 goroutine 借助 channel 写回的 i 值。
// 在读取期间，如果 page 上没有写端写入，主 goroutine 则会阻塞等待，直到有子协程写入，读取打印第 i 个页面爬取完毕。
func SpiderAPage(idx int, page chan int) {
	url := "https://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn=" + strconv.Itoa((idx-1)*50)
	fmt.Printf("Crawling page %d, the url is: %s\n", idx, url)
	result, err := GetPage(url)
	if err != nil {
		fmt.Println("GetPage err", err)
		return
	}
	fileName := strconv.Itoa(idx) + ".html"
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Create file error:", err)
		return
	}

	_, _ = f.WriteString(result)
	_ = f.Close()

	// When a page is done, send to channel
	fmt.Printf("Send %d to page", idx)
	page <- idx
}

func StartWorking(start, end int) {
	// Use channel to avoid the main channel to exit
	page := make(chan int)

	for i := start; i <= end; i++ {
		go SpiderAPage(i, page)
	}
	// The next for loop will exec after the last for loop
	for i := start; i <= end; i++ {
		fmt.Printf("Read for channel <-page %d\n", <-page)
		fmt.Printf("Crawling %d page done ✅ \n", <-page)
	}
}

func main() {
	var s, e int
	fmt.Println("Enter start page")
	_, _ = fmt.Scan(&s)
	fmt.Println("Enter end page")
	_, _ = fmt.Scan(&e)
	StartWorking(s, e)
}
