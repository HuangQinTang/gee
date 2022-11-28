package gee

import (
	"net/http"
	"path"
)

// RouterGroup 路由组
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中间价集合
	parent      *RouterGroup  // 父级路由
	engine      *Engine       // 继承engine方法
}

// Group 创建一个路由组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	//添加到engine路由组集合
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	//log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.add(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use 加入中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) *RouterGroup {
	group.middlewares = append(group.middlewares, middlewares...)
	return group
}

// Static 创建静态资源服务
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root)) //http.Dir是实现来FileSystem接口到string类型
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	// http.StripPrefix主要是用于处理带路由前缀url的问题，比如前端请求的url路径是 /public/a.html
	// 经过了http.StripPrefix处理后，到达http.FileServer处理时的url就被截成了 /a.html,
	// 会在实际文件系统下的目录下( http.Dir("./static") 指定了是在./static相对目录下找) 寻找a.html文件。
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		//检查文件路径及是否有权限访问
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
