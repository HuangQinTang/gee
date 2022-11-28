package gee

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGee(t *testing.T) {
	path := "/test"
	addr := "127.0.0.1:8080"
	content := "hello world"

	r := New()
	r.GET(path, func(c *Context) {
		c.Data(http.StatusOK, []byte(content))
	})

	go r.Run(addr)

	time.Sleep(time.Second * 1)
	res, err := http.Get("http://" + addr + path)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if content != string(body) {
		fmt.Println("预期", content, "---", "响应", string(body))
		t.Errorf("请求失败，返回值不符预期")
	}

	//r:= New()
	//r.GET("/", func(c *Context) {
	//	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	//})
	//
	//r.GET("/hello", func(c *Context) {
	//	// expect /hello?name=geektutu
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	//})
	//
	//r.GET("/hello/:name", func(c *Context) {
	//	// expect /hello/geektutu
	//	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	//})
	//
	//r.GET("/assets/*filepath", func(c *Context) {
	//	c.JSON(http.StatusOK, H{"filepath": c.Param("filepath")})
	//})
	//
	//r.Run(addr)
	// 测试用例
	//curl "http://localhost:8080/hello/geektutu"
	//curl "http://localhost:8080/assets/css/geektutu.css"
}
