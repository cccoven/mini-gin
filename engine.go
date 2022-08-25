package mini_gin

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

// HandleFunc 请求处理函数（Controller）
type HandleFunc func(*Context)

// Engine 框架总引擎
type Engine struct {
	*RouterGroup                // 将 Engine 作为顶层分组，并具有 RouterGroup 的所有能力
	router       *router        // 存储路由
	groups       []*RouterGroup // 存储所有分组
	// Go 内置了 text/template 和 html/template 2 个模板标准库，其中 html/template 为 HTML 提供了较为完整的支持。包括普通变量渲染、列表渲染、对象渲染等
	htmlTemplate *template.Template // 渲染 HTML 使用
	funcMap      template.FuncMap   // 渲染 HTML 使用
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

// Default 默认应用 Logger 和 Recovery 中间件
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplate = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

// Run 启动 HTTP 服务器
func (e *Engine) Run(addr string) error {
	log.Printf("Server is running at %s\n", addr)
	// ListenAndServe 第二个参数传入一个实现了 ServeHTTP 方法的结构体，即可接管 http 标准库的请求
	return http.ListenAndServe(addr, e)
}

// ServeHTTP 在引擎上实现 go 内部的 Handler 接口，接管 http 标准库的请求
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandleFunc
	for _, group := range e.groups {
		// 当请求路径前缀命中分组前缀时，将该分组上应用的中间件全部取出
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, r)
	// 在收到请求时，先将这个分组上的所有中间件加入上下文的处理函数列表（handlers）中
	c.handlers = middlewares
	c.engine = e
	e.router.handle(c)
}
