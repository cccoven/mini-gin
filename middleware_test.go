package mini_gin

import (
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	r := New()
	v1 := r.Group("/v1")
	v1.Use(Logger())
	v1.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "Hello Middleware")
	})
	r.Run(":8080")
}
