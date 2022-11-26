package gee

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGee(t *testing.T) {
	path := "/test"
	addr := "127.0.0.1:8080"
	content := "hello world"

	server := New()
	server.GET(path, func(c *Context) {
		c.Data(http.StatusOK, []byte(content))
	})

	go server.Run(addr)

	time.Sleep(time.Second * 1)
	res, err := http.Get("http://" + addr + path)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if content != string(body) {
		t.Errorf("请求失败，返回值不符预期")
	}
}
