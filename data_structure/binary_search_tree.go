package data_structure

import "reflect"

var _ Tree = (*BinarySearchTree)(nil)

type BinarySearchTree struct {
	root  *TreeNode
	check TypeCheck
}

func (tree *BinarySearchTree) Root() *TreeNode {
	return tree.root
}

func (tree *BinarySearchTree) Query(val Comparable) (*TreeNode, error) {
	if  err := tree.check.Check(reflect.TypeOf(val)); err != nil {
		return nil, err
	}

	cur := tree.root
	for cur != nil {
		if val.Compare(cur.value) < 0 {
			cur = cur.left
		} else if val.Compare(cur.value) > 0 {
			cur = cur.right
		} else {
			return cur, nil
		}
	}
	return nil, nil
}

func (tree *BinarySearchTree) Insert(val Comparable) (*TreeNode, error) {
	if  err := tree.check.Check(reflect.TypeOf(val)); err != nil {
		return nil, err
	}

	parent, cur := tree.root, tree.root
	for cur != nil {
		parent = cur
		if val.Compare(cur.value) < 0 {
			cur = cur.left
		} else if val.Compare(cur.value) > 0 {
			cur = cur.right
		} else {
			return cur, nil
		}
	}

	cur = &TreeNode{value: val}
	cur.parent = parent
	if cur.value.Compare(parent.value) < 0 {
		parent.left = cur
	} else {
		parent.right = cur
	}
	return cur, nil
}

func (tree *BinarySearchTree) Delete(val Comparable) (*TreeNode, error) {
	if node, _ := tree.Query(val); node != nil {
		parent := node.parent
		if node.value.Compare(parent.value) < 0 {
			parent.left = nil
		} else {
			parent.right = nil
		}
		node.parent = nil
		return node, nil
	}
	return nil, nil
}

func (tree *BinarySearchTree) Update(old, new Comparable) (*TreeNode, error) {
	panic("implement me")
}