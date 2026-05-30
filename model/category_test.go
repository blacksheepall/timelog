package model

import (
	"testing"

	"github.com/blacksheepaul/timelog/model/gen"
)

// TestBuildCategoryTreePointers tests that the tree building uses pointers correctly
// so that children are preserved in the returned structure
func TestBuildCategoryTreePointers(t *testing.T) {
	// Create test categories
	rootID := int32(1)
	childID := int32(2)
	grandchildID := int32(3)
	rootColor := "#FF0000"
	rootDesc := "Root category"
	rootPath := "/"
	childColor := "#00FF00"
	childDesc := "Child category"
	childPath := "/Root"
	grandchildColor := "#0000FF"
	grandchildDesc := "Grandchild category"
	grandchildPath := "/Root/Child"
	levelOne := int32(1)
	levelTwo := int32(2)
	levelThree := int32(3)

	categories := []gen.Category{
		{
			ID:          &rootID,
			Name:        "Root",
			Color:       &rootColor,
			Description: &rootDesc,
			ParentID:    nil,
			Level:       &levelOne,
			Path:        &rootPath,
		},
		{
			ID:          &childID,
			Name:        "Child",
			Color:       &childColor,
			Description: &childDesc,
			ParentID:    &rootID,
			Level:       &levelTwo,
			Path:        &childPath,
		},
		{
			ID:          &grandchildID,
			Name:        "Grandchild",
			Color:       &grandchildColor,
			Description: &grandchildDesc,
			ParentID:    &childID,
			Level:       &levelThree,
			Path:        &grandchildPath,
		},
	}

	// Build the tree
	tree := buildCategoryTree(categories)

	// Verify tree structure
	if len(tree) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(tree))
	}

	rootNode := tree[0]
	if *rootNode.Category.ID != rootID {
		t.Errorf("Expected root ID %d, got %d", rootID, *rootNode.Category.ID)
	}

	if len(rootNode.Children) != 1 {
		t.Fatalf("Expected root to have 1 child, got %d. This suggests children array was copied by value instead of using pointers.", len(rootNode.Children))
	}

	childNode := rootNode.Children[0]
	if *childNode.Category.ID != childID {
		t.Errorf("Expected child ID %d, got %d", childID, *childNode.Category.ID)
	}

	if len(childNode.Children) != 1 {
		t.Fatalf("Expected child to have 1 grandchild, got %d. This suggests children array was copied by value instead of using pointers.", len(childNode.Children))
	}

	grandchildNode := childNode.Children[0]
	if *grandchildNode.Category.ID != grandchildID {
		t.Errorf("Expected grandchild ID %d, got %d", grandchildID, *grandchildNode.Category.ID)
	}

	if len(grandchildNode.Children) != 0 {
		t.Errorf("Expected grandchild to have no children, got %d", len(grandchildNode.Children))
	}
}

// TestBuildCategoryTreeMultipleRoots tests tree building with multiple root nodes
func TestBuildCategoryTreeMultipleRoots(t *testing.T) {
	root1ID := int32(1)
	root2ID := int32(2)
	child1ID := int32(3)
	child2ID := int32(4)
	levelOne := int32(1)
	levelTwo := int32(2)

	categories := []gen.Category{
		{
			ID:       &root1ID,
			Name:     "Root1",
			ParentID: nil,
			Level:    &levelOne,
		},
		{
			ID:       &root2ID,
			Name:     "Root2",
			ParentID: nil,
			Level:    &levelOne,
		},
		{
			ID:       &child1ID,
			Name:     "Child1",
			ParentID: &root1ID,
			Level:    &levelTwo,
		},
		{
			ID:       &child2ID,
			Name:     "Child2",
			ParentID: &root2ID,
			Level:    &levelTwo,
		},
	}

	tree := buildCategoryTree(categories)

	if len(tree) != 2 {
		t.Fatalf("Expected 2 root nodes, got %d", len(tree))
	}

	// Verify each root has exactly one child
	for _, root := range tree {
		if len(root.Children) != 1 {
			t.Errorf("Expected root '%s' to have 1 child, got %d", root.Category.Name, len(root.Children))
		}
	}
}

// TestBuildCategoryTreeEmptyInput tests with empty category list
func TestBuildCategoryTreeEmptyInput(t *testing.T) {
	categories := []gen.Category{}
	tree := buildCategoryTree(categories)

	if len(tree) != 0 {
		t.Errorf("Expected empty tree, got %d nodes", len(tree))
	}
}

// TestBuildCategoryTreeSingleRoot tests with only one root category
func TestBuildCategoryTreeSingleRoot(t *testing.T) {
	rootID := int32(1)
	levelOne := int32(1)

	categories := []gen.Category{
		{
			ID:       &rootID,
			Name:     "OnlyRoot",
			ParentID: nil,
			Level:    &levelOne,
		},
	}

	tree := buildCategoryTree(categories)

	if len(tree) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(tree))
	}

	if len(tree[0].Children) != 0 {
		t.Errorf("Expected root to have no children, got %d", len(tree[0].Children))
	}
}
