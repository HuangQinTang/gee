package gee

import (
	"fmt"
	"html/template"
	"net/http"
	"testing"
	"time"
)

func TestGee(t *testing.T) {
	//path := "/test"
	//addr := "127.0.0.1:8080"
	//content := "hello world"

	//r := New()
	//r.GET(path, func(c *Context) {
	//	c.Data(http.StatusOK, []byte(content))
	//})

	//go r.Run(addr)

	//time.Sleep(time.Second * 1)
	//res, err := http.Get("http://" + addr + path)
	//if err != nil {
	//	t.Errorf(err.Error())
	//}
	//defer res.Body.Close()

	//body, _ := ioutil.ReadAll(res.Body)
	//if content != string(body) {
	//	fmt.Println("预期", content, "---", "响应", string(body))
	//	t.Errorf("请求失败，返回值不符预期")
	//}

	addr := "127.0.0.1:8080"
	r := New()
	r.GET("/", func(c *Context) {
		c.Data(http.StatusOK, []byte("Hello Gee"))
	})
	r.GET("/hello", func(c *Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	// 测试参数绑定
	r.GET("/hello/:name", func(c *Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	// 测试模糊匹配
	r.GET("/assets/*filepath", func(c *Context) {
		c.JSON(http.StatusOK, H{"filepath": c.Param("filepath")})
	})

	// 测试分组路由
	v1 := r.Group("/v1").Use(testMiddleware)
	{
		v1.GET("/test", func(c *Context) {
			fmt.Println("/v1/test")
			c.Data(http.StatusOK, []byte("/v1/test"))
		})
	}

	// 测试静态资源服务
	r.Static("/assets", "./static") //相对路径
	//r.Static("/assets", "/usr/geektutu/blog/static")	//绝对路径

	// 测试模版功能
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	//r.Static("/assets", "./static") //上面已设置
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/tmpl", func(c *Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.GET("/students", func(c *Context) {
		c.HTML(http.StatusOK, "arr.tmpl", H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	//测试异常捕获
	r.GET("/panic", func(c *Context) {
		panic("出错来了")
		c.Data(http.StatusOK, []byte(""))
	})
	r.Run(addr)

	// 测试用例
	// curl "http://localhost:8080/hello/geektutu"
	// curl "http://localhost:8080/assets/css/geektutu.css"
	// curl "http://localhost:8080/v1/test"
	// curl http://127.0.0.1:8080/assets/test.js
	// curl http://127.0.0.1:8080/tmpl
	// curl http://127.0.0.1:8080/students
	// curl http://127.0.0.1:8080/panic
}

func testMiddleware(c *Context) {
	fmt.Println("前置逻辑")
	c.Next()
	fmt.Println("后置逻辑")
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
