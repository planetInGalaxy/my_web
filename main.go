/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:41
 * @LastEditTime: 2022-03-20 20:53:56
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"net/http"

	"my_web_frame"
)

func main() {
	r := my_web_frame.Default()
	r.GET("/", func(c *my_web_frame.Context) {
		c.String(http.StatusOK, "Hello \n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *my_web_frame.Context) {
		names := []string{"Jake"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
