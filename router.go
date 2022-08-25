package mini_gin

import (
	"net/http"
	"strings"
)

type router struct {
	// 使用一个 map 来保存每个请求方法对应的路径 trie 树
	// roots key eg: roots['GET'] roots['POST']
	roots map[string]*node
	// 另外使用 map 来保存每个请求（方法加路径）对应的处理函数
	// handlers key eg: handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
	handlers map[string]HandleFunc // 维护所有路由与处理函数的映射
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandleFunc),
	}
}

// parsePattern 将路径解析成数组，数组按下标顺序形成父子树节点
// eg：/p/:lang/doc -> ["", "p", ":lang", "doc"]    /p/*/doc -> ["", "p", "*"]
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 如果路径中带有 * 号，表示 * 号及之后的位置可以匹配上任何值
			// 所以只要将 * 号作为最后一个叶子节点即可
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

// addRoute 向映射中添加一个条目
func (r *router) addRoute(method string, pattern string, handler HandleFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	// 每个请求方法（GET/POST）就是一个 trie 树的根结点
	// eg: r.roots[GET] -> trie    r.roots[POST] -> trie
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute 匹配路由并设置动态参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchPath := parsePattern(path)
	// 存储动态路由的变量与匹配上的路由的部分的映射
	params := make(map[string]string)

	// 取出请求方法对应的 trie 树
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// 从 trie 中查出匹配到的节点
	// eg: /hello/mini-gin 匹配到 /hello/:name 这个节点
	n := root.search(searchPath, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		// 将匹配到的节点遍历
		for idx, part := range parts {
			if part[0] == ':' {
				// 将 : 匹配到的字符串存入 params 中
				// eg: {name: mini-gin}
				params[part[1:]] = searchPath[idx]
			}
			if part[0] == '*' && len(part) > 1 {
				// 如果路径中只有一个 *，那么就不存在动态变量
				// 如果 * 后面存在一个变量，则将该变量后面的所有路径都变成这个变量的值
				// eg: 匹配规则为：/assets/*filepath，传入的路径为 /assets/public/xx，那么 params 则存储 {filepath: public/xx}
				params[part[1:]] = strings.Join(searchPath[idx:], "/")
			}
		}
		
		// 最后将节点与 params 映射表返回，并将 params 存入 context 中以便使用 ctx.Param(key) 来获取动态路由
		return n, params
	}

	return nil, nil
}

// handle 从映射表中取出条目并执行用户的处理函数
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 匹配到了路由，将路由参数存入给上下文并处理请求
		c.Params = params
		key := c.Method + "-" + n.pattern
		// 直接从路由中取出处理函数并执行
		// r.handlers[key](c)
		
		// 查找到路由节点后，从路由中取出对应的处理函添加到上下文的处理函数链中
		// 而在前面的 ServeHTTP 中，已经先将中间件加入进去了，所以用户的处理函数在所有中间件的后面
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 Not Found - %s\n", c.Path)
		})
	}
	
	// 直接由上下文的 Next() 函数开始执行处理函数链
	c.Next()
}
