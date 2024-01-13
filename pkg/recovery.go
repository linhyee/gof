package pkg

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
// in trace (), runtime.Callers (3, pcs [:]) is called. Callers is used to return
// the program counter of the call stack. The 0th Caller is the Callers itself,
// the 1st is the previous trace, and the 2nd Is the defer func on the next level
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) //skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

// recovery
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}
