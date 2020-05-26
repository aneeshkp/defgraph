package report

import (
	"fmt"
)
type EdgeType string
var (
	EdgeTypeLink  EdgeType = "│"
	EdgeTypeMid   EdgeType = "├──"
	EdgeTypeNext  EdgeType = "──"
	EdgeTypeEnd   EdgeType = "└──"
)

type RefResult struct {
	Name string
	Kind string
	Namespace string
	Ownerkind string
	OwnerReference string
	Images []string
}

type Node struct {
	name string
	namespace string
	kind string
	images *[]string
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

func (nodeTable *NodeTable) AddNode(name,namespace, kind, parentName string, images *[]string) {
	fmt.Printf("add: name=%s namespace=%s kind=%s parentId=%s\n", name,namespace, kind, parentName)
	node := &Node{name: name,namespace:namespace,kind:kind, child: []*Node{},images:images}
	if parentName == "" {
		nodeTable.root =append(nodeTable.root,node)
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
		fmt.Printf("%v %v\\%v\\%v\n",string(EdgeTypeLink), node.name , node.namespace,node.kind)
		prefix=string(EdgeTypeEnd)
		if node.images !=nil{
			for _, image := range *node.images {
				fmt.Printf("%v %v\n", string(EdgeTypeLink), image)
			}
		}


	} else {
		fmt.Printf("%v %v\\%v\\%v\n", prefix+string(EdgeTypeNext), node.name,node.namespace,node.kind)
	}
	//println ("")
		for i, n := range node.child{
			if i==0 {
				showNode(n,  prefix + string(EdgeTypeNext))
			}else{
				showNode(n,  prefix + string(EdgeTypeEnd))
			}

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
