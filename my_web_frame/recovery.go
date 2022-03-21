/*
对于一个web框架而言，错误处理机制是必要的。
web服务在运行过程中可能会出现异常，进而导致系统宕机。
net/http中已经实现了recovery，可以保证服务不会挂掉。
Recovery 中间件在这里的作用是保证所有请求都能正常响应，
否则，panic 之后，就没回应了。
原理：
panic 会导致程序被中止，可以手动触发，也会在程序运行过程中遇见错误时自动触发，
会导致程序退出，但是在退出前，会先处理完当前协程上已经 defer 的任务，
执行完成后再退出。Go 语言还提供了 recover 函数，
可以避免因为 panic 发生而导致整个程序终止，recover 函数只在 defer 中生效。
*/

package my_web_frame

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(0, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		// 再执行其他中间件和handler之前，使用defer和recover，
		// 保证错误发生时能够得到程序的控制权
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				// 返回500
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
