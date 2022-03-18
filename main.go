/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:41
 * @LastEditTime: 2022-03-17 23:00:07
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"log"
	"net/http"
	"time"

	"my_web_frame"
)

// 用户自定义的中间件函数
func onlyForV2() my_web_frame.HandlerFunc {
	return func(c *my_web_frame.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {

	r := my_web_frame.New()

	r.Use(my_web_frame.Logger()) // global midlleware
	r.GET("/", func(c *my_web_frame.Context) {
		c.HTML(http.StatusOK, "<h1>Hello my_web_frame</h1>")
	})

	v2 := r.Group("/v2")

	v2.Use(onlyForV2()) // v2 group middleware
	v2.GET("/hello/:name", func(c *my_web_frame.Context) {
		// expect /hello/my_web_frame
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.Run(":9999")
}
