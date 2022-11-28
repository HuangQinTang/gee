package gee

import (
	"net/http"
	"testing"
)

func TestRouterGroup(t *testing.T) {
	r := New()

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})
		v1.GET("/test", func(c *Context) {
			c.Data(http.StatusOK, []byte("v1 test"))
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/test", func(c *Context) {
			c.Data(http.StatusOK, []byte("b2 test"))
		})

		v2New := v2.Group("/new")
		{
			v2New.GET("/test", func(c *Context) {
				c.Data(http.StatusOK, []byte("v2/new/test test"))
			})
		}
	}

	r.Run("127.0.0.1:8080")

	// 测试用例
	// curl http://127.0.0.1:8080/v1/
	// curl http://127.0.0.1:8080/v1/test
	// curl http://127.0.0.1:8080/v2/test
	// curl http://127.0.0.1:8080/v2/new/test

}
