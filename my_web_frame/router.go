/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-16 22:03:51
 * @LastEditTime: 2022-03-17 22:37:58
 * @LastEditors: Please set LastEditors
 */
package my_web_frame

import (
	"net/http"
	"strings"
)

// 定义路由器结构体
type router struct {
	// 使用 roots 来存储每种请求方式(GET、POST...)的 Trie 树根节点
	// 主要用来解析url，将其分段，获取动态url参数
	roots map[string]*node
	// 使用 handlers 存储指定请求方式和URL的 HandlerFunc
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

// 路由器构造方法
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析URL字符串为由/分段的切片，并且到*为止（忽略*后所有字符串）
// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 将路由映射的处理方法注册到router
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	// 插入到Trie树中
	r.roots[method].insert(pattern, parts, 0)
	// 记录特定请求方法+特定类型url 到 handler 的映射
	r.handlers[key] = handler
}

// 解析用户请求的URL并获取特定的URL模板和参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	// 获取路由函数的过程中用一个map记录参数
	params := make(map[string]string)
	// 根据不同的方法返回不同的前缀树
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}
	// 在前缀树中查找符合URL结构的节点
	node := root.search(searchParts, 0)

	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

// 路由控制，根据url和参数，执行特定的处理程序
func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
