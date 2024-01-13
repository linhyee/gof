package pkg

import (
	"net/http"
	"sync"
)

// Pool is the context pool
type Pool struct {
	pool    *sync.Pool
	newFunc func() *Context
}

// creates a new context pool
func newPool(newFunc func() *Context) *Pool {
	p := &Pool{pool: &sync.Pool{}, newFunc: newFunc}
	p.pool.New = func() interface{} { return p.newFunc() }
	return p
}

// Acquire returns a Context from pool
func (p *Pool) Acquire(w http.ResponseWriter, req *http.Request) *Context {
	c := p.pool.Get().(*Context)
	c.beginRequest(w, req)
	return c
}

// Release puts a Context back to its pool, release its resources
func (p *Pool) Release(c *Context) {
	c.endRequest()
	p.pool.Put(c)
}
