package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

// HTTPGetDB get the certain url page, and get the result
func HTTPGetDB(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)

	// for loop to get all the whole page data
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			reak
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		result += string(buf[:n])
	}
	return
}
func Save2File(idx int, filmName, filmScore, peopleNum [][]string) {
	path := "/Users/suzl/dev/golang/go-demos/crawler" + "No-" + strconv.Itoa(idx) + ".txt"
	f, err := os.Create(path)
	if err != nil {
		log.Printf("Error during os.Create file: %v", err)
		return
	}
	defer f.Close()

	_, _ = f.WriteString("FileName" + "\t\t\t" + "Score" + "\t\t\t" + "RatingCount" + "\n")

	for i := 0; i < len(filmName); i++ {
		f.WriteString(filmName[i][1] + "\t\t\t" + filmScore[i][1] + "\t\t" + peopleNum[i][1] + "\n")
	}
}

func SpiderPage(idx int, page chan int) {
	url := "https://movie.douban.com/top250?start=" + strconv.Itoa((idx-1)*25) + "&filter="
	// get url page
	result, err := HTTPGetDB(url)
	if err != nil {
		log.Printf("Error during HTTPGetDB:%v", err)
		return
	}
	// get the film name by regEx
	ret1 := regexp.MustCompile(`<img width="100" alt="(?s:(.*?))"`)
	filmName := ret1.FindAllStringSubmatch(result, -1)

	// get score
	pattern := `<span class="rating_num" property="v:average">(?s:(.*?))</span>`
	ret2 := regexp.MustCompile(pattern)
	filmScore := ret2.FindAllStringSubmatch(result, -1)
	// get rating numbers
	ret3 := regexp.MustCompile(`<span>(?s:(\d*?))人评价</span>`)
	peopleNum := ret3.FindAllStringSubmatch(result, -1)

	Save2File(idx, filmName, filmScore, peopleNum)
	// work with the goroutine, to sync
	page <- idx
}

func startWork(start, end int) {
	log.Printf("Crawling page from %d to %d \n", start, end)
	page := make(chan int) // to avoid the goroutine ends
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		log.Printf("Page finished %d \n", <-page)
	}
}

func main() {
	var start, end int
	log.Println("Please input the start page (>=1):")
	_, _ = fmt.Scan(&start)

	log.Println("Please input the end page (>=start):")
	_, _ = fmt.Scan(&end)
	startWork(start, end)
}
