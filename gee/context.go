package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Path    string
	Method  string
	Params  map[string]string
	Status  int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
	}
}

// 查询表单
func (context *Context) GetForm(key string) string {
	return context.Request.FormValue(key)
}

// 查询参数
func (context *Context) GetQuery(key string) string {
	return context.Request.URL.Query().Get(key)
}

// 查询路由参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 设置状态码
func (context *Context) SetStatus(num int) {
	context.Status = num
	context.Writer.WriteHeader(num)
}

// 设置响应头字段
func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Set(key, value)
}

// 返回响应字符串
func (context *Context) String(code int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.SetStatus(code)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回响应JSON
func (context *Context) Json(code int, obj interface{}) {
	context.SetHeader("Context-Type", "application/json")
	context.SetStatus(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Writer, err.Error(), 500)
	}
}

// 返回响应数组
func (context *Context) Data(code int, data []byte) {
	context.SetStatus(code)
	context.Writer.Write(data)
}

// 返回HTML文件
func (context *Context) HTML(code int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.SetStatus(code)
	context.Writer.Write([]byte(html))
}
