package mini_gin

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestRouter(t *testing.T) {
	r := newRouter()
	// r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	r.addRoute("GET", "/", nil)

	n, params := r.getRoute("GET", "/")
	n, params = r.getRoute("GET", "/hello/mini-gin")
	n, params = r.getRoute("GET", "/assets/public/xx")
	fmt.Println(n, params)
}

func TestParam(t *testing.T) {
	r := New()

	r.GET("/a/b", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/a/b", func(c *Context) {
		c.String(http.StatusOK, "okok")
	})

	r.GET("/hello/:name", func(c *Context) {
		c.String(http.StatusOK, c.Param("name"))
	})

	r.GET("/hello/:name/aaa", func(c *Context) {
		c.String(http.StatusOK, c.Param("name")+"aaa")
	})

	r.GET("/assets/*filepath/test", func(c *Context) {
		c.String(http.StatusOK, c.Param("filepath"))
	})

	r.Run(":8080")
}
