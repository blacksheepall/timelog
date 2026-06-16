package integration_test

import (
	"testing"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestCategoryTreeStructure(t *testing.T) {
	db := testDB().Db()
	db.Exec("DELETE FROM categories")

	svc := getTestService()

	// Create test hierarchy: Root -> Child -> Grandchild
	root := &domain.Category{
		Name: "Root Category",
	}
	if err := svc.CreateCategory(root); err != nil {
		t.Fatalf("Failed to create root category: %v", err)
	}

	child := &domain.Category{
		Name:     "Child Category",
		ParentID: &root.ID,
	}
	if err := svc.CreateCategory(child); err != nil {
		t.Fatalf("Failed to create child category: %v", err)
	}

	grandchild := &domain.Category{
		Name:     "Grandchild Category",
		ParentID: &child.ID,
	}
	if err := svc.CreateCategory(grandchild); err != nil {
		t.Fatalf("Failed to create grandchild category: %v", err)
	}

	// Get the tree
	tree, err := svc.GetCategoryTree()
	if err != nil {
		t.Fatalf("Failed to get category tree: %v", err)
	}

	// Verify tree structure
	if len(tree) != 1 {
		t.Errorf("Expected 1 root node, got %d", len(tree))
	}

	rootNode := tree[0]
	if rootNode.Category.Name != "Root Category" {
		t.Errorf("Expected root name 'Root Category', got '%s'", rootNode.Category.Name)
	}

	if len(rootNode.Children) != 1 {
		t.Fatalf("Expected root to have 1 child, got %d", len(rootNode.Children))
	}

	childNode := rootNode.Children[0]
	if childNode.Category.Name != "Child Category" {
		t.Errorf("Expected child name 'Child Category', got '%s'", childNode.Category.Name)
	}

	if len(childNode.Children) != 1 {
		t.Fatalf("Expected child to have 1 grandchild, got %d", len(childNode.Children))
	}

	grandchildNode := childNode.Children[0]
	if grandchildNode.Category.Name != "Grandchild Category" {
		t.Errorf("Expected grandchild name 'Grandchild Category', got '%s'", grandchildNode.Category.Name)
	}

	if len(grandchildNode.Children) != 0 {
		t.Errorf("Expected grandchild to have no children, got %d", len(grandchildNode.Children))
	}

	// Verify IDs are preserved
	if rootNode.Children[0].Category.ID != child.ID {
		t.Error("Child node ID doesn't match expected child ID")
	}

	if rootNode.Children[0].Children[0].Category.ID != grandchild.ID {
		t.Error("Grandchild node ID doesn't match expected grandchild ID")
	}
}

func TestCategoryTreeMultipleRoots(t *testing.T) {
	db := testDB().Db()
	db.Exec("DELETE FROM categories")

	svc := getTestService()

	// Create multiple root categories with children
	root1 := &domain.Category{
		Name: "Root 1",
	}
	if err := svc.CreateCategory(root1); err != nil {
		t.Fatalf("Failed to create root1: %v", err)
	}

	child1 := &domain.Category{
		Name:     "Child 1",
		ParentID: &root1.ID,
	}
	if err := svc.CreateCategory(child1); err != nil {
		t.Fatalf("Failed to create child1: %v", err)
	}

	root2 := &domain.Category{
		Name: "Root 2",
	}
	if err := svc.CreateCategory(root2); err != nil {
		t.Fatalf("Failed to create root2: %v", err)
	}

	child2 := &domain.Category{
		Name:     "Child 2",
		ParentID: &root2.ID,
	}
	if err := svc.CreateCategory(child2); err != nil {
		t.Fatalf("Failed to create child2: %v", err)
	}

	// Get the tree
	tree, err := svc.GetCategoryTree()
	if err != nil {
		t.Fatalf("Failed to get category tree: %v", err)
	}

	// Verify we have 2 root nodes
	if len(tree) != 2 {
		t.Errorf("Expected 2 root nodes, got %d", len(tree))
	}

	// Verify each root has its child
	for _, rootNode := range tree {
		if len(rootNode.Children) != 1 {
			t.Errorf("Expected root '%s' to have 1 child, got %d", rootNode.Category.Name, len(rootNode.Children))
		}
	}
}
