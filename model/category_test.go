package model_test

import (
	"errors"
	"testing"

	"github.com/blacksheepaul/timelog/internal/testutil"
	. "github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/core/errs"
	"gorm.io/gorm"
)

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

// --- buildCategoryTree tests (preserved from original) ---

func TestBuildCategoryTreePointers(t *testing.T) {
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
		{ID: &rootID, Name: "Root", Color: &rootColor, Description: &rootDesc, ParentID: nil, Level: &levelOne, Path: &rootPath},
		{ID: &childID, Name: "Child", Color: &childColor, Description: &childDesc, ParentID: &rootID, Level: &levelTwo, Path: &childPath},
		{ID: &grandchildID, Name: "Grandchild", Color: &grandchildColor, Description: &grandchildDesc, ParentID: &childID, Level: &levelThree, Path: &grandchildPath},
	}

	tree := BuildCategoryTree(categories)
	if len(tree) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(tree))
	}
	rootNode := tree[0]
	if *rootNode.Category.ID != rootID {
		t.Errorf("Expected root ID %d, got %d", rootID, *rootNode.Category.ID)
	}
	if len(rootNode.Children) != 1 {
		t.Fatalf("Expected root to have 1 child, got %d", len(rootNode.Children))
	}
	childNode := rootNode.Children[0]
	if *childNode.Category.ID != childID {
		t.Errorf("Expected child ID %d, got %d", childID, *childNode.Category.ID)
	}
	if len(childNode.Children) != 1 {
		t.Fatalf("Expected child to have 1 grandchild, got %d", len(childNode.Children))
	}
	grandchildNode := childNode.Children[0]
	if *grandchildNode.Category.ID != grandchildID {
		t.Errorf("Expected grandchild ID %d, got %d", grandchildID, *grandchildNode.Category.ID)
	}
}

func TestBuildCategoryTreeMultipleRoots(t *testing.T) {
	root1ID := int32(1)
	root2ID := int32(2)
	child1ID := int32(3)
	child2ID := int32(4)
	levelOne := int32(1)
	levelTwo := int32(2)

	categories := []gen.Category{
		{ID: &root1ID, Name: "Root1", ParentID: nil, Level: &levelOne},
		{ID: &root2ID, Name: "Root2", ParentID: nil, Level: &levelOne},
		{ID: &child1ID, Name: "Child1", ParentID: &root1ID, Level: &levelTwo},
		{ID: &child2ID, Name: "Child2", ParentID: &root2ID, Level: &levelTwo},
	}

	tree := BuildCategoryTree(categories)
	if len(tree) != 2 {
		t.Fatalf("Expected 2 root nodes, got %d", len(tree))
	}
	for _, root := range tree {
		if len(root.Children) != 1 {
			t.Errorf("Expected root '%s' to have 1 child, got %d", root.Category.Name, len(root.Children))
		}
	}
}

func TestBuildCategoryTreeEmptyInput(t *testing.T) {
	tree := BuildCategoryTree([]gen.Category{})
	if len(tree) != 0 {
		t.Errorf("Expected empty tree, got %d nodes", len(tree))
	}
}

func TestBuildCategoryTreeSingleRoot(t *testing.T) {
	rootID := int32(1)
	levelOne := int32(1)
	categories := []gen.Category{
		{ID: &rootID, Name: "OnlyRoot", ParentID: nil, Level: &levelOne},
	}
	tree := BuildCategoryTree(categories)
	if len(tree) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(tree))
	}
	if len(tree[0].Children) != 0 {
		t.Errorf("Expected root to have no children, got %d", len(tree[0].Children))
	}
}

// --- move tests (preserved from original, now with testutil) ---

func TestMoveCategoryUpdatesDescendantLevels(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

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

	if err := MoveCategory(db, *child.ID, rootB.ID); err != nil {
		t.Fatalf("move category: %v", err)
	}

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

func TestMoveCategoryToRoot(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	root := mustCreateCategory(t, db, "Root", nil)
	child := mustCreateCategory(t, db, "Child", root.ID)
	grandchild := mustCreateCategory(t, db, "Grandchild", child.ID)

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

// --- new bad-case tests ---

func TestValidateLevel(t *testing.T) {
	if err := ValidateLevel(1); err != nil {
		t.Fatalf("level 1 should be valid: %v", err)
	}
	if err := ValidateLevel(3); err != nil {
		t.Fatalf("level 3 should be valid: %v", err)
	}
	if err := ValidateLevel(0); err == nil {
		t.Fatal("level 0 should be invalid")
	}
	if err := ValidateLevel(4); err == nil {
		t.Fatal("level 4 should be invalid")
	}
}

func TestGetFullPath(t *testing.T) {
	rootPath := "/"
	root := &gen.Category{Name: "Root", Path: &rootPath}
	if got := GetFullPath(root); got != "/Root" {
		t.Fatalf("root full path: want /Root, got %s", got)
	}

	childPath := "/Root"
	child := &gen.Category{Name: "Child", Path: &childPath}
	if got := GetFullPath(child); got != "/Root/Child" {
		t.Fatalf("child full path: want /Root/Child, got %s", got)
	}
}

func TestCreateCategoryRejectsInvalidParentID(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	invalid := int32(0)
	cat := &gen.Category{Name: "bad", ParentID: &invalid}
	if err := CreateCategory(db, cat); !errors.Is(err, errs.ErrInvalidParentID) {
		t.Fatalf("expected ErrInvalidParentID, got %v", err)
	}
}

func TestCreateCategoryRejectsMissingParent(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	missing := int32(9999)
	cat := &gen.Category{Name: "orphan", ParentID: &missing}
	if err := CreateCategory(db, cat); err == nil {
		t.Fatal("expected error for missing parent")
	}
}

func TestCreateCategoryRejectsExceedingMaxLevel(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	root := mustCreateCategory(t, db, "Root", nil)
	child := mustCreateCategory(t, db, "Child", root.ID)
	grandchild := mustCreateCategory(t, db, "Grandchild", child.ID)

	fourth := &gen.Category{Name: "TooDeep", ParentID: grandchild.ID}
	if err := CreateCategory(db, fourth); err == nil {
		t.Fatal("expected error exceeding max level")
	}
}

func TestGetCategoryByIDNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	_, err := GetCategoryByID(db, 9999)
	if err == nil {
		t.Fatal("expected error for non-existent category")
	}
}

func TestGetCategoryByNameNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	_, err := GetCategoryByName(db, "missing", nil)
	if err == nil {
		t.Fatal("expected error for non-existent category name")
	}
}

func TestUpdateCategoryPreservesHierarchy(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	root := mustCreateCategory(t, db, "Root", nil)
	child := mustCreateCategory(t, db, "Child", root.ID)

	newParent := int32(9999)
	newLevel := int32(99)
	newPath := "/hacked"
	child.ParentID = &newParent
	child.Level = &newLevel
	child.Path = &newPath
	child.Name = "Renamed"
	if err := UpdateCategory(db, child); err != nil {
		t.Fatalf("UpdateCategory: %v", err)
	}

	got, err := GetCategoryByID(db, *child.ID)
	if err != nil {
		t.Fatalf("GetCategoryByID: %v", err)
	}
	if got.Name != "Renamed" {
		t.Fatalf("name not updated: %s", got.Name)
	}
	if *got.ParentID != *root.ID {
		t.Fatalf("parent_id was changed: %d", *got.ParentID)
	}
	if *got.Level != 2 {
		t.Fatalf("level was changed: %d", *got.Level)
	}
}

func TestMoveCategoryRejectsSelf(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	root := mustCreateCategory(t, db, "Root", nil)
	if err := MoveCategory(db, *root.ID, root.ID); err == nil {
		t.Fatal("expected error moving category to itself")
	}
}

func TestMoveCategoryRejectsDescendant(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	root := mustCreateCategory(t, db, "Root", nil)
	child := mustCreateCategory(t, db, "Child", root.ID)
	if err := MoveCategory(db, *root.ID, child.ID); err == nil {
		t.Fatal("expected error moving category to its own child")
	}
}

func TestMoveCategoryRejectsExceedingMaxLevel(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	root := mustCreateCategory(t, db, "Root", nil)
	child := mustCreateCategory(t, db, "Child", root.ID)
	grandchild := mustCreateCategory(t, db, "Grandchild", child.ID)

	otherRoot := mustCreateCategory(t, db, "OtherRoot", nil)
	deep := mustCreateCategory(t, db, "Deep", otherRoot.ID)
	_ = mustCreateCategory(t, db, "Deeper", deep.ID)

	if err := MoveCategory(db, *otherRoot.ID, grandchild.ID); err == nil {
		t.Fatal("expected error exceeding max level after move")
	}
}

func TestListCategoriesDirectly(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()
	_ = mustCreateCategory(t, db, "A", nil)

	list, err := ListCategories(db)
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	if len(list) < 1 {
		t.Fatalf("expected at least 1 category, got %d", len(list))
	}
}

func TestListCategoriesByLevelDirectly(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()
	_ = mustCreateCategory(t, db, "LevelOne", nil)

	list, err := ListCategoriesByLevel(db, 1)
	if err != nil {
		t.Fatalf("ListCategoriesByLevel: %v", err)
	}
	if len(list) < 1 {
		t.Fatalf("expected at least 1 category, got %d", len(list))
	}
}

func TestGetCategoriesByParentIDDirectly(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()
	root := mustCreateCategory(t, db, "RootParent", nil)
	_ = mustCreateCategory(t, db, "ChildOfRoot", root.ID)

	children, err := GetCategoriesByParentID(db, root.ID)
	if err != nil {
		t.Fatalf("GetCategoriesByParentID: %v", err)
	}
	if len(children) < 1 {
		t.Fatalf("expected at least 1 child, got %d", len(children))
	}

	roots, err := GetCategoriesByParentID(db, nil)
	if err != nil {
		t.Fatalf("GetCategoriesByParentID nil: %v", err)
	}
	if len(roots) < 1 {
		t.Fatalf("expected at least 1 root, got %d", len(roots))
	}
}

func TestGetCategoryTreeDirectly(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()
	root := mustCreateCategory(t, db, "TreeRoot", nil)
	_ = mustCreateCategory(t, db, "TreeChild", root.ID)

	tree, err := GetCategoryTree(db)
	if err != nil {
		t.Fatalf("GetCategoryTree: %v", err)
	}
	if len(tree) < 1 {
		t.Fatalf("expected at least 1 root node, got %d", len(tree))
	}
}
