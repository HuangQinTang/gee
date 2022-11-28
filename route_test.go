package gee

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.add("GET", "/", nil)
	r.add("GET", "/hello/:name", nil)
	r.add("GET", "/hello/b/c", nil)
	r.add("GET", "/hi/:name", nil)
	r.add("GET", "/assets/*filepath", nil)
	r.add("GET", "/test/:name/haha", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

// Read 打印路由树节点信息
func Read(node []*node) {
	for _, v := range node {
		fmt.Printf("%#v\n", v)
		Read(v.children)
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	//for _, v := range r.roots {
	//	Read(v.children)
	//}

	//获取匹配到到叶子节点与路由绑定的参数
	n, ps := r.getRoute("GET", "/hello/geektutu")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	//测试路由是否匹配到
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	//测试参数是否绑定成功
	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}
	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

	//测试绑定函数是否正确
	r.add("get", "/18", eighteen)
	r.add("get", "/:num", num)
	if GetFunctionName(r.handlers["get-/:num"]) != "gee.num" {
		t.Fatal("映射方法不对")
	}
	if GetFunctionName(r.handlers["get-/18"]) != "gee.eighteen" {
		t.Fatal("映射方法不对")
	}

	//测试路由冲突是否报错
	defer func() {
		if fail := recover(); fail != nil {
			if fail != "同级路由冲突" {
				t.Errorf("路由冲突检测失败")
			}
		}
	}()
	r.add("get", "/19", num)
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func num(c *Context) {
	c.Data(http.StatusOK, []byte("num fun"))
}

func eighteen(c *Context) {
	c.Data(http.StatusOK, []byte("eighteen"))
}
