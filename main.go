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
	r.GET("/", func(c *my_web_frame.Context) {
		c.HTML(http.StatusOK, "<h1>Hello my_web_frame</h1>")
	})

	r.GET("/hello", func(c *my_web_frame.Context) {
		// expect /hello?name=planet
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name/123", func(c *my_web_frame.Context) {
		// expect /hello/planet
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *my_web_frame.Context) {
		c.JSON(http.StatusOK, my_web_frame.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
