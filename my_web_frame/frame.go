/*
 * @Description:
 Go语言内置了net/http库，封装了HTTP网络编程的基础的接口，
 本 Web 框架便是基于 net/http 实现的。
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:29
 * @LastEditTime: 2022-03-15 22:25:19
 * @LastEditors: Please set LastEditors
*/
package my_web_frame

import (
	"fmt"
	"net/http"
)

// HandlerFunc 用来定义路由映射的处理方法
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 一个实现了 ServeHTTP 的接口
type Engine struct {
	// 添加了一张路由映射表router，用来注册处理方法
	// key 由请求方法和静态路由地址构成，value是处理方法
	router map[string]HandlerFunc
}

// New 是Engine接口的构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRoute 将路由映射的处理方法注册到router
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET 添加GET请求处理方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 添加POST请求处理方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// ServeHTTP 是必须实现的方法，它现在的功能是，解析请求的路径，
// 并查找路由映射表，如果查到，就执行注册的处理方法；如果查不到，就返回 404
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
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
	http.ListenAndServe(addr, engine)

}
