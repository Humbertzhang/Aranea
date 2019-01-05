package core

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const  (
	CLIENTINNERERROR = 1
)

type Crawler struct {
	URL 		string			`json:"url"`
	Method 		string			`json:"method"`
	ContentType	string 			`json:"content_type"`
	// "fotmat:key=value"
	Cookie 		string 			`json:"cookie"`
	Payload 	string			`json:"payload"`
	Headers 	http.Header 	`json:"headers"`
}

// 这里利用了http包中的Do()方法，省去了很多工作
func (crawler *Crawler) Request() (status int, body string, err error) {
	client := &http.Client{}
	request, err := http.NewRequest(crawler.Method, crawler.URL, strings.NewReader(crawler.Payload))
	if err != nil {
		return CLIENTINNERERROR, "", err
	}

	// 设置头部
	// 这里注意在接受任务的时候要如何把Header转换成http.Header形式
	if crawler.Headers != nil {
		request.Header = crawler.Headers
	}
	if crawler.Cookie != "" {
		request.Header.Set("Cookie", crawler.Cookie)
	}
	if crawler.ContentType != "" {
		request.Header.Set("Content-Type", crawler.ContentType)
	}


	// 发请求
	resp, err := client.Do(request)
	if err != nil {
		return resp.StatusCode, "", err
	}
	defer resp.Body.Close()

	// 读取返回
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CLIENTINNERERROR, "", err
	}

	return resp.StatusCode, string(respbody), nil
}