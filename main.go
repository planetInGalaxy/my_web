/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:41
 * @LastEditTime: 2022-03-16 22:25:16
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"net/http"

	"my_web_frame"
)

func main() {
	r := my_web_frame.New()

	r.GET("/", func(c *my_web_frame.Context) {
		// 打印URL路径
		c.HTML(http.StatusOK, "<h1>Hello World!</h1>")
	})

	r.GET("/hello", func(c *my_web_frame.Context) {
		// expect /hello?name=xxx
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *my_web_frame.Context) {
		c.JSON(http.StatusOK, my_web_frame.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run("localhost:9999")
}
