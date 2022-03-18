/*
 * @Description:
 我们实现分组控制，将具有相同前缀的路由划分为一组，以相同的方式来处理它们。
 中间件可以实现一些扩展方法，将不同的中间件应用到不同的分组上，就可以实现按需扩展。
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:29
 * @LastEditTime: 2022-03-17 22:28:20
 * @LastEditors: Please set LastEditors
*/
package my_web_frame

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 用来定义路由映射的处理方法
type HandlerFunc func(*Context)

// Engine 一个实现了 ServeHTTP 的接口
type Engine struct {
	// 指针继承（嵌入）RouterGroup
	*RouterGroup                // 根群组
	router       *router        // 路由器
	groups       []*RouterGroup // store all groups
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

// Use is defined to add middleware to the group
// 中间件的定义与路由映射的 Handler 一致，处理的输入是Context对象。插入点是框架接收到请求初始化Context对象后，
// 允许用户使用自己定义的中间件做一些额外的处理，例如记录日志等，
// 以及对Context进行二次加工。
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
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
