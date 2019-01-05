package core

import (
	"net/http"
	"strings"
	"testing"
)

const (
	GetURL = "http://www.baidu.com"
	PostURL = "https://postman-echo.com/post"
	PutURL = "https://postman-echo.com/put"
	DeleteURL = "https://postman-echo.com/delete"
)

func TestCrawler_Request_GET(t *testing.T) {
	crawler := &Crawler{}
	crawler.URL = GetURL
	crawler.Method = HTTPGET

	statusCode, _, err := crawler.Request()
	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusOK {
		t.Error("StatusCode is not StatusOK:", statusCode)
	}
}

func TestCrawler_Request_POST(t *testing.T) {
	crawler := &Crawler{}
	crawler.URL = PostURL
	crawler.Method = HTTPPOST
	crawler.Payload = "TestString"
	// 注意ContentType
	crawler.ContentType = "text/plain"
	statusCode, body, err := crawler.Request()

	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusOK {
		t.Error("StatusCode is not StatusOK:", statusCode)
	}
	if !strings.Contains(body, crawler.Payload) {
		t.Error("crawler.Payload Not in echo content.")
	}
}

func TestCrawler_Request_PUT(t *testing.T) {
	crawler := &Crawler{}
	crawler.URL = PutURL
	crawler.Method = HTTPPUT
	crawler.Payload = "TestString"
	crawler.ContentType = "text/plain"
	statusCode, body, err := crawler.Request()

	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusOK {
		t.Error("StatusCode is not StatusOK:", statusCode)
	}
	if !strings.Contains(body, crawler.Payload) {
		t.Error("crawler.Payload Not in echo content.")
	}
}

func TestCrawler_Request_DELETE(t *testing.T) {
	crawler := &Crawler{}
	crawler.URL = DeleteURL
	crawler.Method = HTTPDELETE
	crawler.Payload = "TestString"
	crawler.ContentType = "text/plain"
	statusCode, body, err := crawler.Request()

	if err != nil {
		t.Error(err)
	}
	if statusCode != http.StatusOK {
		t.Error("StatusCode is not StatusOK:", statusCode)
	}
	if !strings.Contains(body, crawler.Payload) {
		t.Error("crawler.Payload Not in echo content.")
	}
}