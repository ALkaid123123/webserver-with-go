package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*trie       //不同的方法对应不同的前缀树
	handlers map[string]HandlerFunc //不同方法+url对应不同回调函数
}

func NewRouter() *router {
	return &router{roots: make(map[string]*trie), handlers: make(map[string]HandlerFunc)}
}

// 将请求URL根据/解析为字符串切片
func parseURL(pattern string) []string {
	parts := strings.Split(pattern, "/")
	res := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			res = append(res, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return res
}

// 配置静态路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &trie{}
	}
	parts := parseURL(pattern)
	r.roots[method].insert(pattern, parts, 0)
}

// 从路由表中查找当前请求对应的参数
func (r *router) getRoute(method string, pattern string) (*trie, map[string]string) {
	parts := parseURL(pattern)
	_, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := r.roots[method]
	target := node.search(parts, 0)
	params := make(map[string]string, 0)
	if target != nil {
		originParts := parseURL(target.pattern)
		for index, part := range originParts {
			if part[0] == ':' {
				params[part[1:]] = parts[index]
			} else if part[0] == '*' {
				params[part[1:]] = strings.Join(parts[index:], "/")
			}
		}
		return target, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node == nil {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	} else {
		c.Params = params
		key := c.Method + "-" + node.pattern
		r.handlers[key](c)
	}
}
