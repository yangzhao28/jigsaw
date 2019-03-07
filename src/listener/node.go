package listener

import "github.com/yangzhao28/jigsaw/src/common"

type node struct {
	prev      *node
	key       string
	nodes     []common.Entry
	wildcards []common.Entry
	children  map[string]*node
}

func newNode(key string, prev *node) *node {
	return &node{
		children: make(map[string]*node),
		prev:     prev,
		key:      key,
	}
}
