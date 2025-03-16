package shell

import (
	"fmt"
	"reflect"

	"github.com/fmotalleb/scrapper-go/utils"
	"github.com/rivo/tview"
)

// newNode creates a new tree node with the given label and collapsible behavior.
func newNode(label string, collapsible bool) *tview.TreeNode {
	node := tview.NewTreeNode(label).
		SetTextStyle(outputTheme).
		SetSelectable(collapsible)
	node.SetSelectedFunc(func() {
		if node.IsExpanded() {
			node.Collapse()
		} else {
			node.Expand()
		}
	})
	return node
}

// buildTree constructs a tree from the provided data and returns the root node.
func buildTree(data any) *tview.TreeNode {
	root := newNode("root", false)
	buildSubTree(root, data)
	return root
}

// buildSubTree recursively builds the tree for a given parent node and data.
func buildSubTree(node *tview.TreeNode, data any, replaceLabel ...bool) {
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Map:
		// Iterate over map keys and build subtree for each key-value pair.
		for _, key := range val.MapKeys() {
			val := val.MapIndex(key).Interface()
			subNode := newNode(fmt.Sprintf("%v: %v", key.Interface(), val), true)
			// If the value is not a string, build a sub-tree for it.
			buildSubTree(subNode, val, false)
			node.AddChild(subNode)
		}
	case reflect.Slice, reflect.Array:
		// Iterate over each element in the slice/array and build a subtree for each.
		for i := 0; i < val.Len(); i++ {
			subNode := newNode(fmt.Sprintf("[%d]", i), true)
			buildSubTree(subNode, val.Index(i).Interface())
			node.AddChild(subNode)
		}
	default:
		// For primitive types, set the text as the value itself.
		if utils.FirstOr(replaceLabel, true) {
			node.SetText(fmt.Sprintf("%v", data))
		}
	}
}
