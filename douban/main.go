package main

import (
	"fmt"
	"github.com/goquery"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type movie struct {
	//排名
	rank int
	name string
}

var movielist = make([]movie, 251)

//fetch 取来
func fetch(url string) *goquery.Document {
	fmt.Println(url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP获取错误", err)
		return nil
	}
	if resp.StatusCode != 200 {
		fmt.Println("返回状态码错误:", resp.StatusCode)
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("读取失败", err)
		return nil
	}
	return doc
}
func parseUrls(url string) {
	doc := fetch(url)
	var indexNum int
	tmp := movie{}
	//pic节点
	doc.Find("ol.grid_view li").Find(".pic").Each(
		func(index int, ele *goquery.Selection) {
			b := ele.Find("em")
			//fmt.Println(b.Text())
			tmp.rank, _ = strconv.Atoi(b.Text())
			a, _ := ele.Find("img").Attr("alt")
			//fmt.Println(a)
			tmp.name = a
			movielist[tmp.rank] = tmp
			indexNum = tmp.rank
		})
	time.Sleep(time.Second * 1)
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			parseUrls("https://movie.douban.com/top250?start=" + strconv.Itoa(25*i))
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("历时:%s\n", elapsed)
	for i := 1; i <= 250; i++ {
		fmt.Printf("第%d名:%s\n", movielist[i].rank, movielist[i].name)
	}
}
