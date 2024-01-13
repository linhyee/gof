package pkg

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"sync/atomic"
)

// HandlerFunc defines the request handler used by pkg
type HandlerFunc func(*Context)

// Engine struct implement the interface of ServeHTTP
type (
	Engine struct {
		*RouterGroup  // abstract engine as group
		router        *router
		groups        []*RouterGroup     // store all groups
		htmlTemplates *template.Template // for html render
		funcMap       template.FuncMap   // for html render
		contextPool   *Pool
		seqId         func() int64
	}

	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		parent      *RouterGroup  // support nesting
		engine      *Engine       // all groups share a Engine instance
	}
)

// Sequence number generator
func SequenceId() func() int64 {
	var id int64
	return func() int64 {
		return atomic.AddInt64(&id, 1)
	}
}

// New function is the constructor of pkg.Engine
func New() *Engine {
	e := &Engine{router: newRouter(), seqId: SequenceId()}
	e.RouterGroup = &RouterGroup{engine: e}
	e.groups = []*RouterGroup{e.RouterGroup}

	e.contextPool = newPool(func() *Context {
		return newContext(e)
	})
	return e
}

// Default use Logger() & Recovery middlewares
func Default() *Engine {
	e := New()
	e.Use(Logger(), Recovery())
	return e
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	e := g.engine
	ng := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: e,
	}
	e.groups = append(e.groups, ng)
	return ng
}

// Use is defined to add middleware to the group
func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// addRoute add a new route
func (g *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

// create a static handler
func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (g *RouterGroup) Static(relativePath, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handler
	g.GET(urlPattern, handler)
}

// for custom render function
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

// load html glob
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

// Run defiens the method to start a http server
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := e.contextPool.Acquire(w, req)
	e.router.handle(c)
	e.contextPool.Release(c)
}
