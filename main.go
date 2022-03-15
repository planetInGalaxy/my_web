/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:41
 * @LastEditTime: 2022-03-15 22:18:30
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"fmt"
	"net/http"

	"my_web_frame"
)

func main() {
	r := my_web_frame.New()

	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		// 打印URL路径
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		// 打印所有请求头字段
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":9999")
}
