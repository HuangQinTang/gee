package gee

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup  //继承路由组方法
	router        *router
	groups        []*RouterGroup     // 当前已有路由组集合
	htmlTemplates *template.Template // 保存模版
	funcMap       template.FuncMap   // 自定义模板渲染函数
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// addRoute 添加路由
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.add(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.add("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.router.add("POST", pattern, handler)
}

// SetFuncMap 设置模版方法
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 加载模版
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// ServeHTTP 实现http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// 添加默认中间件
	middlewares = append(middlewares, Recovery(), Logger())
	for _, group := range engine.groups {
		// 收集当前请求路径需要执行到中间件
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	// 请求响应对象封装到Context
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Logger 中间件 打印请求状态、路径、市场
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

// Recovery 中间件 捕获异常
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

// trace 获取堆栈信息
func trace(message string) string {
	var pcs [32]uintptr
	// 第 0 个 Caller 是 Callers 本身， 第 1 个是上一层 trace，
	// 第 2 个是再上一层的 defer func。因此，为了日志简洁一点，这里我们跳过了前 3 个 Caller。
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
