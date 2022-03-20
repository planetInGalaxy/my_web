/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-16 21:52:51
 * @LastEditTime: 2022-03-20 20:47:09
 * @LastEditors: Please set LastEditors
 */
package my_web_frame

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 起了一个别名 H，构建JSON数据时，显得更简洁。
type H map[string]interface{}

// 包含请求和响应等信息的上下文结构体
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// 提供对动态路由中的参数的访问
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	// engine pointer
	// 模板渲染中需要使用engine.htmlTemplates中注册的模板文件
	engine *Engine
}

// 上下文结构体构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1, // 记录执行流程
	}
}

// 按顺序执行index以后的所有handler函数或中间件
// c.handlers是这样的[A, B, Handler]，
// c.index初始化为-1，调用c.Next()，流程为：
// part1 -> part3 -> Handler -> part 4 -> part2
func (c *Context) Next() {
	// 每次调用Next，流程直接都向后推进
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

// 处理handlers失败的情况
func (c *Context) Fail(code int, err string) {
	// 直接跳过handlers中所有的中间件和handler
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// 获取url中路径参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 获取表单参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取URL查询参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 写入状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置响应头部字段
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 将字符串写入响应体
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 将接口/结构体转为JSON并写入响应体
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 将字节串（字节切片）写入响应体
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 将HTML文本写入响应体
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// 指定模板类型和数据进行渲染
// refer https://golang.org/pkg/html/template/
func (c *Context) Render(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	println("before")
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
	println("after")
}
