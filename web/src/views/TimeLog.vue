<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-gray-900">Time Logs</h1>
      <div class="flex gap-2">
        <button
          @click="handleRefresh"
          class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-600 border border-transparent rounded-md shadow-sm hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500"
        >
          Refresh
        </button>
        <button
          @click="toggleForm"
          class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <PlusIcon class="h-5 w-5 mr-2" />
          New Log
        </button>
      </div>
    </div>

    <!-- Time Filter -->
    <div class="bg-white shadow rounded-lg p-4">
      <div class="flex flex-wrap items-center gap-4">
        <span class="text-sm font-medium text-gray-700">Filter:</span>
        <div class="flex gap-2">
          <button
            v-for="option in filterOptions"
            :key="option.value"
            @click="setFilter(option.value)"
            :class="[
              'px-3 py-1.5 text-sm font-medium rounded-md transition-colors',
              activeFilter === option.value
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
            ]"
          >
            {{ option.label }}
          </button>
        </div>
        <span class="text-sm text-gray-500 ml-auto"> {{ filteredTimeLogs.length }} entries </span>
      </div>
    </div>

    <TimeLogForm
      v-if="showForm"
      :editing-log="editingLog"
      :submitting="submitting"
      :available-categories="categories"
      :available-tasks="tasks"
      :available-constraints="constraints"
      :last-end-time="getLastEndTime()"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />

    <TimeLogList
      :time-logs="filteredTimeLogs"
      :categories="categories"
      :loading="loading"
      :error="error"
      @edit="handleEdit"
      @delete="handleDelete"
    />
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, inject, computed } from 'vue'
  import { PlusIcon } from '@heroicons/vue/24/outline'
  import { startOfDay, startOfWeek, isAfter, parseISO } from 'date-fns'
  import TimeLogList from '@/components/TimeLogList.vue'
  import TimeLogForm from '@/components/TimeLogForm.vue'
  import { timelogAPI, categoryAPI, taskAPI, constraintAPI } from '@/api'
  import { useTimeLogStore } from '@/stores/timelog'
  import type {
    TimeLog,
    Category,
    Task,
    Constraint,
    CreateTimeLogRequest,
    UpdateTimeLogRequest,
  } from '@/types'

  // 注入全局通知功能
  const showNotification = inject('showNotification') as (
    type: 'success' | 'error',
    message: string
  ) => void

  // 使用Pinia store
  const timeLogStore = useTimeLogStore()
  const timeLogs = computed(() => timeLogStore.timeLogs)
  const loading = computed(() => timeLogStore.loading)
  const error = computed(() => timeLogStore.error)

  const categories = ref<Category[]>([])
  const tasks = ref<Task[]>([])
  const constraints = ref<Constraint[]>([])
  const submitting = ref(false)
  const showForm = ref(false)
  const editingLog = ref<TimeLog | undefined>()

  // Filter state
  type FilterValue = 'today' | 'week' | 'all'
  const activeFilter = ref<FilterValue>('today')
  const filterOptions: { label: string; value: FilterValue }[] = [
    { label: 'Today', value: 'today' },
    { label: 'This Week', value: 'week' },
    { label: 'All', value: 'all' },
  ]

  const setFilter = (value: FilterValue) => {
    activeFilter.value = value
  }

  const filteredTimeLogs = computed(() => {
    if (activeFilter.value === 'all') {
      return timeLogs.value
    }

    const now = new Date()
    let filterDate: Date

    if (activeFilter.value === 'today') {
      filterDate = startOfDay(now)
    } else {
      // week - start from Monday
      filterDate = startOfWeek(now, { weekStartsOn: 1 })
    }

    return timeLogs.value.filter(log => {
      const logDate = parseISO(log.start_time)
      return isAfter(logDate, filterDate) || logDate.getTime() === filterDate.getTime()
    })
  })

  const handleRefresh = async () => {
    await timeLogStore.refreshTimeLogs()
    showNotification('success', 'Time logs refreshed successfully')
  }

  const loadCategories = async () => {
    try {
      const response = await categoryAPI.getAll()
      categories.value = response.data || []
    } catch (err) {
      console.error('Error loading categories:', err)
      showNotification('error', 'Failed to load categories')
    }
  }

  const loadTasks = async () => {
    try {
      const response = await taskAPI.getAll()
      // Filter out completed tasks when associating with timelogs
      tasks.value = (response.data || []).filter(task => !task.is_completed)
    } catch (err) {
      console.error('Error loading tasks:', err)
      showNotification('error', 'Failed to load tasks')
    }
  }

  const loadConstraints = async () => {
    try {
      const response = await constraintAPI.getAll(true) // Get only active constraints
      constraints.value = response.data || []
    } catch (err) {
      console.error('Error loading constraints:', err)
      // Don't show notification for constraint loading errors, as it's optional
    }
  }

  const toggleForm = () => {
    showForm.value = !showForm.value
    if (!showForm.value) {
      editingLog.value = undefined
    }
  }

  const handleSubmit = async (data: CreateTimeLogRequest | UpdateTimeLogRequest) => {
    submitting.value = true

    try {
      if (editingLog.value) {
        if (!editingLog.value.id) {
          showNotification('error', 'Invalid time log ID')
          return
        }
        await timelogAPI.update(editingLog.value.id, data as UpdateTimeLogRequest)
        showNotification('success', 'Time log updated successfully')
      } else {
        await timelogAPI.create(data as CreateTimeLogRequest)
        showNotification('success', 'Time log created successfully')
      }

      await timeLogStore.refreshTimeLogs()
      showForm.value = false
      editingLog.value = undefined
    } catch (err) {
      console.error('Error saving time log:', err)
      showNotification('error', 'Failed to save time log')
    } finally {
      submitting.value = false
    }
  }

  const handleEdit = (log: TimeLog) => {
    editingLog.value = log
    showForm.value = true
  }

  const handleCancel = () => {
    showForm.value = false
    editingLog.value = undefined
  }

  const getLastEndTime = (): string | null => {
    if (timeLogs.value.length === 0) {
      return null
    }

    // 按ID排序获取最新的 timelog（ID是自增的，能准确反映创建顺序）
    const sortedLogs = [...timeLogs.value].sort((a, b) => b.id - a.id)

    const lastLog = sortedLogs[0]
    return lastLog?.end_time || null
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this time log?')) {
      return
    }

    try {
      await timelogAPI.delete(id)
      showNotification('success', 'Time log deleted successfully')
      await timeLogStore.refreshTimeLogs()
    } catch (err) {
      console.error('Error deleting time log:', err)
      showNotification('error', 'Failed to delete time log')
    }
  }

  onMounted(() => {
    timeLogStore.loadTimeLogs()
    loadCategories()
    loadTasks()
    loadConstraints()
  })
</script>
