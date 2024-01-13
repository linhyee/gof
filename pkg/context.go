package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type H map[string]interface{}

type Context struct {
	// origin object
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path, Method string
	Params       map[string]string
	// middleware
	handlers []HandlerFunc
	// response info
	StatusCode int
	index, id  int
	// engine pointer
	engine *Engine
}

// newContext new a context for http interface
func newContext(e *Engine) *Context {
	c := &Context{engine: e, index: -1, id: int(e.seqId())}
	return c
}

// String for output
func (c *Context) Info() string {
	return fmt.Sprintf("context{Path:%s,Method:%s,index:%d,id:%d}", c.Path, c.Method, c.index, c.id)
}

// BeginRequest
func (c *Context) beginRequest(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, g := range c.engine.groups {
		if strings.HasPrefix(req.URL.Path, g.prefix) {
			middlewares = append(middlewares, g.middlewares...)
		}
	}
	c.Req = req
	c.Writer = w
	c.Path = req.URL.Path
	c.Method = req.Method
	c.handlers = middlewares
}

// EndRequest
func (c *Context) endRequest() {
	c.Req = nil
	c.Writer = nil
	c.index = -1
}

// Next process control is passed to the next middleware
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Fail fail message response
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// Param get param vlue
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm get value from form
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query get url param's value
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// String string for plain text output
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON json ouput
func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data bytes output
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML html output
func (c *Context) RawHTML(code int, html string) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	c.Writer.Write([]byte(html))
}

// HTML template render
func (c *Context) HTML(code int, name string, data interface{}) {
	c.Writer.WriteHeader(code)
	c.Writer.Header().Set("Content-Type", "text/html")
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
