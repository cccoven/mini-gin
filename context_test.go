package mini_gin

import (
	"net/http"
	"testing"
)

func TestContext(t *testing.T) {
	r := New()

	r.GET("/ping", func(c *Context) {
		c.String(http.StatusOK, "res: %s\n", "pong")
	})

	r.GET("/say", func(c *Context) {
		msg := c.Query("msg")
		c.String(http.StatusOK, msg+"\n")
	})

	r.POST("/login", func(c *Context) {
		phone := c.PostForm("phone")
		password := c.PostForm("password")
		c.JSON(http.StatusOK, H{
			"phone":    phone,
			"password": password,
		})
	})

	r.Run(":8080")
}

func TestTemplate(t *testing.T) {
	r := New()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *Context) {
		c.HTML(http.StatusOK, "test.tmpl", nil)
	})
	r.Run(":8080")
}
