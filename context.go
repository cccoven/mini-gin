package mini_gin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 类型简写
type H map[string]any

// RequestContext 请求信息上下文
type RequestContext struct {
	Path   string            // 请求路径
	Method string            // 请求方法
	Params map[string]string // 存储 RESTful 风格路径参数
}

// ResponseContext 响应信息上下文
type ResponseContext struct {
	StatusCode int
}

// MiddlewareContext 中间件信息
type MiddlewareContext struct {
	handlers []HandleFunc // 存储处理函数列表（中间件与用户的处理函数）
	index    int          // 记录当前执行到第几个中间件
}

// Context 承载每一个接口的请求上下文（请求方法、url、参数等）
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	RequestContext
	ResponseContext
	MiddlewareContext
	engine *Engine // 上下文中持有 engine 成员变量，这样就可以通过上下文访问到 engine 中的 HTML 模板
}

// newContext 创建 context 实例
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		RequestContext: RequestContext{
			Path:   r.URL.Path,
			Method: r.Method,
		},
		MiddlewareContext: MiddlewareContext{
			index: -1,
		},
	}
}

// Next 应用下一个 HandleFunc
func (c *Context) Next() {
	// index 初始化为 -1，从头开始依次执行中间件与用户的处理函数
	c.index++
	s := len(c.handlers)

	for ; c.index < s; c.index++ {
		// 如果在第一个处理函数中（中间件中）调用了 Next() 方法，那么会将 c.index 再次 +1，然后再次进入循环，执行下一个处理函数
		// 执行完成返回到第一个处理函数时，c.index 已经被累加过了，会自动进入到与 c.index 对应的处理函数中或退出
		// 因为在一串处理链中，Context 上下文是一级一级传递下去的，它们共享同一个 index 变量
		c.handlers[c.index](c)
	}
}

// Param 获取 RESTful 风格路径参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Query 获取请求路径参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// PostForm 获取 POST 请求参数（form）
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// SetHeader 设置请求头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// Status 设置响应状态码
func (c *Context) Status(code int) {
	c.ResponseContext.StatusCode = code
	c.Writer.WriteHeader(code)
}

// String HTTP 字符串响应
func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON HTTP JSON 响应
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	// 将参数编码为 JSON 并写入到 Writer 响应
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Data HTTP 字节响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML HTTP html 文档响应
func (c *Context) HTML(code int, name string, data any) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	// c.Writer.Write([]byte(html))
	if err := c.engine.htmlTemplate.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	}
}

// Fail 响应失败
func (c *Context) Fail(code int, reason string) {
	c.Status(code)
	c.Writer.Write([]byte(reason))
}
