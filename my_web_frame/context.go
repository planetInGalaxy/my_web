/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-16 21:52:51
 * @LastEditTime: 2022-03-16 22:02:55
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
	// response info
	StatusCode int
}

// 上下文结构体构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 获取表单参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取URL参数
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
