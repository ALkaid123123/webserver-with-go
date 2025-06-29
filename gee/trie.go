package gee

import (
	"strings"
)

type trie struct {
	pattern  string  //待匹配链接
	part     string  //当前节点需要匹配的字符串
	children []*trie //子节点
}

func (node *trie) searchChild(part string) *trie {
	for _, child := range node.children {
		if child.part == part || strings.HasPrefix(child.part, ":") || strings.HasPrefix(child.part, "*") {
			return child
		}
	}
	return nil
}

func (node *trie) searchChildren(part string) []*trie {
	nodes := make([]*trie, 0)
	for _, child := range node.children {
		if child.part == part || strings.HasPrefix(child.part, ":") || strings.HasPrefix(child.part, "*") {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (node *trie) insert(pattern string, parts []string, index int) {
	if index == len(parts) {
		node.pattern = pattern
		return
	}
	part := parts[index]
	ch := node.searchChild(part)
	if ch == nil {
		// fmt.Print(index)
		ch = &trie{pattern: pattern, part: part}
		node.children = append(node.children, ch)
	}
	ch.insert(pattern, parts, index+1)
}

func (node *trie) search(parts []string, index int) *trie {
	if index == len(parts) || strings.HasPrefix(node.part, "*") {
		if node.pattern == "" {
			return nil
		}
		return node
	}
	part := parts[index]
	children := node.searchChildren(part)
	for _, ch := range children {
		res := ch.search(parts, index+1)
		if res != nil {
			return res
		}
	}
	return nil
}
