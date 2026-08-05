<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-text-primary">Task Management</h1>
      <button
        @click="toggleForm"
        class="inline-flex items-center px-4 py-2 text-sm font-medium text-bg-surface bg-brand border border-transparent rounded-md shadow-sm hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand"
      >
        <PlusIcon class="h-5 w-5 mr-2" />
        New Task
      </button>
    </div>

    <!-- 任务创建/编辑表单 -->
    <div v-if="showForm" ref="taskFormRef" class="bg-bg-surface shadow-md-md rounded-lg p-6">
      <h2 class="text-lg font-medium text-text-primary mb-6">
        {{ isEditing ? 'Edit Task' : 'Create New Task' }}
      </h2>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="title" class="block text-sm font-medium text-text-muted mb-2">
            Title *
          </label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            required
            maxlength="200"
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="Task title"
          />
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-text-muted mb-2">
            Description
          </label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="Describe what this task involves..."
          ></textarea>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label for="category_id" class="block text-sm font-medium text-text-muted mb-2">
              Tag *
            </label>
            <select
              id="category_id"
              v-model="form.category_id"
              required
              class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            >
              <option value="" disabled>Select a tag</option>
              <option
                v-for="tag in availableCategories"
                :key="tag.id"
                :value="tag.id"
                :style="{ color: tag.color }"
              >
                {{ tag.name }} - {{ tag.description }}
              </option>
            </select>
          </div>

          <div>
            <label for="due_date" class="block text-sm font-medium text-text-muted mb-2">
              Due Date *
            </label>
            <input
              id="due_date"
              v-model="form.due_date"
              type="date"
              required
              class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            />
          </div>

          <div>
            <label for="estimated_minutes" class="block text-sm font-medium text-text-muted mb-2">
              Estimated Time (minutes) *
            </label>
            <input
              id="estimated_minutes"
              v-model.number="form.estimated_minutes"
              type="number"
              min="1"
              required
              class="w-full px-3 py-2 border border-default rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
              placeholder="e.g. 60"
            />
          </div>
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

    <!-- 任务列表 -->
    <div class="bg-bg-surface shadow-md-md rounded-lg">
      <div class="px-6 py-4 border-b border-default">
        <div class="flex justify-between items-center">
          <h2 class="text-lg font-medium text-text-primary">Tasks</h2>
          <div class="flex items-center space-x-4">
            <button
              @click="showSuspended = !showSuspended"
              class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors"
              :class="
                showSuspended
                  ? 'bg-yellow-100 text-yellow-800 border border-yellow-200'
                  : 'bg-bg-elevated text-text-secondary border border-default hover:bg-bg-elevated'
              "
            >
              ⏸️ Suspended
            </button>
            <button
              @click="showCompleted = !showCompleted"
              class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors"
              :class="
                showCompleted
                  ? 'bg-green-100 text-green-800 border border-success-border'
                  : 'bg-bg-elevated text-text-secondary border border-default hover:bg-bg-elevated'
              "
            >
              ✅ Completed
            </button>
            <select
              v-model="dateFilter"
              @change="loadTasks"
              class="px-3 py-1 text-sm border border-default rounded-md focus:outline-none focus:ring-2 focus:ring-brand"
            >
              <option value="">All dates</option>
              <option value="today">Today</option>
              <option value="tomorrow">Tomorrow</option>
              <option value="this-week">This week</option>
            </select>
          </div>
        </div>
      </div>

      <div v-if="loading" class="p-6 text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-brand"></div>
        <p class="mt-2 text-text-secondary">Loading...</p>
      </div>

      <div v-else-if="error" class="p-6 text-center text-danger">
        {{ error }}
      </div>

      <div v-else-if="filteredTasks.length === 0" class="p-6 text-center text-text-secondary">
        No tasks found. Create your first one!
      </div>

      <div v-else class="divide-y divide-default">
        <div
          v-for="task in filteredTasks"
          :key="task.id"
          class="p-6 hover:bg-bg-elevated"
          :class="{
            'opacity-60': task.is_completed,
            'bg-yellow-50': task.is_suspended,
            'hover:bg-yellow-100': task.is_suspended,
          }"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <h3
                  class="text-lg font-medium"
                  :class="
                    task.is_completed
                      ? 'line-through text-text-secondary'
                      : task.is_suspended
                        ? 'text-yellow-800'
                        : 'text-text-primary'
                  "
                >
                  {{ task.title }}
                </h3>
                <span
                  v-if="task.category"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-bg-surface"
                  :style="{ backgroundColor: task.category.color }"
                >
                  {{ task.category.name }}
                </span>
                <span
                  v-if="task.is_completed"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                >
                  ✓ Completed
                </span>
                <span
                  v-if="task.is_suspended"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"
                >
                  ⏸ Suspended
                </span>
              </div>

              <p v-if="task.description" class="text-text-secondary mb-3">
                {{ task.description }}
              </p>

              <div class="flex items-center space-x-4 text-sm text-text-secondary">
                <span>Due: {{ formatDate(task.due_date) }}</span>
                <span>Estimated: {{ task.estimated_minutes }}min</span>
                <span v-if="task.completed_at">
                  Completed: {{ formatDateTime(task.completed_at) }}
                </span>
              </div>
            </div>

            <div class="flex items-center space-x-2 ml-4">
              <button
                v-if="!task.is_suspended && !task.is_completed"
                @click="suspendTask(task.id)"
                class="px-3 py-1 text-xs font-medium text-yellow-700 bg-yellow-100 rounded-md hover:bg-yellow-200 focus:outline-none focus:ring-2 focus:ring-yellow-500"
              >
                Suspend
              </button>
              <button
                v-if="task.is_suspended"
                @click="unsuspendTask(task.id)"
                class="px-3 py-1 text-xs font-medium text-blue-700 bg-blue-100 rounded-md hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-brand"
              >
                Unsuspend
              </button>
              <button
                v-if="!task.is_completed && !task.is_suspended"
                @click="completeTask(task.id)"
                class="px-3 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md hover:bg-green-200 focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                Complete
              </button>
              <button
                v-if="task.is_completed"
                @click="incompleteTask(task.id)"
                class="px-3 py-1 text-xs font-medium text-text-muted bg-bg-elevated rounded-md hover:bg-bg-elevated focus:outline-none focus:ring-2 focus:ring-default"
              >
                Reopen
              </button>
              <button
                @click="handleEdit(task)"
                class="text-brand hover:text-blue-900 text-sm font-medium"
              >
                Edit
              </button>
              <button
                @click="handleDelete(task.id)"
                class="text-danger hover:text-red-900 text-sm font-medium"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, computed, inject, nextTick, watch } from 'vue'
  import { PlusIcon } from '@heroicons/vue/24/outline'
  import { formatDate, formatDateTime } from '@/utils/date'
  import { taskAPI, categoryAPI } from '@/api'
  import type { Task, Tag, CreateTaskRequest, UpdateTaskRequest } from '@/types'
  import { useSettings } from '@/composables/useSettings'

  // 注入全局通知功能
  const showNotification = inject('showNotification') as (
    type: 'success' | 'error',
    message: string
  ) => void

  const tasks = ref<Task[]>([])
  const availableCategories = ref<Tag[]>([])
  const loading = ref(false)
  const submitting = ref(false)
  const error = ref<string | null>(null)
  const showForm = ref(false)
  const editingTask = ref<Task | undefined>()
  const taskFormRef = ref<HTMLElement | null>(null)

  // Use settings from composable
  const {
    taskShowCompleted: showCompleted,
    taskShowSuspended: showSuspended,
    taskDateFilter: dateFilter,
  } = useSettings()

  const isEditing = computed(() => !!editingTask.value)

  const filteredTasks = computed(() => {
    // Since we're now filtering on the backend, just return all tasks
    return tasks.value
  })

  const form = reactive({
    title: '',
    description: '',
    category_id: '',
    due_date: '',
    estimated_minutes: 60,
  })

  const resetForm = () => {
    form.title = ''
    form.description = ''
    form.category_id = ''
    form.due_date = new Date().toISOString().split('T')[0] // Today's date
    form.estimated_minutes = 60
  }

  const loadEditingData = () => {
    if (editingTask.value) {
      form.title = editingTask.value.title
      // description 在数据库中可为 NULL，API 会省略该字段
      form.description = editingTask.value.description || ''
      form.category_id = editingTask.value.category_id.toString()
      form.due_date = editingTask.value.due_date.split('T')[0]
      form.estimated_minutes = editingTask.value.estimated_minutes
    } else {
      resetForm()
    }
  }

  // API 调用函数已从 @/api 导入

  const loadTasks = async () => {
    loading.value = true
    error.value = null

    try {
      let dateParam = ''
      const today = new Date().toISOString().split('T')[0]

      if (dateFilter.value === 'today') {
        dateParam = today
      } else if (dateFilter.value === 'tomorrow') {
        const tomorrow = new Date()
        tomorrow.setDate(tomorrow.getDate() + 1)
        dateParam = tomorrow.toISOString().split('T')[0]
      }

      const response = await taskAPI.getAll(dateParam, showSuspended.value, showCompleted.value)
      tasks.value = response.data || []
    } catch (err) {
      error.value = 'Failed to load tasks'
      console.error('Error loading tasks:', err)
      showNotification('error', 'Failed to load tasks')
    } finally {
      loading.value = false
    }
  }

  const loadTags = async () => {
    try {
      const response = await categoryAPI.getAll()
      availableCategories.value = response.data || []
    } catch (err) {
      console.error('Error loading tags:', err)
      showNotification('error', 'Failed to load tags')
    }
  }

  const toggleForm = () => {
    showForm.value = !showForm.value
    if (!showForm.value) {
      editingTask.value = undefined
    }
    loadEditingData()
  }

  const handleSubmit = async () => {
    submitting.value = true

    try {
      if (editingTask.value) {
        const data: UpdateTaskRequest = {
          title: form.title.trim(),
          description: form.description.trim(),
          category_id: Number(form.category_id),
          due_date: form.due_date,
          estimated_minutes: form.estimated_minutes,
        }
        await taskAPI.update(editingTask.value.id, data)
        showNotification('success', 'Task updated successfully')
      } else {
        const data: CreateTaskRequest = {
          title: form.title.trim(),
          description: form.description.trim(),
          category_id: Number(form.category_id),
          due_date: form.due_date,
          estimated_minutes: form.estimated_minutes,
        }
        await taskAPI.create(data)
        showNotification('success', 'Task created successfully')
      }

      await loadTasks()
      showForm.value = false
      editingTask.value = undefined
      resetForm()
    } catch (err) {
      console.error('Error saving task:', err)
      showNotification('error', 'Failed to save task')
    } finally {
      submitting.value = false
    }
  }

  const handleEdit = (task: Task) => {
    editingTask.value = task
    showForm.value = true
    loadEditingData()

    // 滚动到表单位置
    nextTick(() => {
      if (taskFormRef.value) {
        taskFormRef.value.scrollIntoView({
          behavior: 'smooth',
          block: 'start',
        })
      }
    })
  }

  const handleCancel = () => {
    showForm.value = false
    editingTask.value = undefined
    resetForm()
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this task? This action cannot be undone.')) {
      return
    }

    try {
      await taskAPI.delete(id)
      showNotification('success', 'Task deleted successfully')
      await loadTasks()
    } catch (err) {
      console.error('Error deleting task:', err)
      showNotification('error', 'Failed to delete task')
    }
  }

  const completeTask = async (id: number) => {
    try {
      await taskAPI.complete(id)
      showNotification('success', 'Task marked as completed')
      await loadTasks()
    } catch (err) {
      console.error('Error completing task:', err)
      showNotification('error', 'Failed to complete task')
    }
  }

  const incompleteTask = async (id: number) => {
    try {
      await taskAPI.incomplete(id)
      showNotification('success', 'Task reopened')
      await loadTasks()
    } catch (err) {
      console.error('Error reopening task:', err)
      showNotification('error', 'Failed to reopen task')
    }
  }

  const suspendTask = async (id: number) => {
    try {
      await taskAPI.suspend(id)
      showNotification('success', 'Task suspended')
      await loadTasks()
    } catch (err) {
      console.error('Error suspending task:', err)
      showNotification('error', 'Failed to suspend task')
    }
  }

  const unsuspendTask = async (id: number) => {
    try {
      await taskAPI.unsuspend(id)
      showNotification('success', 'Task unsuspended')
      await loadTasks()
    } catch (err) {
      console.error('Error unsuspending task:', err)
      showNotification('error', 'Failed to unsuspend task')
    }
  }

  watch(showCompleted, () => {
    loadTasks()
  })

  watch(showSuspended, () => {
    loadTasks()
  })

  onMounted(() => {
    loadTasks()
    loadTags()
  })
</script>
