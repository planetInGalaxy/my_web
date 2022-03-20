/*
 * @Description:
 * @Author: Tjg
 * @Date: 2022-03-15 21:36:41
 * @LastEditTime: 2022-03-20 20:53:56
 * @LastEditors: Please set LastEditors
 */
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"my_web_frame"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := my_web_frame.New()
	r.Use(my_web_frame.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Jimmy", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *my_web_frame.Context) {
		c.Render(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *my_web_frame.Context) {
		c.Render(http.StatusOK, "arr.tmpl", my_web_frame.H{
			"title":  "hello",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *my_web_frame.Context) {
		c.Render(http.StatusOK, "custom_func.tmpl", my_web_frame.H{
			"title": "hello",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
