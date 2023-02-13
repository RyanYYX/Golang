package data_structure

type Comparable interface {
	Object
	Compare(c Comparable) int
}

type TreeNode struct {
	value  Comparable
	left   *TreeNode
	right  *TreeNode
	parent *TreeNode
}

func (node *TreeNode) Value() Comparable {
	return node.value
}

func (node *TreeNode) Left() *TreeNode {
	return node.left
}

func (node *TreeNode) Right() *TreeNode {
	return node.right
}

func (node *TreeNode) Parent() *TreeNode {
	return node.parent
}

type Tree interface {
	Query(val Comparable) (*TreeNode, error)
	Insert(val Comparable) (*TreeNode, error)
	Delete(val Comparable) (*TreeNode, error)
	Update(old, new Comparable) (*TreeNode, error)
}
