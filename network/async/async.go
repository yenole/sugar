package async

import (
	"sync"
)

type chanx struct {
	C chan func()
}

func (c *chanx) proc(p *sync.Pool) {
g:
	if fn := <-c.C; fn != nil {
		fn()
		p.Put(c)
		goto g
	}
}

type Async struct {
	p *sync.Pool
}

func New() *Async {
	p := &Async{}
	p.p = &sync.Pool{New: p.chanx}
	return p
}

func (a *Async) chanx() interface{} {
	ch := &chanx{C: make(chan func(), 1)}
	go ch.proc(a.p)
	return ch
}

func (a *Async) Do(fn func()) {
	if ch, ok := a.p.Get().(*chanx); ok {
		ch.C <- fn
	}
}
