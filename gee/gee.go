package gee

import (
	"net/http"
)

type H map[string]interface{}
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix   string
	midwares []HandlerFunc
	parents  *RouterGroup
	engine   *Engine
}

type Engine struct {
	*RouterGroup
	routers *router
	groups  []*RouterGroup
}

func New() *Engine {
	engine := &Engine{routers: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (rg *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	newPattern := rg.prefix + pattern
	rg.engine.routers.addRoute(method, newPattern, handler)
}

func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
	rg.addRoute("GET", pattern, handler)
}

func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
	rg.addRoute("POST", pattern, handler)
}

func (engine *Engine) RUN(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewContext(w, r)
	engine.routers.handle(context)
}

func (rg *RouterGroup) Group(name string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix:  rg.prefix + name,
		engine:  engine,
		parents: rg,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}
