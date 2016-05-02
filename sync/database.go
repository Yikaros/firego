package sync

import (
	"strings"
)

type Database struct {
	root *Node
}

func NewDB() *Database {
	return &Database{
		root: &Node{
			Children: map[string]*Node{},
		},
	}
}

func (d *Database) Add(path string, n *Node) {
	if path == "" {
		d.root = n
		return
	}

	rabbitHole := strings.Split(path, "/")
	current := d.root
	for i := 0; i < len(rabbitHole)-1; i++ {
		step := rabbitHole[i]
		next, ok := current.Children[step]
		if !ok {
			next = &Node{
				Parent:   current,
				Key:      step,
				Children: map[string]*Node{},
			}
			current.Children[step] = next
		}
		next.Value = nil // no long has a value since it now has a child
		current, next = next, nil
	}

	lastPath := rabbitHole[len(rabbitHole)-1]
	current.Children[lastPath] = n
	n.Parent = current
}

func (d *Database) Update(path string, n *Node) {
	current := d.root
	rabbitHole := strings.Split(path, "/")

	for i := 0; i < len(rabbitHole); i++ {
		step := rabbitHole[i]
		if step == "" {
			// prevent against empty strings due to strings.Split
			continue
		}
		next, ok := current.Children[step]
		if !ok {
			next = &Node{
				Parent:   current,
				Key:      step,
				Children: map[string]*Node{},
			}
			current.Children[step] = next
		}
		next.Value = nil // no long has a value since it now has a child
		current, next = next, nil
	}

	current.merge(n)

}

func (d *Database) Del(path string) {
	if path == "" {
		d.root = &Node{
			Children: map[string]*Node{},
		}
		return
	}

	rabbitHole := strings.Split(path, "/")
	current := d.root

	// traverse to target node's parent
	var delIdx int
	for ; delIdx < len(rabbitHole)-1; delIdx++ {
		next, ok := current.Children[rabbitHole[delIdx]]
		if !ok {
			// item does not exist, no need to do anything
			return
		}

		current = next
	}

	endNode := current
	leafPath := rabbitHole[len(rabbitHole)-1]
	delete(endNode.Children, leafPath)

	for tmp := endNode.prune(); tmp != nil; tmp = tmp.prune() {
		delIdx--
		endNode = tmp
	}

	if endNode != nil {
		delete(endNode.Children, rabbitHole[delIdx])
	}
}

func (d *Database) Get(path string) *Node {
	current := d.root
	if path == "" {
		return current
	}

	rabbitHole := strings.Split(path, "/")
	for i := 0; i < len(rabbitHole); i++ {
		var ok bool
		current, ok = current.Children[rabbitHole[i]]
		if !ok {
			return nil
		}
	}
	return current
}
