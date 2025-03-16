package shell

import (
	"fmt"
	"reflect"

	"github.com/rivo/tview"
)

func newNode(lable string) *tview.TreeNode {
	return tview.NewTreeNode(lable).
		SetTextStyle(outputTheme).
		SetSelectable(false)
}
func buildTree(data any) *tview.TreeNode {
	root := newNode("root")
	buildSubTree(root, data)
	return root
}
func buildSubTree(node *tview.TreeNode, data any) {
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Map:
		// Iterate over map keys
		for _, key := range val.MapKeys() {
			// Get the value corresponding to the key
			val := val.MapIndex(key).Interface()
			// Create a new node for the key
			subNode := newNode(fmt.Sprintf("%v: %v", key, val))
			// Recursively build the subtree for the value
			buildSubTree(subNode, val)
			// Add the subNode to the parent node
			node.AddChild(subNode)
		}
	case reflect.Slice, reflect.Array:
		// Iterate over each element in the slice/array
		for i := 0; i < val.Len(); i++ {
			// Create a new node for each index
			subNode := newNode(fmt.Sprintf("[%d]", i))
			// Recursively build the subtree for each element
			buildSubTree(subNode, val.Index(i).Interface())
			// Add the subNode to the parent node
			node.AddChild(subNode)
		}
	default:
		// For non-map, non-slice types, set the text as the value itself
		node.SetText(fmt.Sprintf("%v", data))
	}
}
