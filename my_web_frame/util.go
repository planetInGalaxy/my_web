package my_web_frame

import "fmt"

func (n *node) printAll() {
	fmt.Println(*n)
	if n.children == nil {
		return
	}
	for _, child := range n.children {
		child.printAll()
	}
}
