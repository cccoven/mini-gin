package mini_gin

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// trace 打印追踪栈
func trace(msg string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(msg + "\nTraceback: ")

	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s: %d", file, line))
	}

	return str.String()
}

func Recovery() HandleFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		}()

		// 如果执行 Next() 的过程中发生了 panic，那么就会进入 defer 的 recover() 流程
		c.Next()
	}
}
