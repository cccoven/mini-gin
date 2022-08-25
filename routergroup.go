package mini_gin

import (
	"log"
	"net/http"
	"path"
)

// RouterGroup 路由分组
type RouterGroup struct {
	prefix      string       // 路由前缀
	middlewares []HandleFunc // 存储分组上的中间件
	parent      *RouterGroup // 支持嵌套分组
	engine      *Engine      // 具有 Engine 的所有能力
}

// createStaticHandler 创建静态资源处理器
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	// path.Join 将多个路径组成为一个路径
	absolutePath := path.Join(group.prefix, relativePath)
	// Go 内部提供了 http.FileServer 方法来实现一个静态资源服务器
	// http.FileServer 将收到的请求路径与 http.Dir 的路径拼接起来并传给 http.Dir 的 Open 方法打开对应的文件或目录进行处理
	// http.StripPrefix 方法会将请求路径中特定的前缀去掉，然后再进行处理
	// 如：我们请求的路径是 /public/test.js，http.StripPrefix 去掉前缀后的路径为 /test.js，再拼接上 http.Dir 的路径就是 ./public/test.js，最终访问到了文件
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	
	return func(c *Context) {
		file := c.Param("filepath")

		// 检查这个文件是否存在或者是否有权访问
		// 这个 fs 是有 http.Dir 实现的，这里打开文件是会拼接上 http.Dir 中的路径
		// eg：./public/test.js
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return 
		}
		
		// fileServer.ServeHTTP 会把 Req.URL.path 作为文件路径
		// 如果上面不去掉 absolutePath（/public）这个前缀的话，那么就会变成 http.Dir + Req.URL.path（/public/public/test.js)
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static 提供静态资源映射
// 仅仅是解析请求的地址，映射到服务器上文件的真实地址的工作交给 http.FileServer 处理就好了
// 静态资源请求的大致原理就是将请求路径作为文件路径去打开读取文件并将文件通过 http 媒体类型返回
func (group *RouterGroup) Static(relativePath, root string) {
	// http.Dir 是一个类型转换，底层是一个 string，表示文件的起始路径，空即为当前路径
	// http.Dir 实现了 FileSystem 接口，具有一个 Open() 方法，调用 Open() 方法时，传入的参数需要在前面拼接上该起始路径得到实际文件路径
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	// 这里借助我们自己实现的路由来处理路径就可以来，访问并返回文件的工作由 http.FileServer 处理
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

// Use 在路由分组上应用中间件
func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// Group 创建一个路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		// 为了支持嵌套分组，需要将父分组的前缀 + 当前分组的前缀，这样 trie 树和 router handlers 的哈希表中才有完整的路由名
		prefix: group.prefix + prefix, 
		// 存储父分组
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

/* 将路由的功能添加到 RouterGroup 上 */
// addRoute 添加一条路由映射
func (group *RouterGroup) addRoute(method string, comp string, handler HandleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandleFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandleFunc) {
	group.addRoute("POST", pattern, handler)
}
