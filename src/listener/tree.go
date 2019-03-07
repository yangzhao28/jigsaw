package listener

import (
	"fmt"
	"strings"

	"github.com/yangzhao28/jigsaw/src/common"
)

type Tree struct {
	root      *node
	segmenter SegmentFunc
}

func New() *Tree {
	return &Tree{
		root:      newNode("root", nil),
		segmenter: pathSegmenter,
	}
}

func (tree *Tree) addToAllChildren(root *node, e common.Entry) *node {
	for _, n := range root.children {
		n.nodes = append(n.nodes, e)
		tree.addToAllChildren(n, e)
	}
	return root
}

func (tree *Tree) add(root *node, channel string, e common.Entry, existsOnly bool) *node {
	cur := root
	for part, i := tree.segmenter(channel, 0); len(part) != 0; part, i = tree.segmenter(channel, i) {
		switch part {
		case "*":
			// not end part
			if i > 0 {
				for _, n := range cur.children {
					tree.add(n, channel[i:], e, true)
				}
			} else {
				cur.wildcards = append(cur.wildcards, e)
				// end of channel, spread to all children
				for _, n := range cur.children {
					n.nodes = append(n.nodes, e)
					tree.addToAllChildren(n, e)
				}
			}
			return cur
		default:
			n := cur.children[part]
			if n == nil && existsOnly {
				return cur
			}
			if n == nil {
				n = newNode(part, cur)
				cur.children[part] = n
			}
			cur = n
		}
	}
	cur.nodes = append(cur.nodes, e)
	if !existsOnly {
		traceback := cur.prev
		for traceback != nil {
			cur.nodes = append(cur.nodes, traceback.wildcards...)
			traceback = traceback.prev
		}
	}
	return cur
}

func (tree *Tree) Add(channel string, e common.Entry) {
	tree.add(tree.root, channel, e, false)
}

func (tree *Tree) Get(path string) (ret []common.Entry) {
	cur := tree.root
	for part, i := tree.segmenter(path, 0); len(part) != 0; part, i = tree.segmenter(path, i) {
		// fmt.Println(">>>", part, cur.)
		switch part {
		case "*":
			for _, child := range cur.children {
				ret = append(ret, child.nodes...)
			}
			return
		default:
			if n, ok := cur.children[part]; ok {
				cur = n
			} else {
				return nil
			}
		}
	}
	return cur.nodes
}

func (tree *Tree) Reset() {
	tree.root = nil
}

func (tree *Tree) print(b *strings.Builder, lv int, n *node) *strings.Builder {
	for i := 0; i < lv; i++ {
		b.WriteString("    ")
	}
	printNonNil := func(n *node) string {
		if n == nil {
			return "nil"
		}
		if n.prev == nil {
			return "--"
		}
		return n.prev.key
	}

	b.WriteString("\\__")
	b.WriteString(fmt.Sprintf("key={%v} prev={%v} ", n.key, printNonNil(n)))
	b.WriteString(fmt.Sprintf("entry={%v} wildcards={%v}", n.nodes, n.wildcards))
	b.WriteString("\n")
	for _, child := range n.children {
		tree.print(b, lv+1, child)
	}
	return b
}

func (tree *Tree) String() string {
	b := &strings.Builder{}
	return tree.print(b, 0, tree.root).String()
}
