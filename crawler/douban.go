package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "regexp"
    "strconv"
    "time"
)

// 发起请求，获取网页内容
func HttpGet(url string) (result string, err error) {
    resp, err1 := http.Get(url) // 发送get请求
    if err1 != nil {
        err = err1
        return
    }
    defer resp.Body.Close()

    // 读取网页内容
    buf := make([]byte, 4*1024)
    for {
        n, _ := resp.Body.Read(buf)
        if n == 0 {
            break
        }
        result += string(buf[:n]) // 累加读取的内容
    }
    return
}

func CrawlPages(i int, page chan int) {
    log.Println("Crawling page: " + strconv.Itoa(i))
    url := "https://movie.douban.com/top250?start=" + strconv.Itoa(i*25) + "&filter="

    time.Sleep(1 * time.Second)

    log.Println("Crawling url: ", url)
    result, err := HttpGet(url)
    if err != nil {
        log.Println("Error during err: ", err)
        return
    }
    re := regexp.MustCompile(`<span>(.*?)人评价</span>`)
    if re == nil {
        log.Println("Error during re")
        return
    }
    // rating count
    ratingCount := re.FindAllStringSubmatch(result, -1)
    pattern3 := `<span class="rating_num" property="v:average">(.*?)</span>`
    rp3 := regexp.MustCompile(pattern3)
    // score
    fScore := rp3.FindAllStringSubmatch(result, -1)
    pattern4 := `<img width="100" alt="(.*?)" src="(.*?)" class="">`
    // movie name
    rp4 := regexp.MustCompile(pattern4)
    fName := rp4.FindAllStringSubmatch(result, -1)

    StoreToFile(i, ratingCount, fScore, fName)

    page <- i
}

func StoreToFile(i int, ratingCount, fsScore, fName [][]string) {
    f, err := os.Create(strconv.Itoa(i) + ".txt")
    if err != nil {
        log.Println("Error during create file", err)
        return
    }
    defer f.Close()
    // 写入标题
    f.WriteString("电影名称 " + "评分" + "\t" + "评价人数" + "\t" + "\r\n")

    // 写内容
    n := len(ratingCount)
    fmt.Println("n=", n)
    for i := 0; i < n; i++ {
        f.WriteString(fName[i][1] + " " + fsScore[i][1] + "\t" + ratingCount[i][1] + "\t" + "\r\n")
    }
}

func DoWork(start, end int) {
    fmt.Printf("准备爬取第%d页到%d页的网址\n", start, end)
    page := make(chan int)
    for i := start; i <= end; i++ {
        // 定义一个函数，爬主页面
        go CrawlPages(i, page)
    }

    for i := start; i <= end; i++ {
        fmt.Printf("第%d个页面爬取完成\n", <-page)
    }

}
func main() {
    var start, end int
    fmt.Printf("请输入起始页( >= 1) :")
    fmt.Scan(&start)
    fmt.Printf("请输入终止页( >= 起始页) :")
    fmt.Scan(&end)

    DoWork(start, end) // 工作函数
}
