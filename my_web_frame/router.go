/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-16 22:03:51
 * @LastEditTime: 2022-03-16 22:07:37
 * @LastEditors: Please set LastEditors
 */
package my_web_frame

import (
	"log"
	"net/http"
)

// 定义路由器结构体
type router struct {
	handlers map[string]HandlerFunc
}

// 路由器构造方法
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// 将路由映射的处理方法注册到router
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 解析请求的路径，并查找路由映射表，如果查到，就执行注册的处理方法；
// 如果查不到，就返回 404
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
