package main

import (
    "io"
    "log"
    "net/http"
    "os"
    "regexp"
    "strconv"
)

func GetHTTPUrl(url string) (urlString string, err error) {
    resp, err1 := http.Get(url)
    if err1 != nil {
        err = err1
        return
    }
    defer resp.Body.Close()

    buf := make([]byte, 4096)
    for {
        n, err2 := resp.Body.Read(buf)
        if n == 0 {
            break
        }
        if err2 != nil && err2 != io.EOF {
            err = err2
            return
        }
        urlString += string(buf[:n])
    }
    return
}

func SavePic(index int, url string, page chan int) {
    path := "/Users/suzl/dev/golang/go-demos/crawler" + strconv.Itoa(index+1) + ".jpg"
    f, err := os.Create(path)
    if err != nil {
        log.Printf("Error during create file:%v", err)
        return
    }
    defer f.Close()

    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error during http.Get :%v", err)
        return
    }
    defer resp.Body.Close()
    buf := make([]byte, 4096)
    for {
        n, err2 := resp.Body.Read(buf)
        if n == 0 {
            log.Println("Empty resp.Body")
            break
        }
        if err2 != nil && err2 != io.EOF {
            err = err2
            return
        }
        f.Write(buf[:n])
    }
    page <- index
}

func main() {
    url := "https://www.douyu.com/g_yz"

    // get page and save the page to result
    result, err := GetHTTPUrl(url)
    if err != nil {
        log.Printf("Error during http.Get in main %v", err)
        return
    }
    // get all page an save to the result
    ret := regexp.MustCompile(`data-original="(?s:(.*?))"`)
    // compile through the regEx
    alls := ret.FindAllStringSubmatch(result, -1)
    page := make(chan int)
    n := len(alls)
    for idx, imgURL := range alls {
        go SavePic(idx, imgURL[1], page)
    }

    for i := 0; i < n; i++ {
        log.Printf("%d img is compleat\n", <-page)
    }
}
