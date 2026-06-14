package model

import (
	"testing"

	"github.com/blacksheepaul/timelog/model/gen"
	sqlite "github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
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

// TestMoveCategoryUpdatesDescendantLevels verifies that moving a category
// updates both its own level/path and those of all its descendants.
func TestMoveCategoryUpdatesDescendantLevels(t *testing.T) {
	db, err := openTestDB()
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	if err := db.AutoMigrate(&gen.Category{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	// Create two roots, one child under rootA, one grandchild under child.
	rootA := mustCreateCategory(t, db, "RootA", nil)
	rootB := mustCreateCategory(t, db, "RootB", nil)
	child := mustCreateCategory(t, db, "Child", rootA.ID)
	grandchild := mustCreateCategory(t, db, "Grandchild", child.ID)

	if *child.Level != 2 {
		t.Fatalf("expected child level 2, got %d", *child.Level)
	}
	if *grandchild.Level != 3 {
		t.Fatalf("expected grandchild level 3, got %d", *grandchild.Level)
	}

	// Move child (with its descendant) from rootA to rootB.
	if err := MoveCategory(db, *child.ID, rootB.ID); err != nil {
		t.Fatalf("move category: %v", err)
	}

	// Reload and verify.
	movedChild, err := GetCategoryByID(db, *child.ID)
	if err != nil {
		t.Fatalf("get moved child: %v", err)
	}
	if *movedChild.ParentID != *rootB.ID {
		t.Errorf("child parent: want %d, got %d", *rootB.ID, *movedChild.ParentID)
	}
	if *movedChild.Level != 2 {
		t.Errorf("child level: want 2, got %d", *movedChild.Level)
	}
	if *movedChild.Path != "/RootB" {
		t.Errorf("child path: want /RootB, got %s", *movedChild.Path)
	}

	movedGrandchild, err := GetCategoryByID(db, *grandchild.ID)
	if err != nil {
		t.Fatalf("get moved grandchild: %v", err)
	}
	if *movedGrandchild.Level != 3 {
		t.Errorf("grandchild level: want 3, got %d", *movedGrandchild.Level)
	}
	if *movedGrandchild.Path != "/RootB/Child" {
		t.Errorf("grandchild path: want /RootB/Child, got %s", *movedGrandchild.Path)
	}
}

// TestMoveCategoryToRoot verifies that moving a category to root level
// decreases its own level and its descendants' levels accordingly.
func TestMoveCategoryToRoot(t *testing.T) {
	db, err := openTestDB()
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	if err := db.AutoMigrate(&gen.Category{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	root := mustCreateCategory(t, db, "Root", nil)
	child := mustCreateCategory(t, db, "Child", root.ID)
	grandchild := mustCreateCategory(t, db, "Grandchild", child.ID)

	// Move child to root.
	if err := MoveCategory(db, *child.ID, nil); err != nil {
		t.Fatalf("move category to root: %v", err)
	}

	movedChild, err := GetCategoryByID(db, *child.ID)
	if err != nil {
		t.Fatalf("get moved child: %v", err)
	}
	if movedChild.ParentID != nil {
		t.Errorf("child parent: want nil, got %v", *movedChild.ParentID)
	}
	if *movedChild.Level != 1 {
		t.Errorf("child level: want 1, got %d", *movedChild.Level)
	}
	if *movedChild.Path != "/" {
		t.Errorf("child path: want /, got %s", *movedChild.Path)
	}

	movedGrandchild, err := GetCategoryByID(db, *grandchild.ID)
	if err != nil {
		t.Fatalf("get moved grandchild: %v", err)
	}
	if *movedGrandchild.Level != 2 {
		t.Errorf("grandchild level: want 2, got %d", *movedGrandchild.Level)
	}
	if *movedGrandchild.Path != "/Child" {
		t.Errorf("grandchild path: want /Child, got %s", *movedGrandchild.Path)
	}
}

func openTestDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
}

func mustCreateCategory(t *testing.T, db *gorm.DB, name string, parentID *int32) *gen.Category {
	t.Helper()
	color := "#3B82F6"
	desc := ""
	cat := &gen.Category{
		Name:        name,
		Color:       &color,
		Description: &desc,
		ParentID:    parentID,
	}
	if err := CreateCategory(db, cat); err != nil {
		t.Fatalf("create category %s: %v", name, err)
	}
	return cat
}
