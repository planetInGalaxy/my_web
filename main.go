/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:41
 * @LastEditTime: 2022-03-17 23:00:07
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"net/http"

	"my_web_frame"
)

func main() {
	r := my_web_frame.New()
	r.GET("/index", func(c *my_web_frame.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *my_web_frame.Context) {
			c.HTML(http.StatusOK, "<h1>Hello my_web_frame</h1>")
		})

		v1.GET("/hello", func(c *my_web_frame.Context) {
			// expect /hello?name=my_web_frame
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
		v11 := v1.Group("/v11")
		{
			v11.GET("/", func(c *my_web_frame.Context) {
				c.HTML(http.StatusOK, "<h1>Hello my_web_frame v11</h1>")
			})
		}
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *my_web_frame.Context) {
			// expect /hello/my_web_frame
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *my_web_frame.Context) {
			c.JSON(http.StatusOK, my_web_frame.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}
