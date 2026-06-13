<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-text-primary">Category Management</h1>
      <button
        @click="toggleForm"
        class="inline-flex items-center px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand"
      >
        <PlusIcon class="h-5 w-5 mr-2" />
        New Category
      </button>
    </div>

    <!-- Category创建/编辑表单 -->
    <div v-if="showForm" class="bg-bg-surface shadow-md-md rounded-lg p-6">
      <h2 class="text-lg font-medium text-text-primary mb-6">
        {{ isEditing ? 'Edit Category' : 'Create New Category' }}
      </h2>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="name" class="block text-sm font-medium text-text-muted mb-2"> Name * </label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              maxlength="50"
              class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
              placeholder="Category name"
            />
          </div>

          <div>
            <label for="color" class="block text-sm font-medium text-text-muted mb-2">
              Color *
            </label>
            <div class="flex items-center space-x-2">
              <input
                id="color"
                v-model="form.color"
                type="color"
                required
                class="h-10 w-16 border border-default rounded-md"
              />
              <input
                v-model="form.color"
                type="text"
                class="flex-1 px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
                placeholder="#000000"
              />
            </div>
          </div>
        </div>

        <div v-if="!isEditing">
          <label for="parent_id" class="block text-sm font-medium text-text-muted mb-2">
            Parent Category (Optional)
          </label>
          <select
            id="parent_id"
            v-model="form.parent_id"
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
          >
            <option :value="null">-- Root Category --</option>
            <option v-for="cat in availableParents" :key="cat.id" :value="cat.id">
              {{ cat.path === '/' ? '' : cat.path.replace(/\//g, ' / ') + ' / ' }}{{ cat.name }}
            </option>
          </select>
          <p class="mt-1 text-sm text-text-secondary">Maximum depth: 3 levels (Root/Child/Grandchild)</p>
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-text-muted mb-2">
            Description
          </label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            maxlength="200"
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="Describe what this category is used for..."
          ></textarea>
        </div>

        <div class="flex justify-end space-x-4">
          <button
            type="button"
            @click="handleCancel"
            class="px-4 py-2 text-sm font-medium text-text-muted bg-bg-surface border border-default rounded-md shadow-sm hover:bg-bg-elevated focus:outline-none focus:ring-2 focus:ring-brand"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="submitting"
            class="px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ submitting ? 'Saving...' : isEditing ? 'Update' : 'Create' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Categories树形列表 -->
    <div class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <h2 class="text-lg font-medium text-text-primary">All Categories</h2>
      </div>

      <div v-if="loading" class="p-6 text-center">
        <div
          class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-brand"
        ></div>
        <p class="mt-2 text-text-secondary">Loading...</p>
      </div>

      <div v-else-if="error" class="p-6 text-center text-danger">
        {{ error }}
      </div>

      <div v-else-if="categoryTree.length === 0" class="p-6 text-center text-text-secondary">
        No categories found. Create your first one!
      </div>

      <div v-else class="divide-y divide-default">
        <CategoryTreeNode
          v-for="node in categoryTree"
          :key="node.category.id"
          :node="node"
          :level="0"
          @edit="handleEdit"
          @delete="handleDelete"
          @move="openMoveDialog"
        />
      </div>
    </div>

    <!-- 移动分类对话框 -->
    <div
      v-if="showMoveDialog"
      class="fixed inset-0 bg-bg-elevated0 bg-opacity-75 flex items-center justify-center z-50"
    >
      <div class="bg-bg-surface rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-medium text-text-primary mb-4">Move Category</h3>
        <p class="text-sm text-text-secondary mb-4">
          Select a new parent for "{{ movingCategory?.name }}"
        </p>

        <select
          v-model="moveTargetParentId"
          class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand mb-4"
        >
          <option :value="null">-- Root Level --</option>
          <option v-for="cat in availableMoveTargets" :key="cat.id" :value="cat.id">
            {{ cat.path === '/' ? '' : cat.path.replace(/\//g, ' / ') + ' / ' }}{{ cat.name }}
          </option>
        </select>

        <div class="flex justify-end space-x-3">
          <button
            @click="closeMoveDialog"
            class="px-4 py-2 text-sm font-medium text-text-muted bg-bg-surface border border-default rounded-md hover:bg-bg-elevated"
          >
            Cancel
          </button>
          <button
            @click="confirmMove"
            :disabled="moving"
            class="px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md hover:bg-brand-hover disabled:opacity-50"
          >
            {{ moving ? 'Moving...' : 'Move' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, computed, inject } from 'vue'
  import { PlusIcon } from '@heroicons/vue/24/outline'
  import { categoryAPI } from '@/api'
  import CategoryTreeNode from '@/components/CategoryTreeNode.vue'
  import type { Category, CategoryNode, CreateCategoryRequest } from '@/types'

  // 注入全局通知功能
  const showNotification = inject('showNotification') as (
    type: 'success' | 'error',
    message: string
  ) => void

  const categoryTree = ref<CategoryNode[]>([])
  const allCategories = ref<Category[]>([])
  const loading = ref(false)
  const submitting = ref(false)
  const error = ref<string | null>(null)
  const showForm = ref(false)
  const editingCategory = ref<Category | undefined>()

  // 移动对话框状态
  const showMoveDialog = ref(false)
  const movingCategory = ref<Category | null>(null)
  const moveTargetParentId = ref<number | null>(null)
  const moving = ref(false)

  const isEditing = computed(() => !!editingCategory.value)

  const form = reactive({
    name: '',
    color: '#3B82F6',
    description: '',
    parent_id: null as number | null,
  })

  // 可作为父分类的选项（创建新分类时用）
  const availableParents = computed(() => {
    return allCategories.value.filter(c => c.level < 3)
  })

  // 可作为移动目标的选项（排除自己及其子分类）
  const availableMoveTargets = computed(() => {
    if (!movingCategory.value) return []
    return allCategories.value.filter(
      c =>
        c.id !== movingCategory.value!.id &&
        c.level < 3 &&
        !c.path.includes(`/${movingCategory.value!.name}`)
    )
  })

  const resetForm = () => {
    form.name = ''
    form.color = '#3B82F6'
    form.description = ''
    form.parent_id = null
  }

  const loadEditingData = () => {
    if (editingCategory.value) {
      form.name = editingCategory.value.name
      form.color = editingCategory.value.color
      form.description = editingCategory.value.description
    } else {
      resetForm()
    }
  }

  const loadCategories = async () => {
    loading.value = true
    error.value = null

    try {
      const [treeRes, allRes] = await Promise.all([categoryAPI.getTree(), categoryAPI.getAll()])
      categoryTree.value = treeRes.data || []
      allCategories.value = allRes.data || []
    } catch (err) {
      error.value = 'Failed to load categories'
      console.error('Error loading categories:', err)
      showNotification('error', 'Failed to load categories')
    } finally {
      loading.value = false
    }
  }

  const toggleForm = () => {
    showForm.value = !showForm.value
    if (!showForm.value) {
      editingCategory.value = undefined
    }
    loadEditingData()
  }

  const buildCreateRequest = (): CreateCategoryRequest => {
    const data: CreateCategoryRequest = {
      name: form.name.trim(),
      color: form.color,
      description: form.description.trim(),
    }
    if (form.parent_id !== null) {
      data.parent_id = form.parent_id
    }
    return data
  }

  const handleSubmit = async () => {
    submitting.value = true

    try {
      if (editingCategory.value) {
        await categoryAPI.update(editingCategory.value.id, {
          name: form.name.trim(),
          color: form.color,
          description: form.description.trim(),
        })
        showNotification('success', 'Category updated successfully')
      } else {
        await categoryAPI.create(buildCreateRequest())
        showNotification('success', 'Category created successfully')
      }

      await loadCategories()
      showForm.value = false
      editingCategory.value = undefined
      resetForm()
    } catch (err: any) {
      console.error('Error saving category:', err)
      showNotification('error', err.response?.data?.message || 'Failed to save category')
    } finally {
      submitting.value = false
    }
  }

  const handleEdit = (category: Category) => {
    editingCategory.value = category
    showForm.value = true
    loadEditingData()
  }

  const handleCancel = () => {
    showForm.value = false
    editingCategory.value = undefined
    resetForm()
  }

  const handleDelete = async (id: number) => {
    if (
      !confirm(
        'Are you sure you want to delete this category? This will also delete all its sub-categories. This action cannot be undone.'
      )
    ) {
      return
    }

    try {
      await categoryAPI.delete(id)
      showNotification('success', 'Category deleted successfully')
      await loadCategories()
    } catch (err: any) {
      console.error('Error deleting category:', err)
      showNotification('error', err.response?.data?.message || 'Failed to delete category')
    }
  }

  const openMoveDialog = (category: Category) => {
    movingCategory.value = category
    moveTargetParentId.value = category.parent_id || null
    showMoveDialog.value = true
  }

  const closeMoveDialog = () => {
    showMoveDialog.value = false
    movingCategory.value = null
    moveTargetParentId.value = null
  }

  const confirmMove = async () => {
    if (!movingCategory.value) return

    moving.value = true
    try {
      await categoryAPI.move(movingCategory.value.id, moveTargetParentId.value)
      showNotification('success', 'Category moved successfully')
      await loadCategories()
      closeMoveDialog()
    } catch (err: any) {
      console.error('Error moving category:', err)
      showNotification('error', err.response?.data?.message || 'Failed to move category')
    } finally {
      moving.value = false
    }
  }

  onMounted(() => {
    loadCategories()
  })
</script>
