package pkg

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // route to be match, eg: /p/:lang
	part     string  // part of route, eg: :lang
	children []*node // child node, eg: [doc, tutorial, intro]
	isWild   bool    // whether to match exactly, true if part contains: or *
}

// String echo node info
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t, chilren=%d}", n.pattern, n.part, n.isWild, len(n.children))
}

// matchChild the first matching node is used for insertion
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren all nodes that match successfully, used to find
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// travel get nodes that node with pattern
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// insert
func (n *node) insert(pattern string, parts []string, height int) {
	// root '/' also goes here
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
