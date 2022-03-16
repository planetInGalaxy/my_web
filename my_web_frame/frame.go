/*
 * @Description:
 Go语言内置了net/http库，封装了HTTP网络编程的基础的接口，
 本 Web 框架便是基于 net/http 实现的。
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:29
 * @LastEditTime: 2022-03-16 22:24:51
 * @LastEditors: Please set LastEditors
*/
package my_web_frame

import (
	"fmt"
	"net/http"
)

// HandlerFunc 用来定义路由映射的处理方法
type HandlerFunc func(*Context)

// Engine 一个实现了 ServeHTTP 的接口
type Engine struct {
	// 添加路由器
	router *router
}

// New 是Engine接口的构造函数
func New() *Engine {
	return &Engine{router: newRouter()}
}

// addRoute 将路由映射的处理方法注册到router
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET 添加GET请求处理方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 添加POST请求处理方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// ServeHTTP 是必须实现的方法，
// 来自客户端的每个请求都会调用该方法进行进一步处理。它现在的功能是，
// 将请求和响应接口封装为上下文接口，并且调用router结构体实现路由解析，
// 找到注册在相应url的函数，进行调用处理。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
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
	http.ListenAndServe(addr, engine)
}
