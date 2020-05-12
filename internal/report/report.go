package report

import (
	"fmt"
)

type RefResult struct {
	Name string
	Kind string
	Namespace string
	Ownerkind string
	OwnerReference string
}

type Node struct {
	name string
	namespace string
	kind string
	child []*Node
}

type NodeTable struct {
	table map[string]*Node
	root []*Node
}
func NewNodeTable()  *NodeTable{
	//return New(map[string]*Node{})
	m:=make(map[string]*Node)
	r:=[]*Node{}
	return  &NodeTable{table:m,root:r}
}

func (nodeTable *NodeTable) AddNode(name,namespace, kind, parentName string) {
	fmt.Printf("add: name=%s namespace=%s kind=%s parentId=%s\n", name,namespace, kind, parentName)
	node := &Node{name: name,namespace:namespace,kind:kind, child: []*Node{}}
	if parentName == "" {
		//check if this parent already exists
		_, ok := nodeTable.table[parentName]
		if !ok {
			nodeTable.root =append(nodeTable.root,node)
		}

	} else {
		parent, ok := nodeTable.table[parentName]
		if !ok {
		fmt.Printf("add: parentId=%v: not found\n", parentName)
		return
		}
		parent.child = append(parent.child, node)
	}
		nodeTable.table[name] = node
}

func showNode(node *Node, prefix string) {

	if prefix == "" {
		fmt.Printf("%v\\%v\\%v\n\n", node.name , node.namespace,node.kind)
	} else {
		fmt.Printf("%v %v\\%v\n\n", prefix, node.name,node.kind)
	}
	println ("")

		for _, n := range node.child{
			showNode(n, prefix+"--")
		}



}

func (nodeTable *NodeTable)  ShowAll() {

	if nodeTable.root == nil {
		fmt.Printf("show: root node not found\n")
		return
	}
	for _, n := range  nodeTable.root{
		fmt.Println("")
		showNode(n, "")
	}

}
