/*
 * @Description:
text/template和html/template2个模板标准库，
其中html/template为 HTML 提供了较为完整的支持。
包括普通变量渲染、列表渲染、对象渲染等。
本框架的模板渲染直接使用了html/template提供的能力。
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:29
 * @LastEditTime: 2022-03-20 21:06:59
 * @LastEditors: Please set LastEditors
*/
package my_web_frame

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc 用来定义路由映射的处理方法
type HandlerFunc func(*Context)

// Engine 一个实现了 ServeHTTP 的接口
type Engine struct {
	// 指针继承（嵌入）RouterGroup
	*RouterGroup                     // 根群组
	router        *router            // 路由器
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // 将所有的模板加载进内存
	funcMap       template.FuncMap   // 保存所有的自定义模板渲染函数
	/*type template.FuncMap map[string]any*/
}

// 分组
type RouterGroup struct {
	// 以组合的方式添加Engine
	engine      *Engine // all groups share a Engine instance
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting

}

// New 是Engine接口的构造函数
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
// 可以实现嵌套分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use is defined to add middleware to the group
// 中间件的定义与路由映射的 Handler 一致，处理的输入是Context对象。插入点是框架接收到请求初始化Context对象后，
// 允许用户使用自己定义的中间件做一些额外的处理，例如记录日志等，
// 以及对Context进行二次加工。
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute 调用router的方法，将路由函数注册到router
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 添加GET请求处理方法
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 添加POST请求处理方法
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// for custom render function
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// ServeHTTP 是必须实现的方法，
// 来自客户端的每个请求都会调用该方法进行进一步处理。它现在的功能是，
// 根据其所在的群组，按序添加相应的中间件和handler，
// 将请求和响应接口封装为上下文接口，并且调用router结构体实现路由解析，
// 找到注册在相应url的函数，进行调用处理。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	// 添加engine
	c.engine = engine
	engine.router.handle(c)
}

// Run 启动web server，是 ListenAndServe 的包装
func (engine *Engine) Run(addr string) {
	/*
		Handler interface {
			ServeHTTP(w ResponseWriter, r *Request)
		}
	*/
	// ListenAndServe 第二个参数接收的是一个Handler接口，
	// 需要实现ServerHTTP方法。
	// 传入一个实现了 ServeHTTP 接口的实例，所有的HTTP请求，
	// 都会交给该实例进行处理。
	fmt.Println("Server started on", addr)
	engine.router.roots["GET"].printAll()
	fmt.Println(engine.router.handlers)
	http.ListenAndServe(addr, engine)
}
