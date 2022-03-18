// 中间件(middlewares)，简单说，就是非业务的技术类组件。
// Web 框架本身不可能去理解所有的业务，因而不可能实现所有的功能。
// 因此，框架需要有一个插口，允许用户自己定义功能，嵌入到框架中，仿佛这个功能是框架原生支持的一样。

package my_web_frame

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// 等待执行后续所有的中间件或用户的Handler完毕
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
