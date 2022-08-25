package mini_gin

import (
	"net/http"
	"testing"
)

func TestRouterGroup(t *testing.T) {
	r := New()
	v1 := r.Group("/v1")
	v1.GET("/hello", func(c *Context) {
		c.String(200, "hello")
	})
	r.Run(":8080")
}

func TestNestedRouterGroup(t *testing.T) {
	r := New()
	
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong\n")
	})
	
	r.Run(":8080")
}

func TestStatic(t *testing.T) {
	r := New()
	r.Static("/public", "./public")
	r.Run(":8080")
}
