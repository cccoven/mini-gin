package mini_gin

import (
	"fmt"
	"net/http"
	"testing"
)

func TestEngine(t *testing.T) {
	r := New()
	
	r.GET("/ping", func(c *Context) {
		fmt.Fprintln(c.Writer, "pong")
	})
	
	r.Run(":8080")
}

func TestDefault(t *testing.T) {
	r := Default()
	r.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "pong\n")
	})
	r.Run(":8080")
}
